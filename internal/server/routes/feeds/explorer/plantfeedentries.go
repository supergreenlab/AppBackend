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

package feeds

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"upper.io/db.v3/lib/sqlbuilder"
)

var fetchPublicPlantFeedEntries = fetchPublicFeedEntries(func(sess sqlbuilder.Database, w http.ResponseWriter, r *http.Request, p httprouter.Params) sqlbuilder.Selector {
	return sess.Select("fe.*").From("feedentries fe").
		Where("p.id = ?", p.ByName("id")).
		OrderBy("fe.createdat DESC")
})
