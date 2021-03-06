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
	"net/http"
	"strings"

	"github.com/rileyr/middleware/wares"

	"github.com/spf13/pflag"

	"github.com/SuperGreenLab/AppBackend/internal/server/tools"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"upper.io/db.v3/lib/sqlbuilder"
)

var (
	_ = pflag.String("jwtsecret", "", "JWT secret")
	_ = pflag.String("logrequests", "true", "Set to false in production") // TODO move this somewhere else
)

func init() {
	viper.SetDefault("JWTSecret", "")
	viper.SetDefault("LogRequests", "true")
}

// AnonStack - allows anonymous connection
func AnonStack() middleware.Stack {
	anon := middleware.NewStack()
	if viper.GetString("LogRequests") == "true" {
		anon.Use(wares.Logging)
	}
	anon.Use(CreateDBSession)
	return anon
}

// AuthStack - Decodes JWT token, errors on failure
func AuthStack() middleware.Stack {
	auth := middleware.NewStack()
	if viper.GetString("LogRequests") == "true" {
		auth.Use(wares.Logging)
	}
	auth.Use(JwtToken)
	auth.Use(UserIDRequired)
	auth.Use(CreateDBSession)
	return auth
}

// OptionalAuthStack - Decodes JWT token, errors on failure
func OptionalAuthStack() middleware.Stack {
	auth := middleware.NewStack()
	if viper.GetString("LogRequests") == "true" {
		auth.Use(wares.Logging)
	}
	auth.Use(JwtToken)
	auth.Use(CreateDBSession)
	return auth
}

// SetUserID - sets the userID field for the object payload
func SetUserID(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		o := r.Context().Value(ObjectContextKey{}).(appbackend.UserObject)
		uid := r.Context().Value(UserIDContextKey{}).(uuid.UUID)

		o.SetUserID(uid)

		ctx := context.WithValue(r.Context(), ObjectContextKey{}, o)
		fn(w, r.WithContext(ctx), p)
	}
}

// CheckAccessRight - checks if the user has access to the given object
func CheckAccessRight(collection, field string, optional bool, factory func() appbackend.UserObject) middleware.Middleware {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			o := r.Context().Value(ObjectContextKey{}).(appbackend.UserObject)
			uid := r.Context().Value(UserIDContextKey{}).(uuid.UUID)
			sess := r.Context().Value(SessContextKey{}).(sqlbuilder.Database)

			if err := tools.CheckUserID(sess, uid, o, collection, field, optional, factory); err != nil {
				errorMsg := "Object is owned by another user"
				logrus.Errorf("CheckUserID in CheckAccessRight '%s' %q for uid: %s o.GetUserID: %s", errorMsg, err, uid, o.GetUserID())
				http.Error(w, errorMsg, http.StatusUnauthorized)
				return
			}

			fn(w, r, p)
		}
	}
}

// JwtClaimsContextKey - context key which stores the request's jwt.MapClaims
type JwtClaimsContextKey struct{}

// UserIDContextKey - context key which stores the request's userID
type UserIDContextKey struct{}

// JwtToken - decodes the JWT token for the request
func JwtToken(fn httprouter.Handle) httprouter.Handle {
	hmacSampleSecret := []byte(viper.GetString("JWTSecret"))

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		authentication := r.Header.Get("Authentication") // Ooops.. mistyped:/ will remove it when the app uses the right header
		authorization := r.Header.Get("Authorization")
		if authorization != "" {
			authentication = authorization
		}
		tokenString := strings.ReplaceAll(authentication, "Bearer ", "")
		if tokenString == "" || tokenString == "null" {
			fn(w, r, p)
			return
		}

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
			ctx := context.WithValue(r.Context(), JwtClaimsContextKey{}, claims)
			ctx = context.WithValue(ctx, UserIDContextKey{}, uuid.FromStringOrNil(claims["userID"].(string)))
			fn(w, r.WithContext(ctx), p)
		} else {
			logrus.Errorln(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}
}
