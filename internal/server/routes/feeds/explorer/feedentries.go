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
	"fmt"
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"upper.io/db.v3/lib/sqlbuilder"
)

type SelectFeedEntriesParams struct {
	middlewares.SelectParamsOffsetLimit
}

type SelectFeedEntriesEndpointBuilder struct {
	middlewares.DBEndpointBuilder

	Cache    middleware.Middleware
	Selector middleware.Middleware
}

func (dbe SelectFeedEntriesEndpointBuilder) Endpoint() middlewares.Endpoint {
	if dbe.Cache != nil {
		dbe.Pre = append([]middleware.Middleware{dbe.Cache}, dbe.Pre...)
		dbe.Pre[1] = dbe.Selector
	} else {
		dbe.Pre[0] = dbe.Selector
	}
	e := dbe.DBEndpointBuilder.Endpoint()
	e.Output = dbe.DBEndpointBuilder.Output
	return e
}

func (dbe SelectFeedEntriesEndpointBuilder) EnableCache(prefix string) SelectFeedEntriesEndpointBuilder {
	dbe.Cache = middlewares.SelectCacheResult(func(r *http.Request, p httprouter.Params) string {
		params := r.Context().Value(middlewares.QueryObjectContextKey{}).(*SelectFeedEntriesParams)
		return fmt.Sprintf("%s.%d-%d", prefix, params.Offset, params.Limit)
	})
	return dbe
}

func (dbe SelectFeedEntriesEndpointBuilder) JoinSocial() SelectFeedEntriesEndpointBuilder {
	dbe.Pre = append(dbe.Pre, joinFeedEntrySocialSelector)
	return dbe
}

func NewSelectFeedEntriesEndpointBuilder(pre []middleware.Middleware) SelectFeedEntriesEndpointBuilder {
	defaultSelector := func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
			params := r.Context().Value(middlewares.QueryObjectContextKey{}).(*SelectFeedEntriesParams)
			selector := sess.Select("fe.*").From("feedentries fe")
			selector = selector.OrderBy("fe.createdat DESC").Offset(params.GetOffset()).Limit(params.GetLimit())
			ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
			fn(w, r.WithContext(ctx), p)
		}
	}
	return NewSelectFeedEntriesEndpointBuilderWithSelector(defaultSelector, pre)
}

func NewSelectFeedEntriesEndpointBuilderWithSelector(selector middleware.Middleware, pre []middleware.Middleware) SelectFeedEntriesEndpointBuilder {
	pre = append([]middleware.Middleware{
		selector,
		publicFeedEntriesOnly,
		pageOffsetLimit,
	}, pre...)
	post := []middleware.Middleware{
		loadFeedMedias,
	}
	factory := func() interface{} { return &publicFeedEntries{} }
	e := SelectFeedEntriesEndpointBuilder{
		DBEndpointBuilder: middlewares.NewDBEndpointBuilder(
			func() interface{} { return &SelectFeedEntriesParams{} }, nil,
			pre, post,
			middlewares.SelectQuery(factory),
			middlewares.OutputResult("entries")),
		Selector: selector,
	}
	return e
}

func NewSelectFeedEntryEndpointBuilder(pre []middleware.Middleware) SelectFeedEntriesEndpointBuilder {
	defaultSelector := func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
			selector := sess.Select("fe.*").From("feedentries fe")
			ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
			fn(w, r.WithContext(ctx), p)
		}
	}

	pre = append([]middleware.Middleware{
		defaultSelector,
		joinFeedEntrySocialSelector,
		publicFeedEntriesOnly,
	}, pre...)
	post := []middleware.Middleware{
		loadFeedMedia,
	}
	factory := func() interface{} { return &publicFeedEntry{} }
	e := SelectFeedEntriesEndpointBuilder{
		DBEndpointBuilder: middlewares.NewDBEndpointBuilder(
			nil, nil,
			pre, post,
			middlewares.SelectOneQuery(factory),
			middlewares.OutputResult("entry")), // TODO: fix inconsistency with other explorer endpoints
		Selector: defaultSelector,
	}
	return e
}
