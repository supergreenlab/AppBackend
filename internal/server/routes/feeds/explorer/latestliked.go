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
	udb "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var fetchLatestLikedFeedEntries = NewSelectFeedEntriesEndpointBuilderWithSelector(
	middlewares.Selector(func(sess sqlbuilder.Database) sqlbuilder.Selector {
		commentLikes := sess.Select("fe.*", "c.text as comment", "c.id as commentid", "likes.cat as likecat", "likes.userid as likeuserid").From("likes").
			Join("comments c").On("c.id = likes.commentid").
			Join("feedentries fe").On("fe.id = c.feedentryid")
		entryLikes := sess.Select("fe.*", db.Raw("null as comment"), db.Raw("null as commentid"), "likes.cat as likecat", "likes.userid as likeuserid").From("likes").
			Join("feedentries fe").On("fe.id = likes.feedentryid")
		return sess.Select("fe.*", udb.Raw("users.nickname as nickname"), udb.Raw("users.pic as pic")).
			From(db.Raw(fmt.Sprintf("(%s union %s) fe", commentLikes.String(), entryLikes.String()))).
			Join("users").On("users.id = fe.likeuserid").
			OrderBy("fe.likecat desc")
	}),
	[]middleware.Middleware{
		joinPlantForFeedEntry,
		createJoinLatestPlantFeedMedia(false, false, []interface{}{"latestfmrow.thumbnailpath as plantthumbnailpath"}),
		leftJoinLatestFeedMediaForFeedEntry,
	},
).EnableCache("latestLikedFeedEntries").Endpoint().Handle()
