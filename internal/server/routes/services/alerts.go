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

package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/SuperGreenLab/AppBackend/internal/services/alerts"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
)

type AlertsSettings struct {
	MinTempDay float64 `json:"minTempDay"`
	MaxTempDay float64 `json:"maxTempDay"`

	MinTempNight float64 `json:"minTempNight"`
	MaxTempNight float64 `json:"maxTempNight"`

	MinHumiDay float64 `json:"minHumiDay"`
	MaxHumiDay float64 `json:"maxHumiDay"`

	MinHumiNight float64 `json:"minHumiNight"`
	MaxHumiNight float64 `json:"maxHumiNight"`
}

func getAlertsSettings(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		plant := r.Context().Value(middlewares.SelectResultContextKey{}).(*appbackend.Plant)

		box, err := db.GetBox(plant.BoxID)
		if err != nil {
			logrus.Errorln("db.GetBox in services.getAlertsSettings %q %+v", err, plant)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if box.DeviceID.Valid == false {
			errMsg := fmt.Sprintf("Missing device for plant %+v", plant)
			logrus.Infof("in services.getAlertsSettings %q", errMsg)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		device, err := db.GetDevice(box.DeviceID.UUID)
		if err != nil {
			logrus.Errorln("db.GetDevice in services.getAlertsSettings %q %+v", err, plant)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tempAlertSettings, err := alerts.GetTemperatureAlertSettings(device.Identifier, int(*box.DeviceBox))
		if err != nil {
			logrus.Errorln("alerts.GetTemperatureAlertSettings in services.getAlertsSettings %q %+v", err, plant)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		humiAlertSettings, err := alerts.GetHumidityAlertSettings(device.Identifier, int(*box.DeviceBox))
		if err != nil {
			logrus.Errorln("alerts.GetHumidityAlertSettings in services.getAlertsSettings %q %+v", err, plant)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), middlewares.SelectResultContextKey{}, AlertsSettings{
			MinTempDay:   tempAlertSettings.MinDay,
			MaxTempDay:   tempAlertSettings.MaxDay,
			MinTempNight: tempAlertSettings.MinNight,
			MaxTempNight: tempAlertSettings.MaxNight,

			MinHumiDay:   humiAlertSettings.MinDay,
			MaxHumiDay:   humiAlertSettings.MaxDay,
			MinHumiNight: humiAlertSettings.MinNight,
			MaxHumiNight: humiAlertSettings.MaxNight,
		})
		fn(w, r.WithContext(ctx), p)
	}
}

var selectAlertsSettings = middlewares.SelectOneEndpoint(
	"plants",
	func() interface{} { return &appbackend.Plant{} },
	nil,
	[]middleware.Middleware{
		middlewares.FilterID,
		middlewares.FilterUserID,
	},
	[]middleware.Middleware{
		getAlertsSettings,
	},
)

func setAlertsSettings(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		params := r.Context().Value(middlewares.ObjectContextKey{}).(*AlertsSettings)
		plant := r.Context().Value(middlewares.SelectResultContextKey{}).(*appbackend.Plant)

		box, err := db.GetBox(plant.BoxID)
		if err != nil {
			logrus.Errorln("db.GetBox in services.setAlertsSettings %q %+v", err, plant)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if box.DeviceID.Valid == false {
			errMsg := fmt.Sprintf("Missing device for plant %+v", plant)
			logrus.Infof("in services.setAlertsSettings %q", errMsg)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		device, err := db.GetDevice(box.DeviceID.UUID)
		if err != nil {
			logrus.Errorln("db.GetDevice in services.setAlertsSettings %q %+v", err, plant)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tas := alerts.TemperatureAlertSettings{
			MinDay:   params.MinTempDay,
			MaxDay:   params.MaxTempDay,
			MinNight: params.MinTempNight,
			MaxNight: params.MaxTempNight,
		}
		err = alerts.SetTemperatureAlertSettings(device.Identifier, int(*box.DeviceBox), tas)
		if err != nil {
			logrus.Errorln("alerts.SetTemperatureAlertSettings in services.setAlertsSettings %q %+v", err, plant)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		has := alerts.HumidityAlertSettings{
			MinDay:   params.MinHumiDay,
			MaxDay:   params.MaxHumiDay,
			MinNight: params.MinHumiNight,
			MaxNight: params.MaxHumiNight,
		}
		err = alerts.SetHumidityAlertSettings(device.Identifier, int(*box.DeviceBox), has)
		if err != nil {
			logrus.Errorln("alerts.SetHumidityAlertSettings in services.setAlertsSettings %q %+v", err, plant)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		middlewares.OutputOK(w, r, p)
	}
}

var updateAlertsSettings = middlewares.SelectOneEndpoint(
	"plants",
	func() interface{} { return &appbackend.Plant{} },
	nil,
	[]middleware.Middleware{
		middlewares.DecodeJSON(func() interface{} { return &AlertsSettings{} }),
		middlewares.FilterID,
		middlewares.FilterUserID,
	},
	[]middleware.Middleware{
		setAlertsSettings,
	},
)
