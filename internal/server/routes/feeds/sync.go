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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

type syncResponse struct {
	Items interface{} `json:"items"`
}

func syncCollection(collection, id string, factory func() interface{}, customSelect func(sqlbuilder.Selector) sqlbuilder.Selector, postSelect []middleware.Middleware) httprouter.Handle {
	s := middleware.NewStack()

	s.Use(func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(sessContextKey{}).(sqlbuilder.Database)
			ueid := r.Context().Value(userEndIDContextKey{}).(uuid.UUID)
			res := factory()
			selector := sess.Select("a.*").From(fmt.Sprintf("%s a", collection)).Join(fmt.Sprintf("userend_%s b", collection)).On(fmt.Sprintf("b.%s = a.id", id)).Where("b.userendid = ?", ueid).And("dirty = true")
			if customSelect != nil {
				selector = customSelect(selector)
			}
			if err := selector.OrderBy("cat ASC").All(res); err != nil {
				logrus.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), objectContextKey{}, res)
			fn(w, r.WithContext(ctx), p)
		}
	})

	if postSelect != nil {
		for _, m := range postSelect {
			s.Use(m)
		}
	}

	return s.Wrap(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		o := r.Context().Value(objectContextKey{})
		if err := json.NewEncoder(w).Encode(syncResponse{o}); err != nil {
			logrus.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

var syncBoxesHandler = syncCollection("boxes", "boxid", func() interface{} { return &[]Box{} }, nil, nil)
var syncPlantsHandler = syncCollection("plants", "plantid", func() interface{} { return &[]Plant{} }, nil, nil)
var syncTimelapsesHandler = syncCollection("timelapses", "timelapseid", func() interface{} { return &[]Timelapse{} }, nil, nil)
var syncDevicesHandler = syncCollection("devices", "deviceid", func() interface{} { return &[]Device{} }, nil, nil)
var syncFeedsHandler = syncCollection("feeds", "feedid", func() interface{} { return &[]Feed{} }, func(selector sqlbuilder.Selector) sqlbuilder.Selector {
	return selector.And("isnewsfeed", false)
}, nil)
var syncFeedEntriesHandler = syncCollection("feedentries", "feedentryid", func() interface{} { return &[]FeedEntry{} }, func(selector sqlbuilder.Selector) sqlbuilder.Selector {
	return selector.Join("feeds f").On("f.id = a.feedid").Where("f.isnewsfeed", false)
}, nil)
var syncFeedMediasHandler = syncCollection("feedmedias", "feedmediaid", func() interface{} { return &[]FeedMedia{} }, nil, []middleware.Middleware{
	func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			var err error
			feedMedias := r.Context().Value(objectContextKey{}).(*[]FeedMedia)
			for i, fm := range *feedMedias {
				fm, err = loadFeedMediaPublicURLs(fm)
				if err != nil {
					logrus.Errorln(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				(*feedMedias)[i] = fm
			}
			ctx := context.WithValue(r.Context(), objectContextKey{}, feedMedias)
			fn(w, r.WithContext(ctx), p)
		}
	},
})

func syncedHandler(collection, field string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(sessContextKey{}).(sqlbuilder.Database)
		ueid := r.Context().Value(userEndIDContextKey{}).(uuid.UUID)

		var deleted struct {
			Deleted bool `db:"deleted"`
		}
		err := sess.Select("deleted").From(strings.Replace(collection, "userend_", "", 1)).Where("id", p.ByName("id")).One(&deleted)
		if err != nil {
			logrus.Errorln(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if deleted.Deleted == true {
			_, err := sess.DeleteFrom(collection).Where(fmt.Sprintf("%s = ?", field), p.ByName("id")).And("userendid = ?", ueid).Exec()
			if err != nil {
				logrus.Errorln(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			_, err := sess.Update(collection).Set("sent", true, "dirty", false).Where(fmt.Sprintf("%s = ?", field), p.ByName("id")).And("userendid = ?", ueid).Exec()
			if err != nil {
				logrus.Errorln(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}
}

var syncedBoxHandler = syncedHandler("userend_boxes", "boxid")
var syncedPlantHandler = syncedHandler("userend_plants", "plantid")
var syncedTimelapseHandler = syncedHandler("userend_timelapses", "timelapseid")
var syncedDeviceHandler = syncedHandler("userend_devices", "deviceid")
var syncedFeedHandler = syncedHandler("userend_feeds", "feedid")
var syncedFeedEntryHandler = syncedHandler("userend_feedentries", "feedentryid")
var syncedFeedMediaHandler = syncedHandler("userend_feedmedias", "feedmediaid")
