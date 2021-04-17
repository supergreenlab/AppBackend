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
	"net/http"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	sgldb "github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/data/storage"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
	udb "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// TODO add deleted filtering

func filterUserID(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)
		uid := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)
		selector = selector.Where("t.userid = ?", uid)
		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

type SelectPlantsParams struct {
	middlewares.SelectParamsOffsetLimit
}

var selectPlants = middlewares.SelectEndpoint(
	"plants",
	func() interface{} { return &[]db.Plant{} },
	func() interface{} { return &SelectPlantsParams{} },
	[]middleware.Middleware{
		filterUserID,
	},
	[]middleware.Middleware{},
)

type SelectFeedEntriesParams struct {
	middlewares.SelectParamsOffsetLimit
}

var selectFeedEntries = middlewares.SelectEndpoint(
	"feedentries",
	func() interface{} { return &[]db.FeedEntry{} },
	func() interface{} { return &SelectFeedEntriesParams{} },
	[]middleware.Middleware{
		filterUserID,
	},
	[]middleware.Middleware{},
)

type SelectFeedsParams struct {
	middlewares.SelectParamsOffsetLimit
}

var selectFeeds = middlewares.SelectEndpoint(
	"feeds",
	func() interface{} { return &[]db.FeedEntry{} },
	func() interface{} { return &SelectFeedsParams{} },
	[]middleware.Middleware{
		filterUserID,
	},
	[]middleware.Middleware{},
)

type SelectBoxesParams struct {
	middlewares.SelectParamsOffsetLimit
}

var selectBoxes = middlewares.SelectEndpoint(
	"boxes",
	func() interface{} { return &[]db.FeedEntry{} },
	func() interface{} { return &SelectBoxesParams{} },
	[]middleware.Middleware{
		filterUserID,
	},
	[]middleware.Middleware{},
)

type SelectDevicesParams struct {
	middlewares.SelectParamsOffsetLimit
}

var selectDevices = middlewares.SelectEndpoint(
	"devices",
	func() interface{} { return &[]db.Device{} },
	func() interface{} { return &SelectDevicesParams{} },
	[]middleware.Middleware{
		filterUserID,
	},
	[]middleware.Middleware{},
)

type SelectFeedMediasParams struct {
	middlewares.SelectParamsOffsetLimit
}

var selectFeedMedias = middlewares.SelectEndpoint(
	"feedmedias",
	func() interface{} { return &[]db.FeedMedia{} },
	func() interface{} { return &SelectFeedMediasParams{} },
	[]middleware.Middleware{
		filterUserID,
	},
	[]middleware.Middleware{},
)

type SelectFeedEntryCommentsParams struct {
	middlewares.SelectParamsOffsetLimit
	ReplyTo          *string `json:"replyTo"`
	RootCommentsOnly bool    `json:"rootCommentsOnly"`
	AllComments      bool    `json:"allComments"`
}

type Comment struct {
	db.Comment
	From string      `db:"nickname" json:"from"`
	Pic  null.String `db:"pic" json:"pic"`

	Liked    bool `db:"liked" json:"liked"`
	NReplies int  `db:"nreplies" json:"nReplies"`
	NLikes   int  `db:"nlikes" json:"nLikes"`
}

func filterFeedEntryID(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)
		params := r.Context().Value(middlewares.QueryObjectContextKey{}).(*SelectFeedEntryCommentsParams)
		feid := p.ByName("id")
		selector = selector.Where("t.feedentryid = ?", feid)
		if params.ReplyTo != nil {
			selector = selector.Where("t.replyto = ?", *(params.ReplyTo))
		} else if !params.AllComments {
			selector = selector.Where("t.replyto is null")
		}
		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func joinCommentSocialSelector(ctx context.Context, selector sqlbuilder.Selector) sqlbuilder.Selector {
	uid, userIDExists := ctx.Value(middlewares.UserIDContextKey{}).(uuid.UUID)
	selector = selector.Columns(udb.Raw("u.nickname"), udb.Raw("u.pic")).Join("users u").On("t.userid = u.id")

	if userIDExists {
		selector = selector.Columns(udb.Raw("exists(select * from likes l where l.userid = ? and l.commentid = t.id) as liked", uid))
	}
	selector = selector.Columns(udb.Raw("(select count(*) from likes l where l.commentid = t.id) as nlikes")).
		Columns(udb.Raw("(select count(*) from comments c where c.replyto = t.id) as nreplies"))
	return selector
}

func joinCommentSocial(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)
		selector = joinCommentSocialSelector(r.Context(), selector)
		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func picMediaURL(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		result := r.Context().Value(middlewares.SelectResultContextKey{}).(*[]Comment)

		for i, c := range *result {
			if c.Pic.Valid == false {
				continue
			}
			expiry := time.Second * 60 * 60
			url1, err := storage.Client.PresignedGetObject("users", c.Pic.String, expiry, nil)
			if err != nil {
				c.Pic = null.NewString("", false)
				logrus.Errorf("minioClient.PresignedGetObject in picMediaURL %q - %+v", err, c)
			} else {
				c.Pic = null.NewString(url1.RequestURI(), true)
			}
			(*result)[i] = c
		}
		fn(w, r, p)
	}
}

