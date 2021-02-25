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

	"github.com/SuperGreenLab/AppBackend/internal/data/kv"
	"github.com/SuperGreenLab/AppBackend/internal/services/prometheus"
	"github.com/SuperGreenLab/AppBackend/internal/services/pubsub"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func boxIDNumFromMetric(name string) (int, error) {
	parts := strings.Split(name, "_")
	return strconv.Atoi(parts[1])
}

type getMinMaxFunc func(timerPower float64) (float64, float64)

func checkMetric(metricName string, metric pubsub.ControllerIntMetric, getMinMax getMinMaxFunc, getSensorPresentForBox kv.GetSensorPresentForBoxFunc, getAlertStatus kv.GetAlertStatusFunc, setAlertStatus kv.SetAlertStatusFunc, getAlertType kv.GetAlertTypeFunc, setAlertType kv.SetAlertTypeFunc) {
	boxID, err := boxIDNumFromMetric(metric.Key)
	if err != nil {
		logrus.Errorf("%q", err)
		return
	}
	if enabled, err := kv.GetBoxEnabled(metric.ControllerID, boxID); !enabled || err != nil {
		if err != nil {
			logrus.Errorf("%q", err)
		}
		return
	}
	if sht21Present, err := getSensorPresentForBox(metric.ControllerID, boxID); !sht21Present || err != nil {
		if err != nil {
			logrus.Errorf("%q", err)
		}
		return
	}
	timerPower, err := kv.GetTimerPower(metric.ControllerID, boxID)
	if err != nil {
		logrus.Errorf("%q", err)
		return
	}
	minValue, maxValue := getMinMax(timerPower)

	alertStatus, err := getAlertStatus(metric.ControllerID, boxID)
	if err != nil {
		logrus.Errorf("%q", err)
		return
	}

	tooLow := metric.Value <= minValue
	tooHigh := metric.Value >= maxValue
	if tooLow || tooHigh {
		if alertStatus {
			return
		}
		err = setAlertStatus(metric.ControllerID, boxID, true)
		if err != nil {
			logrus.Errorf("%q", err)
			return
		}

		alertType := ""
		if tooLow {
			alertType = "TOO_LOW"
		} else if tooHigh {
			alertType = "TOO_HIGH"
		}
		prometheus.AlertTriggered(metricName, alertType)
		err = setAlertType(metric.ControllerID, boxID, alertType)
		if err != nil {
			logrus.Errorf("%q", err)
			return
		}
		logrus.Infof("%s alert %s: %s{id=%s}=%f (timerPower: %f)", metricName, alertType, metric.Key, metric.ControllerID, metric.Value, timerPower)
	} else {
		if !alertStatus {
			return
		}

		alertType, err := getAlertType(metric.ControllerID, boxID)
		if err != nil && !errors.Is(err, redis.Nil) {
			logrus.Errorf("%q", err)
			return
		}

		if alertType == "TOO_LOW" && metric.Value < minValue*1.15 {
			return
		}
		if alertType == "TOO_HIGH" && metric.Value > maxValue/1.15 {
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
	initTemperature()
	initHumidity()
}
