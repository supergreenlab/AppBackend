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

package devices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/data/kv"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

// TODO DRY with server/routes/feeds/select.go

func filterUserID(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)
		uid := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)
		selector = selector.Where("t.userid = ?", uid)
		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

func filterID(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		selector := r.Context().Value(middlewares.SelectorContextKey{}).(sqlbuilder.Selector)

		id := p.ByName("id")
		selector = selector.Where("t.id = ?", id)
		ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
		fn(w, r.WithContext(ctx), p)
	}
}

type SelectDevicesParamsParams struct {
	Params []string
}

type SelectDevicesParamsResponse struct {
	Params map[string]interface{} `json:"params"`
}

func loadParams(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		params := r.Context().Value(middlewares.QueryObjectContextKey{}).(*SelectDevicesParamsParams)
		device := r.Context().Value(middlewares.SelectResultContextKey{}).(*appbackend.Device)
		for i, k := range params.Params {
			params.Params[i] = fmt.Sprintf("%s.KV.%s", device.Identifier, k)
		}
		keys, err := kv.GetKeys(params.Params) // TODO Is this dangerous?
		if err != nil {
			logrus.Errorf("kv.GetKeys in loadParams %q - %+v %+v", err, params, keys)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		m, err := kv.GetValues(keys)
		if err != nil {
			logrus.Errorf("kv.GetValues in loadParams %q - %+v", err, params)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), middlewares.SelectResultContextKey{}, SelectDevicesParamsResponse{Params: m})
		fn(w, r.WithContext(ctx), p)
	}
}

var selectDeviceParams = middlewares.SelectOneEndpoint(
	"devices",
	func() interface{} { return &appbackend.Device{} },
	func() interface{} { return &SelectDevicesParamsParams{} },
	[]middleware.Middleware{
		filterID,
		filterUserID,
	},
	[]middleware.Middleware{
		loadParams,
	},
)
