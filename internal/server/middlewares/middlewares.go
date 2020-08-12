/*
 * Copyright (C) 2020  SuperGreenLab <towelie@supergreenlab.com>
 * Author: Constantin Clauzel <constantin.clauzel@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/SuperGreenLab/AppBackend/internal/server/tools"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"upper.io/db.v3/lib/sqlbuilder"
)

// ObjectContextKey - context key which stores the decoced object
type ObjectContextKey struct{}

// DecodeJSON - decodes the JSON payload
func DecodeJSON(fnObject func() interface{}) func(fn httprouter.Handle) httprouter.Handle {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			o := fnObject()
			err := tools.DecodeJSONBody(w, r, o)
			if err != nil {
				var mr *tools.MalformedRequest
				if errors.As(err, &mr) {
					logrus.Errorln(err.Error())
					http.Error(w, mr.Msg, mr.Status)
				} else {
					log.Println(err.Error())
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				return
			}
			ctx := context.WithValue(r.Context(), ObjectContextKey{}, o)
			fn(w, r.WithContext(ctx), p)
		}
	}
}

// SetUserID - sets the userID field for the object payload
func SetUserID(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		o := r.Context().Value(ObjectContextKey{}).(db.UserObject)
		uid := r.Context().Value(UserIDContextKey{}).(uuid.UUID)

		o.SetUserID(uid)

		ctx := context.WithValue(r.Context(), ObjectContextKey{}, o)
		fn(w, r.WithContext(ctx), p)
	}
}

// CheckAccessRight - checks if the user has access to the given object
func CheckAccessRight(collection, field string, optional bool, factory func() db.UserObject) middleware.Middleware {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			o := r.Context().Value(ObjectContextKey{}).(db.UserObject)
			uid := r.Context().Value(UserIDContextKey{}).(uuid.UUID)
			sess := r.Context().Value(SessContextKey{}).(sqlbuilder.Database)

			if err := tools.CheckUserID(sess, uid, o, collection, field, optional, factory); err != nil {
				logrus.Errorln(err.Error())
				http.Error(w, "Object is owned by another user", http.StatusUnauthorized)
				return
			}

			fn(w, r, p)
		}
	}
}

// UserIDContextKey - context key which stores the request's userID
type UserIDContextKey struct{}

// UserEndIDContextKey - context key which stores the request's userEndID
type UserEndIDContextKey struct{}

// JwtToken - decodes the JWT token for the request
func JwtToken(fn httprouter.Handle) httprouter.Handle {
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
			logrus.Errorln(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), UserIDContextKey{}, uuid.FromStringOrNil(claims["userID"].(string)))
			if userEndID, ok := claims["userEndID"]; ok == true {
				ctx = context.WithValue(ctx, UserEndIDContextKey{}, uuid.FromStringOrNil(userEndID.(string)))
			}
			fn(w, r.WithContext(ctx), p)
		} else {
			logrus.Errorln(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}
}

// UserEndIDRequired - Checks if the request has a userID
func UserEndIDRequired(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ueid := r.Context().Value(UserEndIDContextKey{})
		if ueid == nil {
			logrus.Errorln("Missing userEndID")
			http.Error(w, "Missing userEndID", http.StatusBadRequest)
			return
		}
		fn(w, r, p)
	}
}

// ObjectIDRequired - Checks if the object's id is set in the payload
func ObjectIDRequired(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		o := r.Context().Value(ObjectContextKey{}).(db.Object)
		if o.GetID().Valid == false {
			logrus.Errorln("Missing object's ID")
			http.Error(w, "Missing object's ID", http.StatusBadRequest)
			return
		}
		fn(w, r, p)
	}
}

// CreateUserEndObjects - creates the UserEnd object associated with the inserted object
func CreateUserEndObjects(collection string, factory func() db.UserEndObject) middleware.Middleware {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(SessContextKey{}).(sqlbuilder.Database)
			uid := r.Context().Value(UserIDContextKey{}).(uuid.UUID)
			ueid := r.Context().Value(UserEndIDContextKey{}).(uuid.UUID)

			id := r.Context().Value(InsertedIDContextKey{}).(uuid.UUID)

			uends := []db.UserEnd{}
			err := sess.Collection("userends").Find("userid", uid).All(&uends)
			if err != nil {
				logrus.Errorln(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for _, uend := range uends {
				ueo := factory()
				ueo.SetObjectID(id)
				ueo.SetUserEndID(uend.ID.UUID)
				if uend.ID.UUID == ueid {
					ueo.SetSent(true)
				} else {
					ueo.SetDirty(true)
				}
				sess.Collection(collection).Insert(ueo)
			}

			fn(w, r, p)
		}
	}
}

// UpdateUserEndObjects - sets the UserEnd object to dirty when updated
func UpdateUserEndObjects(collection, field string) middleware.Middleware {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(SessContextKey{}).(sqlbuilder.Database)
			uid := r.Context().Value(UserIDContextKey{}).(uuid.UUID)
			ueid := r.Context().Value(UserEndIDContextKey{}).(uuid.UUID)

			id := r.Context().Value(UpdatedIDContextKey{}).(uuid.UUID)

			_, err := sess.Update(collection).Set("dirty", true).Where(field, id).And("userendid != ?", ueid).And("userendid in (select id from userends where userid = ?)", uid).Exec()
			if err != nil {
				logrus.Errorln(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fn(w, r, p)
		}
	}
}
