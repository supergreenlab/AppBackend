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
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	cmiddlewares "github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	fmiddlewares "github.com/SuperGreenLab/AppBackend/internal/server/routes/feeds/middlewares"
	"github.com/SuperGreenLab/AppBackend/internal/server/tools"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
	udb "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

type syncResponse struct {
	Items interface{} `json:"items"`
}

func syncCollection(collection, id string, factory func() interface{}, customSelect func(sqlbuilder.Selector) sqlbuilder.Selector, postSelect []middleware.Middleware) httprouter.Handle {
	s := middleware.NewStack()

	s.Use(func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(cmiddlewares.SessContextKey{}).(sqlbuilder.Database)
			ueid := r.Context().Value(fmiddlewares.UserEndIDContextKey{}).(uuid.UUID)
			res := factory()
			selector := sess.Select(udb.Raw("a.*")).From(fmt.Sprintf("%s a", collection)).Join(fmt.Sprintf("userend_%s b", collection)).On(fmt.Sprintf("b.%s = a.id", id)).Where("b.userendid = ?", ueid).And("dirty = true")
			if customSelect != nil {
				selector = customSelect(selector)
			}
			if err := selector.OrderBy("a.cat ASC").All(res); err != nil {
				logrus.Errorf("selector.OrderBy in syncCollection %q - collection: %s id: %s ueid: %s", err, collection, id, ueid)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), cmiddlewares.ObjectContextKey{}, res)
			fn(w, r.WithContext(ctx), p)
		}
	})

	if postSelect != nil {
		for _, m := range postSelect {
			s.Use(m)
		}
	}

	return s.Wrap(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		o := r.Context().Value(middlewares.ObjectContextKey{})
		if err := json.NewEncoder(w).Encode(syncResponse{o}); err != nil {
			logrus.Errorf("json.NewEncoder in syncCollection %q - collection: %s id: %s o: %+v", err, collection, id, o)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

var syncBoxesHandler = syncCollection("boxes", "boxid", func() interface{} { return &[]db.Box{} }, nil, nil)

var syncPlantsHandler = syncCollection("plants", "plantid", func() interface{} { return &[]db.Plant{} }, nil, nil)

var syncTimelapsesHandler = syncCollection("timelapses", "timelapseid", func() interface{} { return &[]db.Timelapse{} }, nil, nil)

var syncDevicesHandler = syncCollection("devices", "deviceid", func() interface{} { return &[]db.Device{} }, nil, nil)

var syncFeedsHandler = syncCollection("feeds", "feedid", func() interface{} { return &[]db.Feed{} }, func(selector sqlbuilder.Selector) sqlbuilder.Selector {
	// TODO this should be filtered on userend creation
	return selector.And("isnewsfeed", false)
}, nil)

var syncFeedEntriesHandler = syncCollection("feedentries", "feedentryid", func() interface{} { return &[]db.FeedEntry{} }, func(selector sqlbuilder.Selector) sqlbuilder.Selector {
	return selector.Join("feeds f").On("f.id = a.feedid").Where("f.isnewsfeed", false)
}, nil)

type FeedMediaWithArchived struct {
	db.FeedMedia
	PlantArchived sql.NullBool `json:"-" db:"plant_archived"`
	BoxArchived   sql.NullBool `json:"-" db:"box_archived"`
}

// TODO DRY with explorer/models.go:123
func (r *FeedMediaWithArchived) SetURLs(paths []string) {
	r.FilePath = paths[0]
	r.ThumbnailPath = paths[1]
}

func (r FeedMediaWithArchived) GetURLs() []tools.S3Path {
	return []tools.S3Path{
		tools.S3Path{
			Path:   &r.FilePath,
			Bucket: "feedmedias",
		},
		tools.S3Path{
			Path:   &r.ThumbnailPath,
			Bucket: "feedmedias",
		},
	}
}

var syncFeedMediasHandler = syncCollection("feedmedias", "feedmediaid", func() interface{} { return &[]FeedMediaWithArchived{} }, func(selector sqlbuilder.Selector) sqlbuilder.Selector {
	selector = selector.Join("feedentries fe").On("fe.id = a.feedentryid")
	selector = selector.Columns(udb.Raw("p.archived as plant_archived")).LeftJoin("plants p").On("p.feedid = fe.feedid")
	selector = selector.Columns(udb.Raw("boxes.archived as box_archived")).LeftJoin("boxes").On("boxes.feedid = fe.feedid")
	return selector
}, []middleware.Middleware{
	func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			var err error
			feedMedias := r.Context().Value(middlewares.ObjectContextKey{}).(*[]FeedMediaWithArchived)
			for i, fm := range *feedMedias {
				if fm.Deleted == false && fm.PlantArchived.Bool == false && fm.BoxArchived.Bool == false {
					err = tools.LoadFeedMediaPublicURLs(&fm)
					if err != nil {
						logrus.Errorf("tools.LoadFeedMediaPublicURLs in syncFeedMediasHandler %q - fm: %+v", err, fm)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				} else {
					logrus.Infof("Skipped %+v", fm)
				}
				// might not be useful anymore
				(*feedMedias)[i] = fm
			}
			ctx := context.WithValue(r.Context(), middlewares.ObjectContextKey{}, feedMedias)
			fn(w, r.WithContext(ctx), p)
		}
	},
})

func syncedHandler(collection, field string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(cmiddlewares.SessContextKey{}).(sqlbuilder.Database)
		ueid := r.Context().Value(fmiddlewares.UserEndIDContextKey{}).(uuid.UUID)

		var o struct {
			Deleted  bool `db:"deleted"`
			Archived bool `db:"archived"`
		}
		fields := []interface{}{"deleted"}
		if strings.Replace(collection, "userend_", "", 1) == "plants" {
			fields = append(fields, "archived")
		}
		err := sess.Select(fields...).From(strings.Replace(collection, "userend_", "", 1)).Where("id", p.ByName("id")).One(&o)
		if err != nil {
			logrus.Errorf("sess.Select in syncedHandler %q - collection: %s field: %s id: %s o: %+v ueid: %s", err, collection, field, p.ByName("id"), o, ueid)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if o.Deleted == true || o.Archived == true {
			_, err := sess.DeleteFrom(collection).Where(fmt.Sprintf("%s = ?", field), p.ByName("id")).And("userendid = ?", ueid).Exec()
			if err != nil {
				logrus.Errorf("sess.DeleteFrom in syncedHandler %q - collection: %s field: %s id: %s o: %+v ueid: %s", err, collection, field, p.ByName("id"), o, ueid)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			_, err := sess.Update(collection).Set("sent", true, "dirty", false).Where(fmt.Sprintf("%s = ?", field), p.ByName("id")).And("userendid = ?", ueid).Exec()
			if err != nil {
				logrus.Errorf("sess.Update in syncedHandler %q - collection: %s field: %s id: %s o: %+v ueid: %s", err, collection, field, p.ByName("id"), o, ueid)
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
