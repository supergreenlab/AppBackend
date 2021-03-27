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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	sgldb "github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3"
	udb "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// TODO use select* middlewares

func loadLastFeedMediaForPlant(sess sqlbuilder.Database, p sgldb.Plant) (sgldb.FeedMedia, error) {
	var err error
	selector := sess.Select("fm.*").
		From("feedmedias fm").
		Join("feedentries fe").On("fm.feedentryid = fe.id and fe.deleted = ?", false).
		Join("plants p").On("fe.feedid = p.feedid").
		Where("p.id = ?", p.ID).And("fm.deleted = ?", false).
		OrderBy("fm.cat desc").Limit(1)
	fm := sgldb.FeedMedia{}
	if err = selector.One(&fm); err != nil {
		return fm, err
	}
	err = loadFeedMediaPublicURLs(&fm)
	if err != nil {
		return fm, err
	}
	return fm, nil
}

type publicListingPlantResult struct {
	ID            string `db:"id" json:"id"`
	Name          string `db:"name" json:"name"`
	FilePath      string `db:"filepath" json:"filePath"`
	ThumbnailPath string `db:"thumbnailpath" json:"thumbnailPath"`
}

func (r *publicListingPlantResult) SetURLs(filePath string, thumbnailPath string) {
	r.FilePath = filePath
	r.ThumbnailPath = thumbnailPath
}

func (r publicListingPlantResult) GetURLs() (filePath string, thumbnailPath string) {
	filePath, thumbnailPath = r.FilePath, r.ThumbnailPath
	return
}

type publicPlantsResult struct {
	Plants []publicListingPlantResult `json:"plants"`
}

