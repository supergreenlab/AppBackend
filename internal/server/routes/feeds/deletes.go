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

	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	fmiddlewares "github.com/SuperGreenLab/AppBackend/internal/server/routes/feeds/middlewares"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"

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
	} `json:"deletes"`
}

var factories map[string]func() appbackend.UserObject = map[string]func() appbackend.UserObject{
	"boxes":       func() appbackend.UserObject { return &appbackend.Box{} },
	"plants":      func() appbackend.UserObject { return &appbackend.Plant{} },
	"timelapses":  func() appbackend.UserObject { return &appbackend.Timelapse{} },
	"devices":     func() appbackend.UserObject { return &appbackend.Device{} },
	"feeds":       func() appbackend.UserObject { return &appbackend.Feed{} },
	"feedentries": func() appbackend.UserObject { return &appbackend.FeedEntry{} },
	"feedmedias":  func() appbackend.UserObject { return &appbackend.FeedMedia{} },
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

	s.Use(middlewares.DecodeJSON(func() interface{} {
		return &deletesRequest{}
	}))

	return s.Wrap(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		uid := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
		deletes := r.Context().Value(middlewares.ObjectContextKey{}).(*deletesRequest)
		ueid, ueidOK := r.Context().Value(fmiddlewares.UserEndIDContextKey{}).(uuid.UUID)

		for _, del := range deletes.Deletes {
			factory, ok := factories[del.Type]
			if ok == false {
				logrus.Warningf("Unknown type %s by %s", del.Type, uid)
				continue
			}
			o := factory()
			err := sess.Collection(del.Type).Find("id", del.ID).One(o)
			if err != nil {
				logrus.Errorf("sess.Collection.Find in createDeleteHandler %q - %+v by %s", err, del, uid)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if uid != o.GetUserID() {
				logrus.Warningf("Object is owned by another user - %+v", del)
				continue
			}

			if _, err := sess.Update(del.Type).Set("deleted", true).Where("id = ?", o.GetID()).Exec(); err != nil {
				logrus.Warningf("sess.Update(del.Type) in createDeleteHandler %q - %+v %+v by %s", err, del, o, uid)
				continue
			}

			collection := fmt.Sprintf("userend_%s", del.Type)
			field := idFields[del.Type]
			ueUpdate := sess.Update(collection).Set("dirty", true).Where(field, del.ID)
			if ueidOK {
				ueUpdate = ueUpdate.And("userendid != ?", ueid)
			}
			ueUpdate = ueUpdate.And("userendid in (select id from userends where userid = ?)", uid)
			if _, err := ueUpdate.Exec(); err != nil {
				logrus.Warningf("sess.Update(collection) in createDeleteHandler %q - %+v by %s", err, del, uid)
				continue
			}

			if ueidOK {
				if _, err := sess.DeleteFrom(collection).Where(fmt.Sprintf("%s = ?", field), del.ID).And("userendid = ?", ueid).Exec(); err != nil {
					logrus.Warningf("sess.DeleteFrom(collection) in createDeleteHandler %q - %+v by %s", err, del, uid)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		}
	})
}

var deletesHandler = createDeleteHandler()
