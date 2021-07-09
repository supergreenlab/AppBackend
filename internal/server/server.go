/*
 * Copyright (C) 2019  SuperGreenLab <towelie@supergreenlab.com>
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

package server

import (
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/server/routes/devices"
	"github.com/SuperGreenLab/AppBackend/internal/server/routes/products"
	"github.com/SuperGreenLab/AppBackend/internal/server/routes/services"
	"github.com/SuperGreenLab/AppBackend/internal/services/prometheus"
	"github.com/rs/cors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/SuperGreenLab/AppBackend/internal/data/storage"

	"github.com/SuperGreenLab/AppBackend/internal/server/routes/feeds"
	"github.com/SuperGreenLab/AppBackend/internal/server/routes/metrics"
	"github.com/SuperGreenLab/AppBackend/internal/server/routes/users"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

var (
	_ = pflag.String("addcors", "false", "Adds cors header, used when not behing cloudflare") // TODO move this somewhere else
)

func init() {
	viper.SetDefault("AddCORS", "false")
}

// Start starts the server
func Start() {
	storage.SetupBucket("feedmedias")
	storage.SetupBucket("users")
	storage.SetupBucket("timelapses")

	router := httprouter.New()

	users.Init(router)
	metrics.Init(router)
	feeds.Init(router)
	products.Init(router)
	devices.Init(router)
	services.Init(router)

	go func() {
		if viper.GetString("AddCORS") == "true" {
			corsOpts := cors.Options{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{
					http.MethodHead,
					http.MethodGet,
					http.MethodPost,
					http.MethodPut,
					http.MethodPatch,
					http.MethodDelete,
				},
				AllowedHeaders:   []string{"*"},
				AllowCredentials: false,
				ExposedHeaders:   []string{"x-sgl-token"},
			}

			log.Fatal(http.ListenAndServe(":8080", cors.New(corsOpts).Handler(prometheus.NewHTTPTiming(router))))
		} else {
			log.Fatal(http.ListenAndServe(":8080", prometheus.NewHTTPTiming(router)))
		}
	}()
}
