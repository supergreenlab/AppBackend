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
	"net/http"
	"strings"

	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	udb "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

type SearchPlantsParams struct {
	SelectPlantsParams

	Q string
}

var searchPublicPlants = NewSelectPlantsEndpointBuilder([]middleware.Middleware{
	func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)
			params := r.Context().Value(middlewares.QueryObjectContextKey{}).(*SearchPlantsParams)

			searchTests := []udb.Compound{}
			qs := strings.Split(params.Q, " ")
			for _, q := range qs {
				searchTests = append(searchTests,
					udb.Or(udb.Raw("p.name ilike '%' || ? || '%'", q)).
						Or(udb.Raw("p.settings::text ilike '%' || ? || '%'", q)).
						Or(udb.Raw("boxes.settings::text ilike '%' || ? || '%'", q)),
				)
			}

			selector = selector.Where(udb.And(searchTests...))
			ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
			fn(w, r.WithContext(ctx), p)
		}
	},
}).SetParam(func() interface{} { return &SearchPlantsParams{} }).Endpoint().Handle()
