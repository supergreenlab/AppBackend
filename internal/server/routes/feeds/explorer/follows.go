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

	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"upper.io/db.v3/lib/sqlbuilder"
)

var fetchLatestUpdatedFollowedPublicPlants = fetchPublicPlants(func(sess sqlbuilder.Database, w http.ResponseWriter, r *http.Request, p httprouter.Params) sqlbuilder.Selector {
	uid := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)
	return sess.Select("plants.id", "plants.name", "plants.settings").
		From("plants").
		Join("follows").On("follows.plantid = plants.id and userid = ?", uid).
		OrderBy("latestfm.cat desc")
})
