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
	"time"

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

func pageOffsetLimit(r *http.Request, selector sqlbuilder.Selector) sqlbuilder.Selector {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	if limit < 0 {
		limit = 0
	} else if limit > 50 {
		limit = 50
	}
	return selector.Offset(offset).Limit(limit)
}

func joinLatestFeedMedia(sess sqlbuilder.Database, selector sqlbuilder.Selector) sqlbuilder.Selector {
	lastFeedEntrySelector := sess.Select("feedid", udb.Raw("max(cat) as cat")).
		From("feedentries").
		Where("deleted = false").
		And(fmt.Sprintf("etype in ('%s')", strings.Join([]string{"FE_MEDIA", "FE_BENDING", "FE_DEFOLATION", "FE_TRANSPLANT", "FE_FIMMING", "FE_TOPPING", "FE_MEASURE"}, "', '"))).
		GroupBy("feedid")
	lastFeedMediaSelector := sess.Select("feedid", udb.Raw("max(feedmedias.cat) as cat")).
		From("feedmedias").
		Join("feedentries").On("feedentries.id = feedmedias.feedentryid").
		Where("feedmedias.deleted = false").
		GroupBy("feedid")

	return selector.Columns("feedmedias.filepath", "feedmedias.thumbnailpath").
		Join(db.Raw(fmt.Sprintf("(%s) latestfe", lastFeedEntrySelector.String()))).Using("feedid").
		Join(db.Raw(fmt.Sprintf("(%s) latestfm", lastFeedMediaSelector.String()))).Using("feedid").
		Join("feedentries").On("feedentries.cat = latestfe.cat").And("feedentries.feedid = plants.feedid").
		Join("feedmedias").On("feedmedias.cat = latestfm.cat").And("latestfm.feedid = plants.feedid")
}

func joinBoxSettings(selector sqlbuilder.Selector) sqlbuilder.Selector {
	return selector.Columns("boxes.settings as boxsettings").
		Join("boxes").On("boxes.id = plants.boxid")
}

func joinFollows(r *http.Request, selector sqlbuilder.Selector) sqlbuilder.Selector {
	uid, userIDExists := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)
	if !userIDExists {
		return selector
	}
	return selector.Columns(db.Raw("(follows.id is not null) as followed")).
		Join("follows").On("follows.plantid = plants.id and follows.userid = ?", uid)
}

type publicPlantResult struct {
	ID            string `db:"id" json:"id"`
	Name          string `db:"name" json:"name"`
	FilePath      string `db:"filepath" json:"filePath"`
	ThumbnailPath string `db:"thumbnailpath" json:"thumbnailPath"`

	Followed bool `db:"followed" json:"followed"`

	Settings    string `db:"settings" json:"settings"`
	BoxSettings string `db:"boxsettings" json:"boxSettings"`
}

func (r *publicPlantResult) SetURLs(filePath string, thumbnailPath string) {
	r.FilePath = filePath
	r.ThumbnailPath = thumbnailPath
}

func (r publicPlantResult) GetURLs() (filePath string, thumbnailPath string) {
	filePath, thumbnailPath = r.FilePath, r.ThumbnailPath
	return
}

type publicPlantsResult struct {
	Plants []publicPlantResult `json:"plants"`
}

func fetchPublicPlants(makeSelector func(sess sqlbuilder.Database, w http.ResponseWriter, r *http.Request, p httprouter.Params) sqlbuilder.Selector) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)

		selector := makeSelector(sess, w, r, p).
			Where("plants.is_public = ?", true).
			And("plants.deleted = ?", false)

		selector = joinLatestFeedMedia(sess, selector)
		selector = joinBoxSettings(selector)
		selector = joinFollows(r, selector)
		selector = pageOffsetLimit(r, selector)

		results := []publicPlantResult{}
		if err := selector.All(&results); err != nil {
			logrus.Errorf("selector.All in fetchPublicPlants %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for i, p := range results {
			err := loadFeedMediaPublicURLs(&p)
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
}

var fetchLatestUpdatedPublicPlants = fetchPublicPlants(func(sess sqlbuilder.Database, w http.ResponseWriter, r *http.Request, p httprouter.Params) sqlbuilder.Selector {
	return sess.Select("plants.id", "plants.name", "plants.settings").
		From("plants").
		OrderBy("latestfm.cat desc")
})

var fetchLatestUpdatedFollowedPublicPlants = fetchPublicPlants(func(sess sqlbuilder.Database, w http.ResponseWriter, r *http.Request, p httprouter.Params) sqlbuilder.Selector {
	uid := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)
	return sess.Select("plants.id", "plants.name", "plants.settings").
		From("plants").
		Join("follows").On("follows.plantid = plants.id and userid = ?", uid).
		OrderBy("latestfm.cat desc")
})

