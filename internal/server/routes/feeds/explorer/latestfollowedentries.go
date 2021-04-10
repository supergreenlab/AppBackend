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
	"strings"

	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"upper.io/db.v3/lib/sqlbuilder"
)

var fetchLatestFollowedFeedEntries = NewSelectFeedEntriesEndpointBuilder([]middleware.Middleware{
	middlewares.Filter(func(p httprouter.Params, selector sqlbuilder.Selector) sqlbuilder.Selector {
		return selector.
			Where(
				fmt.Sprintf("fe.etype in ('%s')",
					strings.Join([]string{"FE_MEDIA", "FE_BENDING", "FE_DEFOLATION", "FE_TRANSPLANT", "FE_FIMMING", "FE_TOPPING", "FE_MEASURE"}, "', '")))
	}),
	joinPlantForFeedEntry,
	joinBoxSettings,
	followedFeedEntriesOnly,
}).JoinSocial().Endpoint().Handle()
