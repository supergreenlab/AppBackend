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
	"log"
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/SuperGreenLab/AppBackend/internal/services/pubsub"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

var upgrader = websocket.Upgrader{}

func streamDeviceMetrics(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer c.Close()
		identifier := p.ByName("id")
		q := fmt.Sprintf("pub.%s.*", identifier)
		ch := pubsub.SubscribeControllerMetric(q)
		for e := range ch {
			err = c.WriteJSON(e)
			if err != nil {
				logrus.Errorf("c.WriteJSON in streamDeviceMetrics %q - device: %s", err, identifier)
				break
			}
		}
	}
}

var streamDeviceHandler = middlewares.NewEndpoint().
	PushMiddleware(func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
			uid := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)
			id := p.ByName("id")

			selector := sess.Select("*").From("devices t").Where("t.userid = ?", uid).And("identifier = ?", id).And("deleted = false")
			ctx := context.WithValue(r.Context(), middlewares.SelectorContextKey{}, selector)
			fn(w, r.WithContext(ctx), p)
		}
	}).
	PushMiddleware(middlewares.SelectOneQuery(func() interface{} { return &appbackend.Device{} })).
	PushMiddleware(streamDeviceMetrics).Handle()
