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

func fetchPublicPlant(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)

	plant := publicPlantResult{}
	selector := sess.Select("plants.id", "plants.name", "plants.settings").
		From("plants").
		Where("plants.is_public = ?", true).
		And("plants.deleted = ?", false).
		And("plants.id = ?", p.ByName("id"))

	selector = joinLatestFeedMedia(sess, selector)
	selector = joinBoxSettings(selector)

	if err := selector.One(&plant); err != nil {
		logrus.Errorf("sess.Select('plants') in fetchPublicPlant %q - id: %s", err, p.ByName("id"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err := loadFeedMediaPublicURLs(&plant)
	if err != nil {
		logrus.Errorf("loadFeedMediaPublicURLs in fetchPublicPlant %q - plant: %+v", err, plant)
	}

	if err := json.NewEncoder(w).Encode(plant); err != nil {
		logrus.Errorf("json.NewEncoder in fetchPublicPlant %q - plant: %+v", err, plant)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
