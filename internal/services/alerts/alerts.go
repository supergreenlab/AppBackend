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

package alerts

import (
	"errors"
	"strconv"
	"strings"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/data/kv"
	"github.com/SuperGreenLab/AppBackend/internal/services/notifications"
	"github.com/SuperGreenLab/AppBackend/internal/services/prometheus"
	"github.com/SuperGreenLab/AppBackend/internal/services/pubsub"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func boxIDNumFromMetric(name string) (int, error) {
	parts := strings.Split(name, "_")
	return strconv.Atoi(parts[1])
}

var (
	alertTypeTooHigh = "TOO_HIGH"
	alertTypeTooLow  = "TOO_LOW"
)

type getMinMaxFunc func(controllerID string, boxID int, timerPower float64) (float64, float64, error)

type getAlertContentFunc func(plant db.Plant, alertType string, timerPower, value, minValue, maxValue float64) (string, string)

func checkMetric(metricName string, getAlertContent getAlertContentFunc, metric pubsub.ControllerIntMetric, getMinMax getMinMaxFunc, getSensorPresentForBox kv.GetSensorPresentForBoxFunc, getAlertStatus kv.GetAlertStatusFunc, setAlertStatus kv.SetAlertStatusFunc, getAlertType kv.GetAlertTypeFunc, setAlertType kv.SetAlertTypeFunc) {
	boxID, err := boxIDNumFromMetric(metric.Key)
	if err != nil {
		logrus.Errorf("boxIDNumFromMetric in checkMetric %q - %+v", err, metric)
		return
	}
	if enabled, err := kv.GetBoxEnabled(metric.ControllerID, boxID); !enabled || err != nil {
		if err != nil {
			logrus.Errorf("kv.GetBoxEnabled in checkMetric %q - %d", err, boxID)
		}
		return
	}
	if sht21Present, err := getSensorPresentForBox(metric.ControllerID, boxID); !sht21Present || err != nil {
		if err != nil {
			logrus.Errorf("getSensorPresentForBox in checkMetric %q - metric: %+v boxID: %d", err, metric, boxID)
		}
		return
	}
	timerPower, err := kv.GetTimerPower(metric.ControllerID, boxID)
	if err != nil {
		logrus.Errorf("kv.GetTimerPower in checkMetric %q - metric: %+v boxID: %d", err, metric, boxID)
		return
	}
	minValue, maxValue, err := getMinMax(metric.ControllerID, boxID, timerPower)
	if err != nil {
		logrus.Errorf("getMinMax in checkMetric %q - metric: %+v boxID: %d", err, metric, boxID)
		return
	}

	alertStatus, err := getAlertStatus(metric.ControllerID, boxID)
	if err != nil {
		logrus.Errorf("getAlertStatus in checkMetric %q - metric: %+v boxID: %d", err, metric, boxID)
		return
	}

	tooLow := metric.Value < minValue
	tooHigh := metric.Value > maxValue
	if tooLow || tooHigh {
		if alertStatus {
			return
		}
		err = setAlertStatus(metric.ControllerID, boxID, true)
		if err != nil {
			logrus.Errorf("setAlertStatus in checkMetric %q - metric: %+v boxID: %d", err, metric, boxID)
			return
		}

		alertType := ""
		if tooLow {
			alertType = alertTypeTooLow
		} else if tooHigh {
			alertType = alertTypeTooHigh
		}
		err = setAlertType(metric.ControllerID, boxID, alertType)
		if err != nil {
			logrus.Errorf("setAlertType in checkMetric %q - metric: %+v boxID: %d alertType: %s", err, metric, boxID, alertType)
			return
		}
		plants, err := db.GetActivePlantsForControllerIdentifier(metric.ControllerID, boxID)
		if err != nil {
			logrus.Errorf("db.GetPlantsForController in checkMetric %q - metric: %+v boxID: %d alertType: %s", err, metric, boxID, alertType)
			return
		}
		logrus.Infof("%s alert %s: %s{id=%s}=%f (timerPower: %f)", metricName, alertType, metric.Key, metric.ControllerID, metric.Value, timerPower)
		for _, plant := range plants {
			if plant.AlertsEnabled == false {
				continue
			}
			title, body := getAlertContent(plant, alertType, timerPower, metric.Value, minValue, maxValue)
			data, notif := NewNotificationDataAlert(title, body, "", plant.ID.UUID)
			notifications.SendNotificationToUser(plant.UserID, data, &notif)
			logrus.Infof("Sending notification %q %q %+v plant: %s feed: %s box: %s", title, body, metric, plant.ID.UUID, plant.FeedID, plant.BoxID)
			prometheus.AlertTriggered(metricName, alertType)
		}
	} else {
		if !alertStatus {
			return
		}

		alertType, err := getAlertType(metric.ControllerID, boxID)
		if err != nil && !errors.Is(err, redis.Nil) {
			logrus.Errorf("%q", err)
			return
		}

		if alertType == alertTypeTooLow && metric.Value < minValue*1.15 {
			return
		}
		if alertType == alertTypeTooHigh && metric.Value > maxValue/1.15 {
			return
		}

		err = setAlertStatus(metric.ControllerID, boxID, false)
		if err != nil {
			logrus.Errorf("%q", err)
			return
		}
		logrus.Infof("End %s alert: %s{id=%s}=%f", metricName, metric.Key, metric.ControllerID, metric.Value)
	}
}

func Init() {
	prometheus.InitNotificationSent(NotificationTypeReminder)
	prometheus.InitNotificationSent(NotificationTypeAlert)

	initTemperature()
	initHumidity()
}
