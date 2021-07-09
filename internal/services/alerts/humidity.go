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

type HumidityAlertSettings struct {
	MinNight float64 `json:"minNight"`
	MinDay   float64 `json:"minDay"`

	MaxNight float64 `json:"maxNight"`
	MaxDay   float64 `json:"maxDay"`
}

func GetHumidityAlertSettings(controllerID string, boxID int) (*HumidityAlertSettings, error) {
	minNight, err := kv.GetAlertMinHumidityNight(controllerID, boxID, minTempNight)
	if err != nil {
		return nil, err
	}
	minDay, err := kv.GetAlertMinHumidityDay(controllerID, boxID, minTempDay)
	if err != nil {
		return nil, err
	}
	maxNight, err := kv.GetAlertMaxHumidityNight(controllerID, boxID, maxTempNight)
	if err != nil {
		return nil, err
	}
	maxDay, err := kv.GetAlertMaxHumidityDay(controllerID, boxID, maxTempDay)
	if err != nil {
		return nil, err
	}

	return &HumidityAlertSettings{
		MinNight: minNight,
		MinDay:   minDay,
		MaxNight: maxNight,
		MaxDay:   maxDay,
	}, nil
}

func SetHumidityAlertSettings(controllerID string, boxID int, as HumidityAlertSettings) error {
	err := kv.SetAlertMinHumidityNight(controllerID, boxID, as.MinNight)
	if err != nil {
		return err
	}
	err = kv.SetAlertMinHumidityDay(controllerID, boxID, as.MinDay)
	if err != nil {
		return err
	}
	err = kv.SetAlertMaxHumidityNight(controllerID, boxID, as.MaxNight)
	if err != nil {
		return err
	}
	err = kv.SetAlertMaxHumidityDay(controllerID, boxID, as.MaxDay)
	if err != nil {
		return err
	}

	return nil
}

func getHumidityMinMax(controllerID string, boxID int, timerPower float64) (float64, float64, error) {
	as, err := GetHumidityAlertSettings(controllerID, boxID)
	if err != nil {
		return 0, 0, err
	}

	return as.MinNight + (as.MinDay-as.MinNight)*timerPower/100, as.MaxNight + (as.MaxDay-as.MaxNight)*timerPower/100, nil
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
	ch, _ := pubsub.SubscribeControllerIntMetric("*.BOX_*_HUMI.metric")
	for metric := range ch {
		checkMetric("HUMI", getHumidityAlertContent, metric, getHumidityMinMax, kv.GetSHT21PresentForBox, kv.GetHumidityAlertStatus, kv.SetHumidityAlertStatus, kv.GetHumidityAlertType, kv.SetHumidityAlertType)
	}
}

func initHumidity() {
	go listenHumidityMetrics()
}
