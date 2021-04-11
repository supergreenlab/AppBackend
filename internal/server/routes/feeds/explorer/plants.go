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

	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"upper.io/db.v3/lib/sqlbuilder"
)

type SelectPlantsParams struct {
	middlewares.SelectParamsOffsetLimit
}

type SelectPlantsEndpointBuilder struct {
	middlewares.DBEndpointBuilder

	Selector middleware.Middleware
}

func (dbe SelectPlantsEndpointBuilder) SetParam(param middlewares.Factory) SelectPlantsEndpointBuilder {
	dbe.Params = middlewares.DecodeQuery(param)
	return dbe
}

func (dbe SelectPlantsEndpointBuilder) Endpoint() middlewares.Endpoint {
	dbe.Pre[0] = dbe.Selector
	e := dbe.DBEndpointBuilder.Endpoint()
	e.Output = dbe.DBEndpointBuilder.Output
	return e
}

func NewSelectPlantsEndpointBuilder(pre []middleware.Middleware) SelectPlantsEndpointBuilder {
	defaultSelector := func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
			selector := sess.Select("p.id", "p.name", "p.settings").From("plants p")
			ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
			fn(w, r.WithContext(ctx), p)
		}
	}
	return NewSelectPlantsEndpointBuilderWithSelector(defaultSelector, pre)
}

func NewSelectPlantsEndpointBuilderWithSelector(selector middleware.Middleware, pre []middleware.Middleware) SelectPlantsEndpointBuilder {
	pre = append([]middleware.Middleware{
		selector,
		publicPlantsOnly,
		joinLatestPlantFeedMedia,
		joinBoxSettings,
		joinFollows,
		pageOffsetLimit,
	}, pre...)
	post := []middleware.Middleware{
		loadFeedMedias,
	}
	factory := func() interface{} { return &publicPlants{} }
	e := SelectPlantsEndpointBuilder{
		DBEndpointBuilder: middlewares.NewDBEndpointBuilder(
			func() interface{} { return &SelectPlantsParams{} }, nil,
			pre, post,
			middlewares.SelectQuery(factory),
			middlewares.OutputResult("plants")),
		Selector: selector,
	}
	return e
}

func NewSelectPlantEndpointBuilder(pre []middleware.Middleware) SelectPlantsEndpointBuilder {
	defaultSelector := func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
			selector := sess.Select("p.id", "p.name", "p.settings").From("plants p")

			ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
			fn(w, r.WithContext(ctx), p)
		}
	}

	pre = append([]middleware.Middleware{
		defaultSelector,
		publicPlantsOnly,
		joinLatestPlantFeedMedia,
		joinBoxSettings,
		joinFollows,
	}, pre...)
	post := []middleware.Middleware{
		loadFeedMedia,
	}
	factory := func() interface{} { return &publicPlant{} }
	e := SelectPlantsEndpointBuilder{
		DBEndpointBuilder: middlewares.NewDBEndpointBuilder(
			nil, nil,
			pre, post,
			middlewares.SelectOneQuery(factory),
			middlewares.OutputSelectOneResult()),
		Selector: defaultSelector,
	}
	return e
}