func fetchPublicPlants(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		logrus.Errorf("strconv.Atoi in fetchPublicPlants %q - offset: %s url: %s", err, r.URL.Query().Get("offset"), r.URL.String())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		logrus.Errorf("strconv.Atoi in fetchPublicPlants %q - limit: %s url: %s", err, r.URL.Query().Get("limit"), r.URL.String())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if limit < 0 {
		limit = 0
	} else if limit > 50 {
		limit = 50
	}

	lastFeedEntrySelector := sess.Select("feedid", udb.Raw("max(cat) as cat")).
		From("feedentries").
		Where("deleted = false").
		And(fmt.Sprintf("etype in ('%s')", strings.Join([]string{"FE_MEDIA", "FE_BENDING", "FE_DEFOLATION", "FE_TRANSPLANT", "FE_FIMMING", "FE_TOPPING", "FE_MEASURE"}, "', '"))).
		GroupBy("feedid")
	selector := sess.Select("plants.id", "plants.name", "feedmedias.filepath", "feedmedias.thumbnailpath").
		From("plants").
		Join(db.Raw(fmt.Sprintf("(%s) latest", lastFeedEntrySelector.String()))).Using("feedid").
		Join("feedentries").On("feedentries.cat = latest.cat").And("feedentries.feedid = plants.feedid").
		Join("feedmedias").On("feedmedias.feedentryid = feedentries.id").
		Where("plants.is_public = ?", true).
		And("plants.deleted = ?", false).
		OrderBy("latest.cat desc").
		Offset(offset).Limit(limit)

	results := []publicListingPlantResult{}
	if err := selector.All(&results); err != nil {
		logrus.Errorf("selector.All in fetchPublicPlants %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, p := range results {
		err = loadFeedMediaPublicURLs(&p)
		if err != nil {
			logrus.Errorf("loadFeedMediaPublicURLs in fetchPublicPlants %q - p: %+v", err, p)
			continue
		}
		results[i] = p
	}

	if err := json.NewEncoder(w).Encode(publicPlantsResult{results}); err != nil {
		logrus.Errorf("json.NewEncoder in fetchPublicPlants %q - results: %+v", err, results)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type publicPlantResult struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	FilePath      string `json:"filePath"`
	ThumbnailPath string `json:"thumbnailPath"`
	Settings      string `json:"settings"`
	BoxSettings   string `json:"boxSettings"`
}

func fetchPublicPlant(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)

	plant := sgldb.Plant{}
	if err := sess.Select("*").From("plants").Where("is_public = ?", true).And("deleted = ?", false).And("id = ?", p.ByName("id")).One(&plant); err != nil {
		logrus.Errorf("sess.Select('plants') in fetchPublicPlant %q - id: %s", err, p.ByName("id"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fm, err := loadLastFeedMediaForPlant(sess, plant)
	if err != nil {
		logrus.Errorf("loadLastFeedMediaForPlant in fetchPublicPlant %q - plant: %+v", err, plant)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	box := sgldb.Box{}
	if err := sess.Select("*").From("boxes").And("id = ?", plant.BoxID).One(&box); err != nil {
		logrus.Errorf("sess.Select('boxes') in fetchPublicPlant %q - plant: %+v", err, plant)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := publicPlantResult{plant.ID.UUID.String(), plant.Name, fm.FilePath, fm.ThumbnailPath, plant.Settings, box.Settings}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logrus.Errorf("json.NewEncoder in fetchPublicPlant %q - result: %+v", err, result)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type publicFeedEntry struct {
	sgldb.FeedEntry

	Liked      bool `db:"liked" json:"liked"`
	Bookmarked bool `db:"bookmarked" json:"bookmarked"`
	NComments  int  `db:"ncomments" json:"nComments"`
	NLikes     int  `db:"nlikes" json:"nLikes"`
}

type publicFeedEntriesResult struct {
	Entries []publicFeedEntry `json:"entries"`
}

func fetchPublicFeedEntries(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
	uid, userIDExists := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		logrus.Errorf("strconv.Atoi in fetchPublicFeedEntries %q - offset: %s url: %s", err, r.URL.Query().Get("offset"), r.URL.String())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		logrus.Errorf("strconv.Atoi in fetchPublicFeedEntries %q - limit: %s url: %s", err, r.URL.Query().Get("limit"), r.URL.String())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if limit < 0 {
		limit = 0
	} else if limit > 50 {
		limit = 50
	}

	feedEntries := []publicFeedEntry{}
	selector := sess.Select("fe.*").From("feedentries fe")
	if userIDExists {
		selector = selector.Columns(udb.Raw("exists(select * from likes l where l.userid = ? and l.feedentryid = fe.id) as liked", uid)).
			Columns(udb.Raw("exists(select * from bookmarks b where b.userid = ? and b.feedentryid = fe.id) as bookmarked", uid))
	}
	selector = selector.Columns(udb.Raw("(select count(*) from likes l where l.feedentryid = fe.id) as nlikes")).
		Columns(udb.Raw("(select count(*) from comments c where c.feedentryid = fe.id) as ncomments")).
		Join("feeds f").On("fe.feedid = f.id").
		Join("plants p").On("p.feedid = f.id").
		Where("p.is_public = ?", true).
		And("p.id = ?", p.ByName("id")).
		And("fe.etype not in ('FE_TOWELIE_INFO', 'FE_PRODUCTS')").
		And("fe.deleted = ?", false).
		And("p.deleted = ?", false).
		OrderBy("fe.createdat DESC").Offset(offset).Limit(limit)
	if err := selector.All(&feedEntries); err != nil {
		logrus.Errorf("selector.All in fetchPublicFeedEntries %q - limit: %d offset: %d id: %s", err, limit, offset, p.ByName("id"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := publicFeedEntriesResult{feedEntries}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logrus.Errorf("json.NewEncoder in fetchPublicFeedEntries %q - %+v", err, result)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type publicFeedEntryResult struct {
	Entry publicFeedEntry `json:"entry"`
}

func fetchPublicFeedEntry(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
	uid, userIDExists := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)

	feedEntry := publicFeedEntry{}
	selector := sess.Select("fe.*").From("feedentries fe")
	if userIDExists {
		selector = selector.Columns(udb.Raw("exists(select * from likes l where l.userid = ? and l.feedentryid = fe.id) as liked", uid)).
			Columns(udb.Raw("exists(select * from bookmarks b where b.userid = ? and b.feedentryid = fe.id) as bookmarked", uid))
	}
	selector = selector.Columns(udb.Raw("(select count(*) from likes l where l.feedentryid = fe.id) as nlikes")).
		Columns(udb.Raw("(select count(*) from comments c where c.feedentryid = fe.id) as ncomments")).
		Join("feeds f").On("fe.feedid = f.id").
		Join("plants p").On("p.feedid = f.id").
		Where("p.is_public = ?", true).
		And("fe.id = ?", p.ByName("id")).
		And("fe.etype not in ('FE_TOWELIE_INFO', 'FE_PRODUCTS')").
		And("fe.deleted = ?", false).
		And("p.deleted = ?", false)
	if err := selector.One(&feedEntry); err != nil {
		logrus.Errorf("selector.One in fetchPublicFeedEntry %q - id: %s", err, p.ByName("id"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := publicFeedEntryResult{feedEntry}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logrus.Errorf("json.NewEncoder in fetchPublicFeedEntry %q - %+v", err, result)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type publicFeedMediasResult struct {
	Medias []sgldb.FeedMedia `json:"medias"`
}

func fetchPublicFeedMedias(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)

	feedMedias := []sgldb.FeedMedia{}
	selector := sess.Select("fm.*").From("feedmedias fm").
		Join("feedentries fe").On("fm.feedentryid = fe.id").
		Join("feeds f").On("fe.feedid = f.id").
		Join("plants p").On("p.feedid = f.id").
		Where("p.is_public = ?", true).
		And("fe.id = ?", p.ByName("id")).
		And("fm.deleted = ?", false)
	if err := selector.All(&feedMedias); err != nil {
		logrus.Errorf("selector.All in fetchPublicFeedMedias %q - id: %s", err, p.ByName("id"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var err error
	for i, fm := range feedMedias {
		err = loadFeedMediaPublicURLs(&fm)
		if err != nil {
			logrus.Errorf("loadFeedMediaPublicURLs in fetchPublicFeedMedias %q - %+v", err, fm)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// might not be useful anymore
		feedMedias[i] = fm
	}

	result := publicFeedMediasResult{feedMedias}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logrus.Errorf("json.NewEncoder in fetchPublicFeedMedias %q - %+v", err, result)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func fetchPublicFeedMedia(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)

	feedMedia := sgldb.FeedMedia{}
	selector := sess.Select("fm.*").From("feedmedias fm").
		Join("feedentries fe").On("fm.feedentryid = fe.id").
		Join("feeds f").On("fe.feedid = f.id").
		Join("plants p").On("p.feedid = f.id").
		Where("p.is_public = ?", true).
		And("fm.id = ?", p.ByName("id")).
		And("fm.deleted = ?", false)
	if err := selector.One(&feedMedia); err != nil {
		logrus.Errorf("selector.One in fetchPublicFeedMedia %q - id: %s", err, p.ByName("id"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var err error
	err = loadFeedMediaPublicURLs(&feedMedia)
	if err != nil {
		logrus.Errorf("loadFeedMediaPublicURLs in fetchPublicFeedMedia %q - %+v", err, feedMedia)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(feedMedia); err != nil {
		logrus.Errorf("json.NewEncoder in fetchPublicFeedMedia %q - %+v", err, feedMedia)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
