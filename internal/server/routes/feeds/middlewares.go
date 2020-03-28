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
		uid := r.Context().Value(userIDContextKey{}).(string)

		reflect.ValueOf(o).Elem().FieldByName("UserID").SetString(uid)

		ctx := context.WithValue(r.Context(), objectContextKey{}, o)
		fn(w, r.WithContext(ctx), p)
	}
}

type insertedIDContextKey struct{}

func insertObject(collection string) func(fn httprouter.Handle) httprouter.Handle {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			o := r.Context().Value(objectContextKey{})
			sess := r.Context().Value(sessContextKey{}).(sqlbuilder.Database)
			users := sess.Collection(collection)
			id, err := users.Insert(o)
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
			ctx := context.WithValue(r.Context(), userIDContextKey{}, claims["userID"])
			ctx = context.WithValue(ctx, userEndIDContextKey{}, claims["userEndID"])
			fn(w, r.WithContext(ctx), p)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}
}
