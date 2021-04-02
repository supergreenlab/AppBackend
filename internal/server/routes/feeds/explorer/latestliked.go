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
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

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
