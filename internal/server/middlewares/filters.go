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

package middlewares

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"upper.io/db.v3/lib/sqlbuilder"
)

func FilterUserID(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(SelectorContextKey{}).(sqlbuilder.Selector)
		uid := r.Context().Value(UserIDContextKey{}).(uuid.UUID)
		selector = selector.Where("t.userid = ?", uid)
		ctx := context.WithValue(r.Context(), SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func FilterID(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(SelectorContextKey{}).(sqlbuilder.Selector)

		id := p.ByName("id")
		selector = selector.Where("t.id = ?", id)
		ctx := context.WithValue(r.Context(), SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}
