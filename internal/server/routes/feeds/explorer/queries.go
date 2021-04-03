/*
 * Copyright (C) 2021  SuperGreenLab <towelie@supergreenlab.com>
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

package explorer

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"upper.io/db.v3"
	udb "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func pageOffsetLimit(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)
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
		selector = selector.Offset(offset).Limit(limit)

		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func joinLatestPlantFeedMedia(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)

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

		selector = selector.Columns("feedmedias.filepath", "feedmedias.thumbnailpath").
			Join(db.Raw(fmt.Sprintf("(%s) latestfe", lastFeedEntrySelector.String()))).Using("feedid").
			Join(db.Raw(fmt.Sprintf("(%s) latestfm", lastFeedMediaSelector.String()))).Using("feedid").
			Join("feedentries").On("feedentries.cat = latestfe.cat").And("feedentries.feedid = p.feedid").
			Join("feedmedias").On("feedmedias.cat = latestfm.cat").And("latestfm.feedid = p.feedid").
			OrderBy("latestfm.cat desc")

		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func joinBoxSettings(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)

		selector = selector.Columns("boxes.settings as boxsettings").
			Join("boxes").On("boxes.id = p.boxid")

		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func joinFollows(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)
		uid, userIDExists := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)
		if !userIDExists {
			fn(w, r, p)
			return
		}

		selector = selector.Columns(db.Raw("(follows.id is not null) as followed")).
			LeftJoin("follows").On("follows.plantid = p.id and follows.userid = ?", uid)

		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func joinFeedEntrySocialSelector(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)
		uid, userIDExists := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)

		// TODO optimize with joins?
		if userIDExists {
			selector = selector.Columns(udb.Raw("exists(select * from likes l where l.userid = ? and l.feedentryid = fe.id) as liked", uid)).
				Columns(udb.Raw("exists(select * from bookmarks b where b.userid = ? and b.feedentryid = fe.id) as bookmarked", uid))
		}

		selector = selector.Columns(udb.Raw("(select count(*) from likes l where l.feedentryid = fe.id) as nlikes")).
			Columns(udb.Raw("(select count(*) from comments c where c.feedentryid = fe.id) as ncomments"))

		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func publicPlantsOnly(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)

		selector = selector.Where("p.is_public = true").
			And("p.deleted = false")

		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func followedPlantsOnly(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)

		selector = selector.Where("follows.id is not null")

		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func publicFeedEntriesOnly(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)

		selector = selector.Join("feeds f").On("fe.feedid = f.id").
			Join("plants p").On("p.feedid = f.id").
			Where("p.is_public = true").
			And("fe.etype not in ('FE_TOWELIE_INFO', 'FE_PRODUCTS')").
			And("fe.deleted = false").
			And("p.deleted = false")

		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func publicFeedMediasOnly(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)

		selector = selector.Join("feedentries fe").On("fm.feedentryid = fe.id").
			Join("feeds f").On("fe.feedid = f.id").
			Join("plants p").On("p.feedid = f.id").
			Where("p.is_public = true").
			And("fm.deleted = false")

		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func joinLatestFeedMediaForFeedEntry(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)

		lastFeedMediaSelector := sess.Select("feedentryid", udb.Raw("max(feedmedias.cat) as cat")).
			From("feedmedias").
			Where("feedmedias.deleted = false").
			GroupBy("feedentryid")

		selector = selector.Columns("feedmedias.filepath", "feedmedias.thumbnailpath").
			Join(db.Raw(fmt.Sprintf("(%s) latestfm", lastFeedMediaSelector.String()))).On("latestfm.feedentryid = fe.id").
			Join("feedmedias").On("feedmedias.cat = latestfm.cat").And("latestfm.feedentryid = fe.id")

		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func joinPlantForFeedEntry(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)

		selector = selector.Columns("p.name", "p.id as plantid").
			Join("plants p").On("p.feedid = fe.feedid")

		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}
