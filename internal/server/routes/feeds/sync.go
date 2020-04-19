package feeds

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
				customSelect(selector)
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
var syncFeedEntriesHandler = syncCollection("feedentries", "feedentryid", func() interface{} { return &[]FeedEntry{} }, nil, nil)
var syncFeedMediasHandler = syncCollection("feedmedias", "feedmediaid", func() interface{} { return &[]FeedMedia{} }, nil, []middleware.Middleware{
	func(fn httprouter.Handle) httprouter.Handle {
		expiry := time.Second * 60 * 60
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			minioClient := createMinioClient()
			feedMedias := r.Context().Value(objectContextKey{}).(*[]FeedMedia)
			for i, fm := range *feedMedias {
				url1, err := minioClient.PresignedGetObject("feedmedias", fm.FilePath, expiry, nil)
				if err != nil {
					logrus.Errorln(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				fm.FilePath = url1.RequestURI()

				url2, err := minioClient.PresignedGetObject("feedmedias", fm.ThumbnailPath, expiry, nil)
				if err != nil {
					logrus.Errorln(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				fm.ThumbnailPath = url2.RequestURI()
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
		_, err := sess.Update(collection).Set("sent", true, "dirty", false).Where(fmt.Sprintf("%s = ?", field), p.ByName("id")).And("userendid = ?", ueid).Exec()
		if err != nil {
			logrus.Errorln(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
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
