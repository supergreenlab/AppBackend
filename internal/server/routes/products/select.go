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

package products

import (
	"encoding/json"
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

type searchProductsResult struct {
	Products []db.Products `json:"products"`
}

func searchProducts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	terms := r.URL.Query().Get("terms")
	if terms == "" {
		http.Error(w, "missing 'terms' parameter", http.StatusInternalServerError)
		return
	}
	category := r.URL.Query().Get("category")
	sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
	selector := sess.Select("p.*")
	selector = selector.From("products p")
	selector = selector.Where("p.name % ?", terms)
	if category != "" {
		selector = selector.Where("jsonb_exists(p.categories, ?)", category)
	}
	products := []db.Products{}
	if err := selector.All(&products); err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(searchProductsResult{products}); err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
