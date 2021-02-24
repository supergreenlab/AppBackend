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

package kv

import (
	"fmt"
	"time"
)

func GetTemperature(controllerID string, box int) (float64, error) {
	key := fmt.Sprintf("%s.KV.BOX_%d_TEMP", controllerID, box)
	return GetNum(key)
}

func GetBoxTempSource(controllerID string, box int) (int, error) {
	key := fmt.Sprintf("%s.KV.BOX_%d_TEMP_SOURCE", controllerID, box)
	return GetInt(key)
}

func GetSHT21Present(controllerID string, sht21 int) (bool, error) {
	key := fmt.Sprintf("%s.KV.SHT21_%d_PRESENT", controllerID, sht21)
	return GetBool(key)
}

func GetSHT21PresentForBox(controllerID string, box int) (bool, error) {
	tempSource, err := GetBoxTempSource(controllerID, box)
	if err != nil {
		return false, err
	}
	if tempSource == 0 {
		return false, nil
	}
	return GetSHT21Present(controllerID, tempSource-1)
}

func GetTimerPower(controllerID string, box int) (float64, error) {
	key := fmt.Sprintf("%s.KV.BOX_%d_TIMER_OUTPUT", controllerID, box)
	return GetNum(key)
}

func GetTemperatureAlertStatus(controllerID string, box int) (bool, error) {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_TEMP", controllerID, box)
	return GetBool(key)
}

func SetTemperatureAlertStatus(controllerID string, box int, value bool) error {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_TEMP", controllerID, box)
	return SetBool(key, value, time.Duration(30)*time.Minute)
}

func GetBoxEnabled(controllerID string, box int) (bool, error) {
	key := fmt.Sprintf("%s.KV.BOX_%d_ENABLED", controllerID, box)
	return GetBool(key)
}
