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

package explorer

import (
	"encoding/json"
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

func fetchPublicPlants(makeSelector func(sess sqlbuilder.Database, w http.ResponseWriter, r *http.Request, p httprouter.Params) sqlbuilder.Selector) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)

		selector := makeSelector(sess, w, r, p).
			Where("plants.is_public = ?", true).
			And("plants.deleted = ?", false)

		selector = joinLatestFeedMedia(sess, selector)
		selector = joinBoxSettings(selector)
		selector = joinFollows(r, selector)
		selector = pageOffsetLimit(r, selector)

		results := []publicPlantResult{}
		if err := selector.All(&results); err != nil {
			logrus.Errorf("selector.All in fetchPublicPlants %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for i, p := range results {
			err := explorer.LoadFeedMediaPublicURLs(&p)
			if err != nil {
				logrus.Errorf("tools.LoadFeedMediaPublicURLs in fetchPublicPlants %q - p: %+v", err, p)
				continue
			}
			results[i] = p
		}

		if err := json.NewEncoder(w).Encode(publicPlantsResult{results}); err != nil {
			logrus.Errorf("json.NewEncoder in fetchPublicPlants %q - results: %+v", err, results)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
