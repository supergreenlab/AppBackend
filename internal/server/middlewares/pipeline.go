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
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
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
