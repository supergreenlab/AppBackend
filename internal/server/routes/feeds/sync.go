package feeds

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"upper.io/db.v3/lib/sqlbuilder"
)

func syncCollection(collection, id string, factory func() interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(sessContextKey{}).(sqlbuilder.Database)
		ueid := r.Context().Value(userEndIDContextKey{}).(uuid.UUID)
		res := factory()
		if err := sess.Select("a.*").From(fmt.Sprintf("%s a", collection)).Join(fmt.Sprintf("userend_%s b", collection)).On(fmt.Sprintf("b.%s = a.id", id)).Where("b.userendid = ?", ueid).And("dirty = true").All(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

var syncBoxesHandler = syncCollection("boxes", "boxid", func() interface{} { return &[]Box{} })
var syncPlantsHandler = syncCollection("plants", "plantid", func() interface{} { return &[]Plant{} })
var syncTimelapsesHandler = syncCollection("timelapses", "timelapseid", func() interface{} { return &[]Timelapse{} })
var syncDevicesHandler = syncCollection("devices", "deviceid", func() interface{} { return &[]Device{} })
var syncFeedsHandler = syncCollection("feeds", "feedid", func() interface{} { return &[]Feed{} })
var syncFeedEntriesHandler = syncCollection("feedentries", "feedentryid", func() interface{} { return &[]FeedEntry{} })
var syncFeedMediasHandler = syncCollection("feedmedias", "feedmediaid", func() interface{} { return &[]FeedMedia{} })

func syncedHandler(collection, field string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(sessContextKey{}).(sqlbuilder.Database)
		ueid := r.Context().Value(userEndIDContextKey{}).(uuid.UUID)
		_, err := sess.Update(collection).Set("sent", true, "dirty", false).Where(fmt.Sprintf("%s = ?", field), p.ByName("id")).And("userendid = ?", ueid).Exec()
		if err != nil {
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
