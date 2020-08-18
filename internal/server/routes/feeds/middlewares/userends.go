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
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	cmiddlewares "github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

// CreateUserEndObjects - creates the UserEnd object associated with the inserted object
func CreateUserEndObjects(collection string, factory func() db.UserEndObject) middleware.Middleware {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(cmiddlewares.SessContextKey{}).(sqlbuilder.Database)
			uid := r.Context().Value(cmiddlewares.UserIDContextKey{}).(uuid.UUID)
			ueid := r.Context().Value(UserEndIDContextKey{}).(uuid.UUID)

			id := r.Context().Value(cmiddlewares.InsertedIDContextKey{}).(uuid.UUID)

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
			sess := r.Context().Value(cmiddlewares.SessContextKey{}).(sqlbuilder.Database)
			uid := r.Context().Value(cmiddlewares.UserIDContextKey{}).(uuid.UUID)
			ueid := r.Context().Value(UserEndIDContextKey{}).(uuid.UUID)

			id := r.Context().Value(cmiddlewares.UpdatedIDContextKey{}).(uuid.UUID)

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
