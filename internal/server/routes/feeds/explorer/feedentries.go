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

	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

func fetchPublicFeedEntries(makeSelector func(sess sqlbuilder.Database, w http.ResponseWriter, r *http.Request, p httprouter.Params) sqlbuilder.Selector) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)

		selector := makeSelector(sess, w, r, p)

		selector = joinFeedEntrySocialSelector(r, selector)
		selector = publicFeedEntriesOnly(selector)
		selector = pageOffsetLimit(r, selector)

		feedEntries := []publicFeedEntry{}
		if err := selector.All(&feedEntries); err != nil {
			logrus.Errorf("selector.All in fetchPublicFeedEntries %q - id: %s", err, p.ByName("id"))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for i, p := range feedEntries {
			err := loadFeedMediaPublicURLs(&p)
			if err != nil {
				logrus.Errorf("loadFeedMediaPublicURLs in fetchPublicFeedEntries %q - p: %+v", err, p)
				continue
			}
			feedEntries[i] = p
		}

		result := publicFeedEntriesResult{feedEntries}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			logrus.Errorf("json.NewEncoder in fetchPublicFeedEntries %q - %+v", err, result)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
