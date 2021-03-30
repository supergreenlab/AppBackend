/*
 * Copyright (C) 2020  SuperGreenLab <towelie@supergreenlab.com>
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

package middlewares

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

type Endpoint struct {
	Middlewares []middleware.Middleware
	Output      httprouter.Handle
}

func (e Endpoint) Handle() httprouter.Handle {
	s := middleware.NewStack()

	for _, m := range e.Middlewares {
		s.Use(m)
	}
	return s.Wrap(e.Output)
}

func NewEndpoint() Endpoint {
	return Endpoint{Middlewares: []middleware.Middleware{}}
}

type Factory func() interface{}

type DBEndpointBuilder struct {
	Pre  []middleware.Middleware
	DBFn middleware.Middleware
	Post []middleware.Middleware
}

func (dbe *DBEndpointBuilder) AddPre(pre middleware.Middleware) {
	dbe.Pre = append(dbe.Pre, pre)
}

func (dbe DBEndpointBuilder) Endpoint() Endpoint {
	e := NewEndpoint()
	e.Middlewares = append(e.Middlewares, dbe.Pre...)
	e.Middlewares = append(e.Middlewares, dbe.DBFn)
	e.Middlewares = append(e.Middlewares, dbe.Post...)
	return e
}

func NewDBEndpointBuilder(param Factory, input Factory, pre, post []middleware.Middleware, dbfn middleware.Middleware) DBEndpointBuilder {
	e := DBEndpointBuilder{DBFn: dbfn}
	e.Pre = []middleware.Middleware{}
	if param != nil {
		e.Pre = append(e.Pre, DecodeQuery(param))
	}
	if input != nil {
		e.Pre = append(e.Pre, DecodeJSON(input))
	}
	if pre != nil {
		e.Pre = append(e.Pre, pre...)
	}
	e.Post = []middleware.Middleware{}
	if post != nil {
		e.Post = append(e.Post, post...)
	}
	return e
}

type InsertEndpointBuilder struct {
	DBEndpointBuilder

	Collection string
}

func (dbe InsertEndpointBuilder) Endpoint() Endpoint {
	e := dbe.DBEndpointBuilder.Endpoint()
	e.Middlewares = append(e.Middlewares, PublishInsert(dbe.Collection))
	e.Output = OutputObjectID
	return e
}

func NewInsertEndpointBuilder(collection string, input Factory, pre, post []middleware.Middleware) InsertEndpointBuilder {
	e := InsertEndpointBuilder{
		DBEndpointBuilder: NewDBEndpointBuilder(nil, input, pre, post, InsertObject(collection)),
		Collection:        collection,
	}
	return e
}

type UpdateEndpointBuilder struct {
	DBEndpointBuilder

	Collection string
}

func (dbe UpdateEndpointBuilder) Endpoint() Endpoint {
	e := dbe.DBEndpointBuilder.Endpoint()
	e.Output = OutputOK
	return e
}

func NewUpdateEndpointBuilder(collection string, input Factory, pre, post []middleware.Middleware) UpdateEndpointBuilder {
	e := UpdateEndpointBuilder{
		DBEndpointBuilder: NewDBEndpointBuilder(nil, input, pre, post, UpdateObject(collection)),
		Collection:        collection,
	}
	return e
}

type SelectEndpointBuilder struct {
	DBEndpointBuilder

	Selector middleware.Middleware

	Collection string
}

func (dbe SelectEndpointBuilder) Endpoint() Endpoint {
	dbe.Pre[1] = dbe.Selector
	e := dbe.DBEndpointBuilder.Endpoint()
	e.Output = OutputSelectResult(dbe.Collection)
	return e
}

func NewSelectEndpointBuilder(collection string, param, factory Factory, pre, post []middleware.Middleware) SelectEndpointBuilder {
	defaultSelector := func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(SessContextKey{}).(sqlbuilder.Database)
			params := r.Context().Value(QueryObjectContextKey{}).(SelectParams)
			selector := sess.Select("t.id as objectid", db.Raw("t.*")).From(collection + " t")
			selector = selector.OrderBy("t.cat DESC").Offset(params.GetOffset()).Limit(params.GetLimit())
			ctx := context.WithValue(r.Context(), SelectorContextKey{}, selector)
			fn(w, r.WithContext(ctx), p)
		}
	}
	e := SelectEndpointBuilder{
		DBEndpointBuilder: NewDBEndpointBuilder(param, nil, append([]middleware.Middleware{defaultSelector}, pre...), post, SelectQuery(factory)),
		Collection:        collection,
	}
	e.Selector = defaultSelector
	return e
}

type SelectOneEndpointBuilder struct {
	DBEndpointBuilder

	Collection string
}

func (dbe SelectOneEndpointBuilder) Endpoint() Endpoint {
	e := dbe.DBEndpointBuilder.Endpoint()
	e.Output = OutputSelectOneResult(dbe.Collection)
	return e
}

func NewSelectOneEndpointBuilder(collection string, param, factory Factory, pre, post []middleware.Middleware) SelectOneEndpointBuilder {
	e := SelectOneEndpointBuilder{
		DBEndpointBuilder: NewDBEndpointBuilder(param, nil, pre, post, SelectOneQuery(factory)),
		Collection:        collection,
	}
	return e
}

type CountEndpointBuilder struct {
	DBEndpointBuilder

	Selector middleware.Middleware

	Collection string
}

func (dbe CountEndpointBuilder) Endpoint() Endpoint {
	dbe.Pre[1] = dbe.Selector
	e := dbe.DBEndpointBuilder.Endpoint()
	e.Output = OutputSelectOneResult(dbe.Collection)
	return e
}

func NewCountEndpointBuilder(collection string, param Factory, pre, post []middleware.Middleware) CountEndpointBuilder {
	defaultSelector := func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(SessContextKey{}).(sqlbuilder.Database)
			selector := sess.Select(db.Raw("COUNT(*) AS n")).From(collection + " t")
			ctx := context.WithValue(r.Context(), SelectorContextKey{}, selector)
			fn(w, r.WithContext(ctx), p)
		}
	}
	factory := func() interface{} { return &Count{} }
	e := CountEndpointBuilder{
		DBEndpointBuilder: NewDBEndpointBuilder(param, nil, append([]middleware.Middleware{defaultSelector}, pre...), post, SelectOneQuery(factory)),
		Collection:        collection,
	}
	e.Selector = defaultSelector
	return e
}

// InsertEndpoint - insert an object
func InsertEndpoint(
	collection string,
	factory Factory,
	pre []middleware.Middleware,
	post []middleware.Middleware,
) httprouter.Handle {
	return NewInsertEndpointBuilder(collection, factory, pre, post).Endpoint().Handle()
}

// UpdateEndpoint - updates and object
func UpdateEndpoint(
	collection string,
	factory func() interface{},
	pre []middleware.Middleware,
	post []middleware.Middleware,
) httprouter.Handle {
	return NewUpdateEndpointBuilder(collection, factory, pre, post).Endpoint().Handle()
}

type SelectParams interface {
	GetOffset() int
	GetLimit() int
}

type SelectParamsOffsetLimit struct {
	Offset int
	Limit  int
}

func (p *SelectParamsOffsetLimit) GetOffset() int {
	return p.Offset
}

func (p *SelectParamsOffsetLimit) GetLimit() int {
	return p.Limit
}

// SelectEndpoint - select objects
func SelectEndpoint(
	collection string,
	factory func() interface{},
	paramFactory func() interface{},
	pre []middleware.Middleware,
	post []middleware.Middleware,
) httprouter.Handle {
	return NewSelectEndpointBuilder(collection, paramFactory, factory, pre, post).Endpoint().Handle()
}

// SelectEndpoint - select objects
func SelectOneEndpoint(
	collection string,
	factory func() interface{},
	paramFactory func() interface{},
	pre []middleware.Middleware,
	post []middleware.Middleware,
) httprouter.Handle {
	return NewSelectOneEndpointBuilder(collection, paramFactory, factory, pre, post).Endpoint().Handle()
}

type Count struct {
	N int `db:"n" json:"n"`
}

// CountEndpoint - select objects
func CountEndpoint(
	collection string,
	paramFactory func() interface{},
	pre []middleware.Middleware,
	post []middleware.Middleware,
) httprouter.Handle {
	return NewCountEndpointBuilder(collection, paramFactory, pre, post).Endpoint().Handle()
}
