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
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"upper.io/db.v3/lib/sqlbuilder"
)

var fetchLatestCommentedFeedEntries = NewSelectFeedEntriesEndpointBuilder([]middleware.Middleware{
	middlewares.Filter(func(p httprouter.Params, selector sqlbuilder.Selector) sqlbuilder.Selector {
		return selector.Columns(
			"comments.id as commentid",
			"comments.text as comment",
			"comments.ctype as commenttype",
			"comments.cat as commentdate",
			"comments.replyto as commentreplyto",
			"users.nickname as nickname",
			"users.pic as pic",
			"pfeo.settings as plantsettings",
			"boxes.settings as boxsettings").
			Join("boxes").On("boxes.id = pfeo.boxid").
			Join("comments").On("comments.feedentryid = fe.id").
			Join("users").On("users.id = comments.userid").
			OrderBy("comments.cat DESC")
	}),
	joinPlantForFeedEntry,
	createJoinLatestPlantFeedMedia(false, false, []interface{}{"latestfmrow.thumbnailpath as plantthumbnailpath"}),
	leftJoinLatestFeedMediaForFeedEntry,
}).EnableCache("latestCommentedFeedEntries").Endpoint().Handle()
