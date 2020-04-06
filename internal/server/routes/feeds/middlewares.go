package feeds

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/rileyr/middleware"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

type sessContextKey struct{}

func createDBSession(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var err error
		sess, err := postgresql.Open(settings)
		if err != nil {
			logrus.Errorf("db.Open(): %q\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer sess.Close()

		ctx := context.WithValue(r.Context(), sessContextKey{}, sess)
		fn(w, r.WithContext(ctx), p)
	}
}

type objectContextKey struct{}

func decodeJSON(fnObject func() interface{}) func(fn httprouter.Handle) httprouter.Handle {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			o := fnObject()
			err := decodeJSONBody(w, r, o)
			if err != nil {
				var mr *malformedRequest
				if errors.As(err, &mr) {
					http.Error(w, mr.msg, mr.status)
				} else {
					log.Println(err.Error())
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				return
			}
			ctx := context.WithValue(r.Context(), objectContextKey{}, o)
			fn(w, r.WithContext(ctx), p)
		}
	}
}

func setUserID(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		o := r.Context().Value(objectContextKey{})
		uid := r.Context().Value(userIDContextKey{}).(uuid.UUID)

		reflect.ValueOf(o).Elem().FieldByName("UserID").Set(reflect.ValueOf(uid))

		ctx := context.WithValue(r.Context(), objectContextKey{}, o)
		fn(w, r.WithContext(ctx), p)
	}
}

func checkAccessRight(collection, field string, optional bool, factory func() interface{}) middleware.Middleware {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			o := r.Context().Value(objectContextKey{})
			uid := r.Context().Value(userIDContextKey{}).(uuid.UUID)
			sess := r.Context().Value(sessContextKey{}).(sqlbuilder.Database)

			var id uuid.UUID
			idFieldValue := reflect.ValueOf(o).Elem().FieldByName(field).Interface()
			if v, ok := idFieldValue.(uuid.UUID); ok == true {
				id = v
			} else if v, ok := idFieldValue.(uuid.NullUUID); ok == true {
				if !v.Valid && !optional {
					http.Error(w, "Access denied", http.StatusUnauthorized)
					return
				} else if !v.Valid && optional {
					fn(w, r, p)
					return
				}
				id = v.UUID
			}

			parent := factory()
			err := sess.Collection(collection).Find("id", id).One(parent)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			uidParent := reflect.ValueOf(parent).Elem().FieldByName("UserID").Interface().(uuid.UUID)

			if !uuid.Equal(uid, uidParent) {
				http.Error(w, "Access denied", http.StatusUnauthorized)
				return
			}

			fn(w, r, p)
		}
	}
}

type insertedIDContextKey struct{}

func insertObject(collection string) func(fn httprouter.Handle) httprouter.Handle {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			o := r.Context().Value(objectContextKey{})
			sess := r.Context().Value(sessContextKey{}).(sqlbuilder.Database)
			col := sess.Collection(collection)
			id, err := col.Insert(o)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), insertedIDContextKey{}, string(id.([]uint8)))
			fn(w, r.WithContext(ctx), p)
		}
	}
}

type userIDContextKey struct{}
type userEndIDContextKey struct{}

func jwtToken(fn httprouter.Handle) httprouter.Handle {
	hmacSampleSecret := []byte(viper.GetString("JWTSecret"))

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		authentication := r.Header.Get("Authentication")
		tokenString := strings.ReplaceAll(authentication, "Bearer ", "")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return hmacSampleSecret, nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), userIDContextKey{}, uuid.FromStringOrNil(claims["userID"].(string)))
			if userEndID, ok := claims["userEndID"]; ok == true {
				ctx = context.WithValue(ctx, userEndIDContextKey{}, uuid.FromStringOrNil(userEndID.(string)))
			}
			fn(w, r.WithContext(ctx), p)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}
}

func userEndIDRequired(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ueid := r.Context().Value(userEndIDContextKey{})
		if ueid == nil {
			http.Error(w, "Missing userEndID", http.StatusBadRequest)
			return
		}
		fn(w, r, p)
	}
}
