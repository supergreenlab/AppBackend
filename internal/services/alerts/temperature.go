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
	"github.com/SuperGreenLab/AppBackend/internal/data/kv"
	"github.com/SuperGreenLab/AppBackend/internal/services/prometheus"
	"github.com/SuperGreenLab/AppBackend/internal/services/pubsub"
)

func getTemperatureMinMax(timerPower float64) (float64, float64) {
	var minNight, maxNight float64 = 15, 25
	var minDay, maxDay float64 = 18, 32
	return minNight + (minDay-minNight)*timerPower/100, maxNight + (maxDay-maxNight)*timerPower/100
}

func listenTemperatureMetrics() {
	prometheus.InitAlertTriggered("TEMP", "TOO_LOW")
	prometheus.InitAlertTriggered("TEMP", "TOO_HIGH")
	ch := pubsub.SubscribeControllerIntMetric("*.BOX_*_TEMP")
	for metric := range ch {
		go checkMetric("TEMP", metric, getTemperatureMinMax, kv.GetSHT21PresentForBox, kv.GetTemperatureAlertStatus, kv.SetTemperatureAlertStatus, kv.GetTemperatureAlertType, kv.SetTemperatureAlertType)
	}
}

func initTemperature() {
	go listenTemperatureMetrics()
}
