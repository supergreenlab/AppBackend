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
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// ObjectIDRequired - Checks if the object's id is set in the payload
func ObjectIDRequired(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		o := r.Context().Value(ObjectContextKey{}).(db.Object)
		if o.GetID().Valid == false {
			logrus.Errorf("Missing object's ID - %+v", o)
			http.Error(w, "Missing object's ID", http.StatusBadRequest)
			return
		}
		fn(w, r, p)
	}
}

// UserIDRequired - Checks if the request has a userID
func UserIDRequired(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ueid := r.Context().Value(UserIDContextKey{})
		if ueid == nil {
			logrus.Errorln("Missing userID")
			http.Error(w, "Missing userID", http.StatusBadRequest)
			return
		}
		fn(w, r, p)
	}
}
