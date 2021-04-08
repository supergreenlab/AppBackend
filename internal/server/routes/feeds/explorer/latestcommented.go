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
	udb "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var fetchLatestCommentedFeedEntries = NewSelectFeedEntriesEndpointBuilder([]middleware.Middleware{
	middlewares.Filter(func(p httprouter.Params, selector sqlbuilder.Selector) sqlbuilder.Selector {
		return selector.Columns(
			udb.Raw("comments.id as commentid"),
			udb.Raw("comments.text as comment"),
			udb.Raw("comments.ctype as commenttype"),
			udb.Raw("comments.cat as commentdate"),
			udb.Raw("users.nickname as nickname"),
			udb.Raw("users.pic as pic"),
			udb.Raw("pfeo.settings as plantsettings"),
			udb.Raw("boxes.settings as boxsettings")).
			Join("boxes").On("boxes.id = pfeo.boxid").
			Join("comments").On("comments.feedentryid = fe.id").
			Join("users").On("users.id = comments.userid").
			OrderBy("comments.cat DESC")
	}),
	joinLatestFeedMediaForFeedEntry,
	joinPlantForFeedEntry,
}).Endpoint().Handle()
