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
	"github.com/SuperGreenLab/AppBackend/internal/server/tools"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

type SelectFeedEntriesParams struct {
	middlewares.SelectParamsOffsetLimit
}

type SelectFeedEntriesEndpointBuilder struct {
	middlewares.DBEndpointBuilder

	Selector middleware.Middleware
}

func (dbe SelectFeedEntriesEndpointBuilder) Endpoint() middlewares.Endpoint {
	dbe.Pre[0] = dbe.Selector
	e := dbe.DBEndpointBuilder.Endpoint()
	e.Output = dbe.DBEndpointBuilder.Output
	return e
}

func NewSelectFeedEntriesEndpointBuilder(pre []middleware.Middleware) SelectFeedEntriesEndpointBuilder {
	defaultSelector := func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
			params := r.Context().Value(middlewares.QueryObjectContextKey{}).(SelectFeedEntriesParams)
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
		joinFeedEntrySocialSelector,
		publicFeedEntriesOnly,
		pageOffsetLimit,
	}, pre...)
	post := []middleware.Middleware{
		func(fn httprouter.Handle) httprouter.Handle {
			return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				feedEntries := r.Context().Value(middlewares.SelectResultContextKey{}).(*[]publicFeedEntry)
				for i, p := range *feedEntries {
					err := tools.LoadFeedMediaPublicURLs(&p)
					if err != nil {
						logrus.Errorf("tools.LoadFeedMediaPublicURLs in fetchPublicFeedEntries %q - p: %+v", err, p)
						continue
					}
					(*feedEntries)[i] = p
				}
				ctx := context.WithValue(r.Context(), middlewares.SelectResultContextKey{}, feedEntries)
				fn(w, r.WithContext(ctx), p)
			}
		},
	}
	factory := func() interface{} { return &[]publicFeedEntry{} }
	e := SelectFeedEntriesEndpointBuilder{
		DBEndpointBuilder: middlewares.NewDBEndpointBuilder(
			func() interface{} { return SelectFeedEntriesParams{} },
			nil,
			pre, post,
			middlewares.SelectQuery(factory),
			middlewares.OutputResult("feedentries")),
	}
	return e
}
