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

package middlewares

import (
	"fmt"
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/services/pubsub"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
)

type InsertMessage struct {
	ID     uuid.UUID   `json:"id"`
	Object interface{} `json:"object"`
}

func PublishInsert(collection string) middleware.Middleware {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			id := r.Context().Value(InsertedIDContextKey{}).(uuid.UUID)
			o := r.Context().Value(ObjectContextKey{})

			msg := InsertMessage{id, o}
			if err := pubsub.PublishObject(fmt.Sprintf("insert.%s", collection), msg); err != nil {
				logrus.Errorf("PublishObject in PublishInsert %q", err)
			}
			fn(w, r, p)
		}
	}
}
