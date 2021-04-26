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
	"fmt"

	"github.com/SuperGreenLab/AppBackend/internal/data/kv"
	"github.com/SuperGreenLab/AppBackend/internal/services/prometheus"
	"github.com/SuperGreenLab/AppBackend/internal/services/pubsub"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
)

const (
	minTempDay   = 18
	maxTempDay   = 32
	minTempNight = 15
	maxTempNight = 25
)

func toDegF(degC float64) float64 {
	return (degC * 9 / 5) + 32
}

func getTemperatureMinMax(controllerID string, boxID int, timerPower float64) (float64, float64, error) {
	minNight, err := kv.GetAlertMinTemperatureNight(controllerID, boxID, minTempNight)
	if err != nil {
		return 0, 0, err
	}
	minDay, err := kv.GetAlertMinTemperatureDay(controllerID, boxID, minTempDay)
	if err != nil {
		return 0, 0, err
	}
	maxNight, err := kv.GetAlertMaxTemperatureNight(controllerID, boxID, maxTempNight)
	if err != nil {
		return 0, 0, err
	}
	maxDay, err := kv.GetAlertMaxTemperatureDay(controllerID, boxID, maxTempDay)
	if err != nil {
		return 0, 0, err
	}

	return minNight + (minDay-minNight)*timerPower/100, maxNight + (maxDay-maxNight)*timerPower/100, nil
}

func getTemperatureAlertContent(plant appbackend.Plant, alertType string, timerPower, value, minValue, maxValue float64) (string, string) {
	alertTypesToText := map[string]string{
		alertTypeTooHigh: "too hot",
		alertTypeTooLow:  "too cold",
	}

	title := fmt.Sprintf("Temperature alert")
	body := fmt.Sprintf("Your plant %s is %s\nIt's currently at %d°C (%d°F)", plant.Name, alertTypesToText[alertType], int(value), int(toDegF(value)))
	if alertType == alertTypeTooHigh {
		body = fmt.Sprintf("%s, try to keep it below %d°C (%d°F)", body, int(maxValue), int(toDegF(maxValue)))
	} else {
		body = fmt.Sprintf("%s, try to keep it above %d°C (%d°F)", body, int(minValue), int(toDegF(minValue)))
	}
	if timerPower == 0 {
		body = fmt.Sprintf("%s during the night.", body)
	} else {
		body = fmt.Sprintf("%s during the day.", body)
	}
	return title, body
}

func listenTemperatureMetrics() {
	prometheus.InitAlertTriggered("TEMP", alertTypeTooLow)
	prometheus.InitAlertTriggered("TEMP", alertTypeTooHigh)
	ch := pubsub.SubscribeControllerIntMetric("*.BOX_*_TEMP")
	for metric := range ch {
		checkMetric("TEMP", getTemperatureAlertContent, metric, getTemperatureMinMax, kv.GetSHT21PresentForBox, kv.GetTemperatureAlertStatus, kv.SetTemperatureAlertStatus, kv.GetTemperatureAlertType, kv.SetTemperatureAlertType)
	}
}

func initTemperature() {
	go listenTemperatureMetrics()
}
