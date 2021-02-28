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
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

func CheckPlantArchivedForPlant(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
		o := r.Context().Value(middlewares.ObjectContextKey{}).(*db.Plant)
		plant := db.Plant{}
		if err := sess.Select("archived").From("plants").Where("id = ?", o.ID).One(&plant); err != nil {
			logrus.Errorf("sess.Select in CheckPlantArchivedForPlant %q - %+v", err, o)
			http.Error(w, "Uknown plant", http.StatusBadRequest)
			return
		}
		if plant.Archived == false {
			fn(w, r, p)
		}
	}
}

func CheckPlantArchivedForTimelapse(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
		o := r.Context().Value(middlewares.ObjectContextKey{}).(*db.Timelapse)
		plant := db.Plant{}
		if err := sess.Select("archived").From("plants").Where("id = ?", o.PlantID).One(&plant); err != nil {
			logrus.Errorf("sess.Select in CheckPlantArchivedForTimelapse %q - %+v", err, o)
			http.Error(w, "Uknown plant", http.StatusBadRequest)
			return
		}
		if plant.Archived == false {
			fn(w, r, p)
		}
	}
}

func CheckPlantArchivedForFeed(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
		o := r.Context().Value(middlewares.ObjectContextKey{}).(*db.Feed)
		plant := db.Plant{}
		if err := sess.Select("archived").From("plants").Where("feedid = ?", o.ID).One(&plant); err != nil && err.Error() != "upper: no more rows in this result set" {
			logrus.Errorf("sess.Select in CheckPlantArchivedForFeed %q - %+v", err, o)
			http.Error(w, "Uknown plant", http.StatusBadRequest)
			return
		}
		if plant.Archived == false {
			fn(w, r, p)
		}
	}
}

func CheckPlantArchivedForFeedEntry(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
		o := r.Context().Value(middlewares.ObjectContextKey{}).(*db.FeedEntry)
		plant := db.Plant{}
		if err := sess.Select("archived").From("plants").Join("feedentries").On("plants.feedid = feedentries.feedid").Where("feedentries.id = ?", o.ID).One(&plant); err != nil && err.Error() != "upper: no more rows in this result set" {
			logrus.Errorf("sess.Select in CheckPlantArchivedForFeedEntry %q - %+v", err, o)
			http.Error(w, "Uknown plant", http.StatusBadRequest)
			return
		}
		if plant.Archived == false {
			fn(w, r, p)
		}
	}
}

func CheckPlantArchivedForFeedMedia(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
		o := r.Context().Value(middlewares.ObjectContextKey{}).(*db.FeedMedia)
		plant := db.Plant{}
		if err := sess.Select("archived").From("plants").Join("feedentries").On("plants.feedid = feedentries.feedid").Join("feedmedias").On("feedmedias.feedentryid = feedentries.id").Where("feedmedias.id = ?", o.ID).One(&plant); err != nil && err.Error() != "upper: no more rows in this result set" {
			logrus.Errorf("sess.Select in CheckPlantArchivedForFeedMedia %q - %+v", err, o)
			http.Error(w, "Uknown plant", http.StatusBadRequest)
			return
		}
		if plant.Archived == false {
			fn(w, r, p)
		}
	}
}
