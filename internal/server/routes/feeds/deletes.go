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

package feeds

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

type deletesRequest struct {
	Deletes []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}
}

var factories map[string]func() UserObject = map[string]func() UserObject{
	"boxes":       func() UserObject { return &Box{} },
	"plants":      func() UserObject { return &Plant{} },
	"timelapses":  func() UserObject { return &Timelapse{} },
	"devices":     func() UserObject { return &Device{} },
	"feeds":       func() UserObject { return &Feed{} },
	"feedentries": func() UserObject { return &FeedEntry{} },
	"feedmedias":  func() UserObject { return &FeedMedia{} },
}

var idFields map[string]string = map[string]string{
	"boxes":       "boxid",
	"plants":      "plantid",
	"timelapses":  "timelapseid",
	"devices":     "deviceid",
	"feeds":       "feedid",
	"feedentries": "feedentryid",
	"feedmedias":  "feedmediaid",
}

func createDeleteHandler() httprouter.Handle {
	s := middleware.NewStack()

	s.Use(decodeJSON(func() interface{} {
		return &deletesRequest{}
	}))

	return s.Wrap(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		uid := r.Context().Value(userIDContextKey{}).(uuid.UUID)
		sess := r.Context().Value(sessContextKey{}).(sqlbuilder.Database)
		deletes := r.Context().Value(objectContextKey{}).(*deletesRequest)
		ueid := r.Context().Value(userEndIDContextKey{}).(uuid.UUID)

		for _, del := range deletes.Deletes {
			factory, ok := factories[del.Type]
			if ok == false {
				logrus.Warningf("Unknown type %s", del.Type)
				continue
			}
			o := factory()
			err := sess.Collection(del.Type).Find("id", del.ID).One(o)
			if err != nil {
				logrus.Errorln(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if uid != o.GetUserID() {
				logrus.Warningf("Wrong userID %s %s", del.Type, del.ID)
				continue
			}

			if _, err := sess.Update(del.Type).Set("deleted", true).Where("id = ?", o.GetID()).Exec(); err != nil {
				logrus.Warning(err.Error())
				continue
			}

			field := idFields[del.Type]
			if _, err := sess.Update(fmt.Sprintf("userend_%s", del.Type)).Set("dirty", true).Where(field, del.ID).And("userendid != ?", ueid).And("userendid in (select id from userends where userid = ?)", uid).Exec(); err != nil {
				logrus.Warning(err.Error())
				continue
			}
		}
	})
}

var deletesHandler = createDeleteHandler()
