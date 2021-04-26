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

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

// SessContextKey - context key which stores the DB session object
type SessContextKey struct{}

// CreateDBSession - Creates a DB session and stores it in the context
func CreateDBSession(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		/*sess, err := postgresql.Open(db.Settings)
		if err != nil {
			logrus.Errorf("db.Open(): %q\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer sess.Close()*/

		ctx := context.WithValue(r.Context(), SessContextKey{}, db.Sess)
		fn(w, r.WithContext(ctx), p)
	}
}

// InsertedIDContextKey - context key which stores the inserted object's ID
type InsertedIDContextKey struct{}

// InsertObject - Insert the payload object to DB
func InsertObject(collection string) middleware.Middleware {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			o := r.Context().Value(ObjectContextKey{})
			sess := r.Context().Value(SessContextKey{}).(sqlbuilder.Database)
			col := sess.Collection(collection)
			id, err := col.Insert(o)
			if err != nil {
				logrus.Errorf("Insert in InsertObject %q - %s %+v", err, collection, o)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), InsertedIDContextKey{}, uuid.FromStringOrNil(string(id.([]uint8))))
			fn(w, r.WithContext(ctx), p)
		}
	}
}

// UpdatedIDContextKey - context key which stores the updated object's ID
type UpdatedIDContextKey struct{}

// UpdateObject - Updates the db object with JSON payload object
func UpdateObject(collection string) func(fn httprouter.Handle) httprouter.Handle {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			o := r.Context().Value(ObjectContextKey{}).(appbackend.Object)
			sess := r.Context().Value(SessContextKey{}).(sqlbuilder.Database)
			col := sess.Collection(collection)
			err := col.Find(o.GetID()).Update(o)
			if err != nil {
				logrus.Errorf("Find in UpdateObject %q - %s %+v", err, collection, o)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), UpdatedIDContextKey{}, o.GetID().UUID)
			fn(w, r.WithContext(ctx), p)
		}
	}
}

type SelectorContextKey struct{}
type SelectResultContextKey struct{}

func SelectQuery(factory func() interface{}) func(fn httprouter.Handle) httprouter.Handle {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			selector := r.Context().Value(SelectorContextKey{}).(sqlbuilder.Selector)
			results := factory()
			if err := selector.All(results); err != nil {
				logrus.Errorf("All in SelectQuery %q", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), SelectResultContextKey{}, results)
			fn(w, r.WithContext(ctx), p)
		}
	}
}

func SelectOneQuery(factory func() interface{}) func(fn httprouter.Handle) httprouter.Handle {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			selector := r.Context().Value(SelectorContextKey{}).(sqlbuilder.Selector)
			results := factory()
			if err := selector.One(results); err != nil {
				logrus.Errorf("One in SelectOneQuery %q", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), SelectResultContextKey{}, results)
			fn(w, r.WithContext(ctx), p)
		}
	}
}

type FilterFn func(p httprouter.Params, selector sqlbuilder.Selector) sqlbuilder.Selector

func Filter(filterFn FilterFn) middleware.Middleware {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			selector := r.Context().Value(SelectorContextKey{}).(sqlbuilder.Selector)
			selector = filterFn(p, selector)
			ctx := context.WithValue(r.Context(), SelectorContextKey{}, selector)
			fn(w, r.WithContext(ctx), p)
		}
	}
}

type SelectorFn func(sqlbuilder.Database) sqlbuilder.Selector

func Selector(selectorFn SelectorFn) middleware.Middleware {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(SessContextKey{}).(sqlbuilder.Database)
			selector := selectorFn(sess)
			ctx := context.WithValue(r.Context(), SelectorContextKey{}, selector)
			fn(w, r.WithContext(ctx), p)
		}
	}
}
