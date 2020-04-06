package feeds

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"upper.io/db.v3/lib/sqlbuilder"
)

type syncData struct {
	Boxes       []Box       `json:"boxes"`
	Plants      []Plant     `json:"plants"`
	Timelapses  []Box       `json:"timelapses"`
	Devices     []Device    `json:"devices"`
	Feeds       []Feed      `json:"feeds"`
	FeedEntries []FeedEntry `json:"feedEntries"`
	FeedMedias  []FeedMedia `json:"feedMedias"`
}

func syncCollection(w http.ResponseWriter, r *http.Request, collection, id string, res interface{}) error {
	var err error
	sess := r.Context().Value(sessContextKey{}).(sqlbuilder.Database)
	ueid := r.Context().Value(userEndIDContextKey{}).(uuid.UUID)
	err = sess.Select("a.*").From(fmt.Sprintf("%s a", collection)).Join(fmt.Sprintf("userend_%s b", collection)).On(fmt.Sprintf("b.%s = a.id", id)).Where("b.userendid = ?", ueid).And("dirty = true").All(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

func syncHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := syncData{}

	if err := syncCollection(w, r, "boxes", "boxid", &res.Boxes); err != nil {
		return
	}
	if err := syncCollection(w, r, "plants", "plantid", &res.Plants); err != nil {
		return
	}
	if err := syncCollection(w, r, "timelapses", "timelapseid", &res.Timelapses); err != nil {
		return
	}
	if err := syncCollection(w, r, "devices", "deviceid", &res.Devices); err != nil {
		return
	}
	if err := syncCollection(w, r, "feeds", "feedid", &res.Feeds); err != nil {
		return
	}
	if err := syncCollection(w, r, "feedentries", "feedentryid", &res.FeedEntries); err != nil {
		return
	}
	if err := syncCollection(w, r, "feedmedias", "feedmediaid", &res.FeedMedias); err != nil {
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
