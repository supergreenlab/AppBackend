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

package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/gofrs/uuid"

	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"upper.io/db.v3/lib/sqlbuilder"
)

type loginParams struct {
	Handle   string `json:"handle"`
	Password string `json:"password"`
}

func loginHandler() httprouter.Handle {
	s := middleware.NewStack()

	s.Use(middlewares.DecodeJSON(func() interface{} { return &loginParams{} }))

	return s.Wrap(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		hmacSampleSecret := []byte(viper.GetString("JWTSecret"))
		lp := r.Context().Value(middlewares.ObjectContextKey{}).(*loginParams)
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)

		u := db.User{}
		err := sess.Select("id", "password").From("users").Where("lower(nickname) = ?", strings.ToLower(lp.Handle)).One(&u)
		if err != nil {
			logrus.Errorln(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(lp.Password))
		if err != nil {
			logrus.Errorln(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID": u.ID.UUID.String(),
		})
		tokenString, err := token.SignedString(hmacSampleSecret)
		if err != nil {
			logrus.Errorln(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("x-sgl-token", tokenString)

		w.WriteHeader(http.StatusOK)
	})
}

func fillUserEnd(sess sqlbuilder.Database, ueid uuid.UUID, collection string, all db.Objects, factory func() db.UserEndObject) {
	all.Each(func(a db.Object) {
		ueo := factory()
		ueo.SetUserEndID(ueid)
		ueo.SetObjectID(a.GetID().UUID)
		ueo.SetDirty(true)
		sess.Collection(fmt.Sprintf("userend_%s", collection)).Insert(ueo)
	})
}

var createUserHandler = middlewares.InsertEndpoint(
	"users",
	func() interface{} { return &db.User{} },
	[]middleware.Middleware{
		func(fn httprouter.Handle) httprouter.Handle {
			return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				u := r.Context().Value(middlewares.ObjectContextKey{}).(*db.User)
				sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
				n, err := sess.Collection("users").Find().Where("lower(nickname) = ?", u.Nickname).Count() // TODO this is stupid
				if err != nil {
					logrus.Errorln(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if n > 0 {
					logrus.Errorln("User already exists")
					http.Error(w, "User already exists", http.StatusBadRequest)
					return
				}

				bc, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
				u.Password = string(bc)
				if err != nil {
					logrus.Errorln(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				fn(w, r, p)
			}
		},
	},
	nil,
)

func meHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
	uid := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)

	user := db.User{}
	err := sess.Collection("users").Find().Where("id = ?", uid).One(&user)
	if err != nil {
		logrus.Errorln(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = ""
	if err := json.NewEncoder(w).Encode(user); err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
