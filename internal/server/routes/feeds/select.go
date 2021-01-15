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

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	udb "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

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

func filterFeedEntryID(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)
		feid := p.ByName("id")
		selector = selector.Where("t.feedentryid = ?", feid)
		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func joinUser(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)
		selector = selector.Columns(udb.Raw("u.nickname")).Join("users u").On("t.userid = u.id")
		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

type SelectFeedEntryCommentsParams struct {
	middlewares.SelectParamsOffsetLimit
}

type Comment struct {
	db.Comment
	From string `db:"nickname" json:"from"`
}

var selectFeedEntryComments = middlewares.SelectEndpoint(
	"comments",
	func() interface{} { return &[]Comment{} },
	func() interface{} { return &SelectFeedEntryCommentsParams{} },
	[]middleware.Middleware{
		filterFeedEntryID,
		joinUser,
	},
	[]middleware.Middleware{},
)

var countFeedEntryComments = middlewares.CountEndpoint(
	"comments",
	func() interface{} { return &SelectFeedEntryCommentsParams{} },
	[]middleware.Middleware{
		filterFeedEntryID,
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
