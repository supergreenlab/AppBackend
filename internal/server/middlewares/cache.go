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
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/kv"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type CacheKeyFn func(r *http.Request, p httprouter.Params) string

func SelectCacheResult(cacheKeyFn CacheKeyFn) func(fn httprouter.Handle) httprouter.Handle {
	return func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			key := cacheKeyFn(r, p)
			key = fmt.Sprintf("cache.%s", key)
			ctx := context.WithValue(r.Context(), CacheKeyContextKey{}, key)
			if v, err := kv.GetString(key); err != nil {
				logrus.Errorf("kv.GetString in SelectCacheResult %q - key: %s", err, key)
				lastKey := fmt.Sprintf("%s.last", key)
				if vlast, err := kv.GetString(lastKey); err != nil {
					logrus.Errorf("kv.GetString in SelectCacheResult %q - lastKey: %s", err, lastKey)
					fn(w, r.WithContext(ctx), p)
				} else {
					go func() {
						working, _ := kv.GetBool(fmt.Sprintf("%s.working", key))
						if !working {
							kv.SetBool(fmt.Sprintf("%s.working", key), true, time.Second*4) // TODO find something better
							fn(httptest.NewRecorder(), r.WithContext(ctx), p)
							kv.SetBool(fmt.Sprintf("%s.working", key), false, 0)
						} else {
							logrus.Info("skipping work")
						}
					}()
					w.Write([]byte(vlast))
				}
			} else {
				w.Write([]byte(v))
			}
		}
	}
}