func fetchPublicPlant(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)

	plant := publicPlantResult{}
	selector := sess.Select("plants.id", "plants.name", "plants.settings").
		From("plants").
		Where("plants.is_public = ?", true).
		And("plants.deleted = ?", false).
		And("plants.id = ?", p.ByName("id"))

	selector = joinLatestFeedMedia(sess, selector)
	selector = joinBoxSettings(selector)

	if err := selector.One(&plant); err != nil {
		logrus.Errorf("sess.Select('plants') in fetchPublicPlant %q - id: %s", err, p.ByName("id"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err := loadFeedMediaPublicURLs(&plant)
	if err != nil {
		logrus.Errorf("loadFeedMediaPublicURLs in fetchPublicPlant %q - plant: %+v", err, plant)
	}

	if err := json.NewEncoder(w).Encode(plant); err != nil {
		logrus.Errorf("json.NewEncoder in fetchPublicPlant %q - plant: %+v", err, plant)
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

	// TODO make an interface based middleware to unify the select* middlewares
	PlantID       uuid.NullUUID `db:"plantid,omitempty" json:"plantID,omitempty"`
	PlantName     string        `db:"name,omitempty" json:"plantName,omitempty"`
	CommentID     uuid.NullUUID `db:"commentid,omitempty" json:"commentID,omitempty"`
	Comment       *string       `db:"comment,omitempty" json:"comment,omitempty"`
	LikeDate      *time.Time    `db:"likecat,omitempty" json:"likeDate,omitempty"`
	ThumbnailPath *string       `db:"thumbnailpath,omitempty" json:"thumbnailpath,omitempty"`
}

func (r *publicFeedEntry) SetURLs(_, thumbnailPath string) {
	if thumbnailPath != "" {
		*r.ThumbnailPath = thumbnailPath
	}
}

func (r publicFeedEntry) GetURLs() (filePath, thumbnailPath string) {
	filePath, thumbnailPath = "", ""
	if r.ThumbnailPath != nil {
		thumbnailPath = *r.ThumbnailPath
	}
	return
}

type publicFeedEntriesResult struct {
	Entries []publicFeedEntry `json:"entries"`
}

func joinFeedEntrySocialSelector(r *http.Request, selector sqlbuilder.Selector) sqlbuilder.Selector {
	uid, userIDExists := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)

	// TODO optimize with joins?
	if userIDExists {
		selector = selector.Columns(udb.Raw("exists(select * from likes l where l.userid = ? and l.feedentryid = fe.id) as liked", uid)).
			Columns(udb.Raw("exists(select * from bookmarks b where b.userid = ? and b.feedentryid = fe.id) as bookmarked", uid))
	}

	return selector.Columns(udb.Raw("(select count(*) from likes l where l.feedentryid = fe.id) as nlikes")).
		Columns(udb.Raw("(select count(*) from comments c where c.feedentryid = fe.id) as ncomments"))
}

func publicFeedEntriesOnly(selector sqlbuilder.Selector) sqlbuilder.Selector {
	return selector.Join("feeds f").On("fe.feedid = f.id").
		Join("plants p").On("p.feedid = f.id").
		Where("p.is_public = ?", true).
		And("fe.etype not in ('FE_TOWELIE_INFO', 'FE_PRODUCTS')").
		And("fe.deleted = ?", false).
		And("p.deleted = ?", false)
}

func joinLatestFeedMediaForFeedEntry(sess sqlbuilder.Database, selector sqlbuilder.Selector) sqlbuilder.Selector {
	lastFeedMediaSelector := sess.Select("feedentryid", udb.Raw("max(feedmedias.cat) as cat")).
		From("feedmedias").
		Where("feedmedias.deleted = false").
		GroupBy("feedentryid")

	return selector.Columns("feedmedias.filepath", "feedmedias.thumbnailpath").
		Join(db.Raw(fmt.Sprintf("(%s) latestfm", lastFeedMediaSelector.String()))).On("latestfm.feedentryid = fe.id").
		Join("feedmedias").On("feedmedias.cat = latestfm.cat").And("latestfm.feedentryid = fe.id")
}

func joinPlantForFeedEntry(selector sqlbuilder.Selector) sqlbuilder.Selector {
	return selector.Columns("plants.name", "plants.id as plantid").
		Join("plants").On("plants.feedid = fe.feedid")
}

func fetchPublicFeedEntries(makeSelector func(sess sqlbuilder.Database, w http.ResponseWriter, r *http.Request, p httprouter.Params) sqlbuilder.Selector) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)

		selector := makeSelector(sess, w, r, p)

		selector = joinFeedEntrySocialSelector(r, selector)
		selector = publicFeedEntriesOnly(selector)
		selector = pageOffsetLimit(r, selector)

		feedEntries := []publicFeedEntry{}
		if err := selector.All(&feedEntries); err != nil {
			logrus.Errorf("selector.All in fetchPublicFeedEntries %q - id: %s", err, p.ByName("id"))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for i, p := range feedEntries {
			err := loadFeedMediaPublicURLs(&p)
			if err != nil {
				logrus.Errorf("loadFeedMediaPublicURLs in fetchPublicFeedEntries %q - p: %+v", err, p)
				continue
			}
			feedEntries[i] = p
		}

		result := publicFeedEntriesResult{feedEntries}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			logrus.Errorf("json.NewEncoder in fetchPublicFeedEntries %q - %+v", err, result)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

var fetchPublicPlantFeedEntries = fetchPublicFeedEntries(func(sess sqlbuilder.Database, w http.ResponseWriter, r *http.Request, p httprouter.Params) sqlbuilder.Selector {
	return sess.Select("fe.*").From("feedentries fe").
		Where("p.id = ?", p.ByName("id")).
		OrderBy("fe.createdat DESC")
})

var fetchLatestCommentedFeedEntries = fetchPublicFeedEntries(func(sess sqlbuilder.Database, w http.ResponseWriter, r *http.Request, p httprouter.Params) sqlbuilder.Selector {
	selector := sess.Select("fe.*", "comments.text as comment").From("feedentries fe").
		Join("comments").On("comments.feedentryid = fe.id").
		OrderBy("comments.cat DESC")
	selector = joinLatestFeedMediaForFeedEntry(sess, selector)
	selector = joinPlantForFeedEntry(selector)
	return selector
})

var fetchLatestLikedFeedEntries = fetchPublicFeedEntries(func(sess sqlbuilder.Database, w http.ResponseWriter, r *http.Request, p httprouter.Params) sqlbuilder.Selector {
	selector := sess.Select("fe.*", "comments.text as comment", "comments.id as commentid").From("likes").
		LeftJoin("comments").On("comments.id = likes.commentid").
		Join("feedentries fe").On("fe.id = likes.feedentryid or fe.id = comments.feedentryid").
		OrderBy("likes.cat DESC")
	selector = joinLatestFeedMediaForFeedEntry(sess, selector)
	selector = joinPlantForFeedEntry(selector)
	logrus.Info(selector.String())
	return selector
})

type publicFeedEntryResult struct {
	Entry publicFeedEntry `json:"entry"`
}

func fetchPublicFeedEntry(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)

	feedEntry := publicFeedEntry{}
	selector := sess.Select("fe.*").From("feedentries fe").
		Where("fe.id = ?", p.ByName("id"))

	selector = joinFeedEntrySocialSelector(r, selector)
	selector = publicFeedEntriesOnly(selector)

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
