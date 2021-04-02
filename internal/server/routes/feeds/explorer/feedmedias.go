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

package feeds

import (
	"encoding/json"
	"net/http"

	sgldb "github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

func fetchPublicFeedMedias(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)

	feedMedias := []sgldb.FeedMedia{}
	selector := sess.Select("fm.*").From("feedmedias fm").
		Join("feedentries fe").On("fm.feedentryid = fe.id").
		Join("feeds f").On("fe.feedid = f.id").
		Join("plants p").On("p.feedid = f.id").
		Where("p.is_public = ?", true).
		And("fe.id = ?", p.ByName("id")).
		And("fm.deleted = ?", false)
	if err := selector.All(&feedMedias); err != nil {
		logrus.Errorf("selector.All in fetchPublicFeedMedias %q - id: %s", err, p.ByName("id"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var err error
	for i, fm := range feedMedias {
		err = loadFeedMediaPublicURLs(&fm)
		if err != nil {
			logrus.Errorf("loadFeedMediaPublicURLs in fetchPublicFeedMedias %q - %+v", err, fm)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// might not be useful anymore
		feedMedias[i] = fm
	}

	result := publicFeedMediasResult{feedMedias}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logrus.Errorf("json.NewEncoder in fetchPublicFeedMedias %q - %+v", err, result)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
