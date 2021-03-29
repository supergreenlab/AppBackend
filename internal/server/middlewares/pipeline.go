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
	return s.Wrap(OutputObjectID)
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

func NewDBEndpointBuilder(param Factory, input Factory) DBEndpointBuilder {
	e := DBEndpointBuilder{}
	e.Pre = []middleware.Middleware{
		DecodeQuery(param),
		DecodeJSON(input),
	}
	e.Post = []middleware.Middleware{}
	return e
}

func (dbe DBEndpointBuilder) Endpoint() Endpoint {
	e := NewEndpoint()
	e.Middlewares = append(e.Middlewares, dbe.Pre...)
	e.Middlewares = append(e.Middlewares, dbe.DBFn)
	e.Middlewares = append(e.Middlewares, dbe.Post...)
	return e
}

type InsertEndpointBuilder struct {
	DBEndpointBuilder

	Collection string
}

func NewInsertEndpointBuilder(collection string, param Factory, input Factory) InsertEndpointBuilder {
	e := InsertEndpointBuilder{
		DBEndpointBuilder: NewDBEndpointBuilder(param, input),
		Collection:        collection,
	}
	e.DBFn = InsertObject(collection)
	return e
}

func (dbe InsertEndpointBuilder) Endpoint() Endpoint {
	e := dbe.DBEndpointBuilder.Endpoint()
	e.Middlewares = append(e.Middlewares, PublishInsert(dbe.Collection))
	e.Output = OutputObjectID
	return e
}

// InsertEndpoint - insert an object
func InsertEndpoint(
	collection string,
	factory func() interface{},
	pre []middleware.Middleware,
	post []middleware.Middleware,
) httprouter.Handle {
	s := middleware.NewStack()

	s.Use(DecodeJSON(factory))
	if pre != nil {
		for _, m := range pre {
			s.Use(m)
		}
	}
	s.Use(InsertObject(collection))

	if post != nil {
		for _, m := range post {
			s.Use(m)
		}
	}
	s.Use(PublishInsert(collection))

	return s.Wrap(OutputObjectID)
}

// UpdateEndpoint - updates and object
func UpdateEndpoint(
	collection string,
	factory func() interface{},
	pre []middleware.Middleware,
	post []middleware.Middleware,
) httprouter.Handle {
	s := middleware.NewStack()

	s.Use(DecodeJSON(factory))
	if pre != nil {
		for _, m := range pre {
			s.Use(m)
		}
	}
	s.Use(UpdateObject(collection))

	if post != nil {
		for _, m := range post {
			s.Use(m)
		}
	}

	return s.Wrap(OutputOK)
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
	s := middleware.NewStack()

	s.Use(DecodeQuery(paramFactory))

	s.Use(func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(SessContextKey{}).(sqlbuilder.Database)
			params := r.Context().Value(QueryObjectContextKey{}).(SelectParams)
			selector := sess.Select("t.id as objectid", db.Raw("t.*")).From(collection + " t")
			selector = selector.OrderBy("t.cat DESC").Offset(params.GetOffset()).Limit(params.GetLimit())
			ctx := context.WithValue(r.Context(), SelectorContextKey{}, selector)
			fn(w, r.WithContext(ctx), p)
		}
	})

	if pre != nil {
		for _, m := range pre {
			s.Use(m)
		}
	}
	s.Use(SelectQuery(factory))

	if post != nil {
		for _, m := range post {
			s.Use(m)
		}
	}

	return s.Wrap(OutputSelectResult(collection))
}

// SelectEndpoint - select objects
func SelectOneEndpoint(
	collection string,
	factory func() interface{},
	paramFactory func() interface{},
	pre []middleware.Middleware,
	post []middleware.Middleware,
) httprouter.Handle {
	s := middleware.NewStack()

	s.Use(DecodeQuery(paramFactory))

	if pre != nil {
		for _, m := range pre {
			s.Use(m)
		}
	}
	s.Use(SelectOneQuery(factory))

	if post != nil {
		for _, m := range post {
			s.Use(m)
		}
	}

	return s.Wrap(OutputSelectOneResult(collection))
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
	s := middleware.NewStack()

	s.Use(DecodeQuery(paramFactory))

	s.Use(func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(SessContextKey{}).(sqlbuilder.Database)
			selector := sess.Select(db.Raw("COUNT(*) AS n")).From(collection + " t")
			ctx := context.WithValue(r.Context(), SelectorContextKey{}, selector)
			fn(w, r.WithContext(ctx), p)
		}
	})

	if pre != nil {
		for _, m := range pre {
			s.Use(m)
		}
	}
	s.Use(SelectOneQuery(func() interface{} { return &Count{} }))

	if post != nil {
		for _, m := range post {
			s.Use(m)
		}
	}

	return s.Wrap(OutputSelectOneResult(collection))
}
