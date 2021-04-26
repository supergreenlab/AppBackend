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
	minHumiDay   = 25
	maxHumiDay   = 75
	minHumiNight = 35
	maxHumiNight = 85
)

func getHumidityMinMax(controllerID string, boxID int, timerPower float64) (float64, float64, error) {
	minNight, err := kv.GetAlertMinHumidityNight(controllerID, boxID, minHumiNight)
	if err != nil {
		return 0, 0, err
	}
	minDay, err := kv.GetAlertMinHumidityDay(controllerID, boxID, minHumiDay)
	if err != nil {
		return 0, 0, err
	}
	maxNight, err := kv.GetAlertMaxHumidityNight(controllerID, boxID, maxHumiNight)
	if err != nil {
		return 0, 0, err
	}
	maxDay, err := kv.GetAlertMaxHumidityDay(controllerID, boxID, maxHumiDay)
	if err != nil {
		return 0, 0, err
	}

	return minNight + (minDay-minNight)*timerPower/100, maxNight + (maxDay-maxNight)*timerPower/100, nil
}

func getHumidityAlertContent(plant appbackend.Plant, alertType string, timerPower, value, minValue, maxValue float64) (string, string) {
	alertTypesToText := map[string]string{
		alertTypeTooHigh: "too humid",
		alertTypeTooLow:  "too dry",
	}

	title := fmt.Sprintf("Humidity alert")
	body := fmt.Sprintf("Your plant %s is %s\nIt's currently at %d%%", plant.Name, alertTypesToText[alertType], int(value))
	if alertType == alertTypeTooHigh {
		body = fmt.Sprintf("%s, try to keep it below %d%%", body, int(maxValue))
	} else {
		body = fmt.Sprintf("%s, try to keep it above %d%%", body, int(minValue))
	}
	if timerPower == 0 {
		body = fmt.Sprintf("%s during the night.", body)
	} else {
		body = fmt.Sprintf("%s during the day.", body)
	}
	return title, body
}

func listenHumidityMetrics() {
	prometheus.InitAlertTriggered("HUMI", alertTypeTooLow)
	prometheus.InitAlertTriggered("HUMI", alertTypeTooHigh)
	ch := pubsub.SubscribeControllerIntMetric("*.BOX_*_HUMI")
	for metric := range ch {
		checkMetric("HUMI", getHumidityAlertContent, metric, getHumidityMinMax, kv.GetSHT21PresentForBox, kv.GetHumidityAlertStatus, kv.SetHumidityAlertStatus, kv.GetHumidityAlertType, kv.SetHumidityAlertType)
	}
}

func initHumidity() {
	go listenHumidityMetrics()
}
