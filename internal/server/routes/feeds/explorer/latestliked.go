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
	"fmt"

	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/rileyr/middleware"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var fetchLatestLikedFeedEntries = NewSelectFeedEntriesEndpointBuilderWithSelector(
	middlewares.Selector(func(sess sqlbuilder.Database) sqlbuilder.Selector {
		commentLikes := sess.Select("fe.*", "comments.text as comment", "comments.id as commentid", "likes.cat as dateLiked").From("likes").
			Join("comments").On("comments.id = likes.commentid").
			Join("feedentries fe").On("fe.id = comments.feedentryid")
		entryLikes := sess.Select("fe.*", db.Raw("'' as comment"), db.Raw("null as commentid"), "likes.cat as dateLiked").From("likes").
			Join("feedentries fe").On("fe.id = likes.feedentryid")
		return sess.Select("*").From(db.Raw(fmt.Sprintf("(%s union %s) fe", commentLikes.String(), entryLikes.String()))).OrderBy("dateLiked desc")
	}),
	[]middleware.Middleware{
		joinLatestFeedMediaForFeedEntry,
		joinPlantForFeedEntry,
	},
).Endpoint().Handle()
