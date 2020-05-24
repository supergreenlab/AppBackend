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

package feeds

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

func fetchPublicPlant(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sess := r.Context().Value(sessContextKey{}).(sqlbuilder.Database)

	plant := Plant{}
	if err := sess.Select("*").From("plants").Where("is_public = ?", true).And("id = ?", p.ByName("id")).One(&plant); err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(plant); err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func fetchPublicFeedEntries(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sess := r.Context().Value(sessContextKey{}).(sqlbuilder.Database)

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit == 0 || limit > 50 {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	feedEntries := []FeedEntry{}
	if err := sess.Select("*").From("feedentries fe").Join("feeds f").On("fe.feedid = f.id").Join("plants p").On("p.feedid = f.id").Where("p.is_public = ?", true).And("p.id = ?", p.ByName("id")).And("fe.etype not in ('FE_TOWELIE_INFO', 'FE_PRODUCTS')").OrderBy("fe.cat DESC").Offset(offset).Limit(limit).All(&feedEntries); err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(feedEntries); err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