func selectRepliesForComments(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		params := r.Context().Value(middlewares.QueryObjectContextKey{}).(*SelectFeedEntryCommentsParams)

		if params.AllComments || params.RootCommentsOnly || params.ReplyTo != nil {
			fn(w, r, p)
			return
		}

		result := r.Context().Value(middlewares.SelectResultContextKey{}).(*[]Comment)
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)

		ids := []uuid.UUID{}
		for _, c := range *result {
			if c.ReplyTo.Valid == false {
				ids = append(ids, c.ID.UUID)
			}
		}

		replies := &[]Comment{}
		selector := joinCommentSocialSelector(r.Context(), sess.Select("t.*").From("comments t").Where("replyto in ?", ids))
		if err := selector.All(replies); err != nil {
			logrus.Errorf("selector.All in selectRepliesForComments %q - %+v", err, ids)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		*result = append(*result, *replies...)
		ctx := context.WithValue(r.Context(), middlewares.SelectResultContextKey{}, result)
		fn(w, r.WithContext(ctx), p)
	}
}

var selectFeedEntryComments = middlewares.SelectEndpoint(
	"comments",
	func() interface{} { return &[]Comment{} },
	func() interface{} { return &SelectFeedEntryCommentsParams{} },
	[]middleware.Middleware{
		filterFeedEntryID,
		joinCommentSocial,
	},
	[]middleware.Middleware{
		selectRepliesForComments,
		picMediaURL,
	},
)

func filterCommentID(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)
		cid := p.ByName("id")
		selector = selector.Where("t.id = ?", cid)
		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

var selectComment = middlewares.SelectEndpoint(
	"comments",
	func() interface{} { return &[]Comment{} },
	func() interface{} { return &SelectFeedEntryCommentsParams{} },
	[]middleware.Middleware{
		filterCommentID,
		joinCommentSocial,
	},
	[]middleware.Middleware{
		selectRepliesForComments,
		picMediaURL,
	},
)

var countFeedEntryComments = middlewares.CountEndpoint(
	"comments",
	func() interface{} { return &SelectFeedEntryCommentsParams{} },
	[]middleware.Middleware{
		filterFeedEntryID,
	},
	[]middleware.Middleware{},
)

type SelectFeedEntrySocialParams struct{}

type FeedEntrySocial struct {
	Liked      bool `db:"liked" json:"liked"`
	Bookmarked bool `db:"bookmarked" json:"bookmarked"`
	NLikes     int  `db:"nlikes" json:"nLikes"`
	NComments  int  `db:"ncomments" json:"nComments"`
}

func feedEntrySocialSelect(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
		uid, userIDExists := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)

		feid := p.ByName("id")

		selector := sess.Select()
		// TODO DRY this with explorer middleware
		if userIDExists {
			selector = selector.Columns(udb.Raw("exists(select * from likes l where l.userid = ? and l.feedentryid = ?) as liked", uid, feid)).
				Columns(udb.Raw("exists(select * from bookmarks b where b.userid = ? and b.feedentryid = ?) as bookmarked", uid, feid))
		}
		selector = selector.Columns(udb.Raw("(select count(*) from likes l where l.feedentryid = ?) as nlikes", feid)).
			Columns(udb.Raw("(select count(*) from comments c where c.feedentryid = ?) as ncomments", feid))

		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

var selectFeedEntrySocial = middlewares.SelectOneEndpoint(
	"comments",
	func() interface{} { return &FeedEntrySocial{} },
	func() interface{} { return &SelectFeedEntrySocialParams{} },
	[]middleware.Middleware{
		feedEntrySocialSelect,
	},
	[]middleware.Middleware{},
)

type SelectBookmarksParams struct {
	middlewares.SelectParamsOffsetLimit
}

func joinFeedEntry(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)
		uid, userIDExists := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)

		selector = selector.Columns("fe.*", "p.settings as plantsettings").Join("feedentries fe").On("t.feedentryid = fe.id").
			Columns("p.id as plantid").Join("plants p").On("p.feedid = fe.feedid").Where("p.deleted = ?", false).And("p.is_public = ?", true)

		if userIDExists {
			selector = selector.Columns(udb.Raw("exists(select * from likes l where l.userid = ? and l.feedentryid = fe.id) as liked", uid)).
				Columns(udb.Raw("exists(select * from bookmarks b where b.userid = ? and b.feedentryid = fe.id) as bookmarked", uid))
		}
		selector = selector.Columns(udb.Raw("(select count(*) from likes l where l.feedentryid = fe.id) as nlikes")).
			Columns(udb.Raw("(select count(*) from comments c where c.feedentryid = fe.id) as ncomments")).
			OrderBy("fe.createdat DESC")

		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

type publicFeedEntryBookmark struct {
	sgldb.FeedEntry

	Liked      bool `db:"liked" json:"liked"`
	Bookmarked bool `db:"bookmarked" json:"bookmarked"`
	NComments  int  `db:"ncomments" json:"nComments"`
	NLikes     int  `db:"nlikes" json:"nLikes"`

	// TODO this will be a problem with box entries
	PlantID       string `db:"plantid" json:"plantID"`
	PlantSettings string `db:"plantsettings" json:"plantSettings"`
}

var selectBookmarks = middlewares.SelectEndpoint(
	"bookmarks",
	func() interface{} { return &[]publicFeedEntryBookmark{} },
	func() interface{} { return &SelectBookmarksParams{} },
	[]middleware.Middleware{
		filterUserID,
		joinFeedEntry,
	},
	[]middleware.Middleware{},
)

type SelectTimelapsesParams struct {
	middlewares.SelectParamsOffsetLimit
}

var selectTimelapses = middlewares.SelectEndpoint(
	"timelapses",
	func() interface{} { return &[]db.Timelapse{} },
	func() interface{} { return &SelectTimelapsesParams{} },
	[]middleware.Middleware{
		filterUserID,
	},
	[]middleware.Middleware{},
)
