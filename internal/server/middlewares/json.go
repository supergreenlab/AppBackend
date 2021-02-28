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
	"errors"
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/server/tools"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// ObjectContextKey - context key which stores the decoced object
type ObjectContextKey struct{}

// DecodeJSON - decodes the JSON payload
func DecodeJSON(fnObject func() interface{}) func(fn httprouter.Handle) httprouter.Handle {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			o := fnObject()
			err := tools.DecodeJSONBody(w, r, o)
			if err != nil {
				var mr *tools.MalformedRequest
				if errors.As(err, &mr) {
					logrus.Errorf("tools.DecodeJSONBody in DecodeJSON %q - %s", err, r.URL.String())
					http.Error(w, mr.Msg, mr.Status)
				} else {
					logrus.Errorf("tools.DecodeJSONBody in DecodeJSON %q - %s", err, r.URL.String())
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				return
			}
			ctx := context.WithValue(r.Context(), ObjectContextKey{}, o)
			fn(w, r.WithContext(ctx), p)
		}
	}
}
