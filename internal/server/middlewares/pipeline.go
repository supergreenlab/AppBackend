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
	"upper.io/db.v3/lib/sqlbuilder"
)

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

type SelectParams struct {
	Offset int
	Limit  int
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

	s.Use(func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(SessContextKey{}).(sqlbuilder.Database)
			ctx := context.WithValue(r.Context(), SelectorContextKey{}, sess.Select("*").From(collection))
			fn(w, r.WithContext(ctx), p)
		}
	})

	s.Use(DecodeQuery(paramFactory))
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
