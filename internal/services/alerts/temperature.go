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
	"strconv"
	"strings"

	"github.com/SuperGreenLab/AppBackend/internal/data/kv"
	"github.com/SuperGreenLab/AppBackend/internal/services/prometheus"
	"github.com/SuperGreenLab/AppBackend/internal/services/pubsub"
	"github.com/sirupsen/logrus"
)

func boxIDNumFromMetric(name string) (int, error) {
	parts := strings.Split(name, "_")
	return strconv.Atoi(parts[1])
}

func listenTemperatureMetrics() {
	ch := pubsub.SubscribeControllerIntMetric("*.BOX_*_TEMP")
	for metric := range ch {
		boxID, err := boxIDNumFromMetric(metric.Key)
		if err != nil {
			logrus.Errorf("%q\n", err)
			continue
		}
		if enabled, err := kv.GetBoxEnabled(metric.ControllerID, boxID); !enabled || err != nil {
			if err != nil {
				logrus.Errorf("%q\n", err)
			}
			continue
		}
		if sht21Present, err := kv.GetSHT21PresentForBox(metric.ControllerID, boxID); !sht21Present || err != nil {
			if err != nil {
				logrus.Errorf("%q\n", err)
			}
			continue
		}
		timerPower, err := kv.GetTimerPower(metric.ControllerID, boxID)
		if err != nil {
			logrus.Errorf("%q\n", err)
			continue
		}
		var minTemp, maxTemp float64
		if timerPower == 0 {
			minTemp = 15
			maxTemp = 25
		} else if timerPower != 0 {
			minTemp = 21
			maxTemp = 30
		}

		alertStatus, err := kv.GetTemperatureAlertStatus(metric.ControllerID, boxID)
		if err != nil {
			logrus.Errorf("%q\n", err)
			continue
		}

		tooLow := metric.Value <= minTemp
		tooHigh := metric.Value >= maxTemp
		if tooLow || tooHigh {
			if alertStatus {
				continue
			}
			err = kv.SetTemperatureAlertStatus(metric.ControllerID, boxID, true)
			if err != nil {
				logrus.Errorf("%q\n", err)
				continue
			}

			alertType := ""
			if tooLow {
				alertType = "TOO_LOW"
			} else if tooHigh {
				alertType = "TOO_HIGH"
			}
			logrus.Infof("Temp alert %s: %s{id=%s}=%f\n", alertType, metric.Key, metric.ControllerID, metric.Value)
			prometheus.AlertTriggered("TEMP", alertType, metric.ControllerID, strconv.Itoa(boxID))
		} else {
			if !alertStatus {
				continue
			}

			err = kv.SetTemperatureAlertStatus(metric.ControllerID, boxID, true)
			if err != nil {
				logrus.Errorf("%q\n", err)
				continue
			}
			logrus.Infof("End temp alert: %s{id=%s}=%f\n", metric.Key, metric.ControllerID, metric.Value)
		}
	}
}

func initTemperature() {
	go listenTemperatureMetrics()
}
