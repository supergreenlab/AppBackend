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
	"math/rand"
	"time"
)

func GetAlertMinTemperatureDay(controllerID string, box int, def float64) (float64, error) {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_MIN_TEMP_DAY", controllerID, box)
	return GetNum(key, def)
}

func GetAlertMaxTemperatureDay(controllerID string, box int, def float64) (float64, error) {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_MAX_TEMP_DAY", controllerID, box)
	return GetNum(key, def)
}

func GetAlertMinTemperatureNight(controllerID string, box int, def float64) (float64, error) {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_MIN_TEMP_NIGHT", controllerID, box)
	return GetNum(key, def)
}

func GetAlertMaxTemperatureNight(controllerID string, box int, def float64) (float64, error) {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_MAX_TEMP_NIGHT", controllerID, box)
	return GetNum(key, def)
}

func GetAlertMinHumidityDay(controllerID string, box int, def float64) (float64, error) {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_MIN_HUMI_DAY", controllerID, box)
	return GetNum(key, def)
}

func GetAlertMaxHumidityDay(controllerID string, box int, def float64) (float64, error) {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_MAX_HUMI_DAY", controllerID, box)
	return GetNum(key, def)
}

func GetAlertMinHumidityNight(controllerID string, box int, def float64) (float64, error) {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_MIN_HUMI_NIGHT", controllerID, box)
	return GetNum(key, def)
}

func GetAlertMaxHumidityNight(controllerID string, box int, def float64) (float64, error) {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_MAX_HUMI_NIGHT", controllerID, box)
	return GetNum(key, def)
}

func GetTemperature(controllerID string, box int) (float64, error) {
	key := fmt.Sprintf("%s.KV.BOX_%d_TEMP", controllerID, box)
	return GetNum(key, 0)
}

func GetBoxTempSource(controllerID string, box int) (int, error) {
	key := fmt.Sprintf("%s.KV.BOX_%d_TEMP_SOURCE", controllerID, box)
	return GetInt(key, 0)
}

type GetSensorPresentFunc func(controllerID string, sht21 int) (bool, error)

func GetSHT21Present(controllerID string, sht21 int) (bool, error) {
	key := fmt.Sprintf("%s.KV.SHT21_%d_PRESENT", controllerID, sht21)
	return GetBool(key)
}

type GetSensorPresentForBoxFunc func(controllerID string, box int) (bool, error)

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
	return GetNum(key, 0)
}

func GetLedBox(controllerID string, led int) (int, error) {
	key := fmt.Sprintf("%s.KV.LED_%d_DIM", controllerID, led)
	n, err := r.Get(key).Int()
	return n, err
}

type GetAlertStatusFunc func(controllerID string, box int) (bool, error)

func GetTemperatureAlertStatus(controllerID string, box int) (bool, error) {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_TEMP", controllerID, box)
	return GetBool(key)
}

func GetHumidityAlertStatus(controllerID string, box int) (bool, error) {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_HUMI", controllerID, box)
	return GetBool(key)
}

type SetAlertStatusFunc func(controllerID string, box int, value bool) error

func SetTemperatureAlertStatus(controllerID string, box int, value bool) error {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_TEMP", controllerID, box)
	return SetBool(key, value, time.Duration(30+rand.Int()%15)*time.Minute)
}

func SetHumidityAlertStatus(controllerID string, box int, value bool) error {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_HUMI", controllerID, box)
	return SetBool(key, value, time.Duration(30+rand.Int()%15)*time.Minute)
}

type GetAlertTypeFunc func(controllerID string, box int) (string, error)

func GetTemperatureAlertType(controllerID string, box int) (string, error) {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_TEMP_TYPE", controllerID, box)
	return GetString(key)
}

func GetHumidityAlertType(controllerID string, box int) (string, error) {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_HUMI_TYPE", controllerID, box)
	return GetString(key)
}

type SetAlertTypeFunc func(controllerID string, box int, atype string) error

func SetTemperatureAlertType(controllerID string, box int, atype string) error {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_TEMP_TYPE", controllerID, box)
	return SetString(key, atype)
}

func SetHumidityAlertType(controllerID string, box int, atype string) error {
	key := fmt.Sprintf("%s.ALERT.BOX_%d_HUMI_TYPE", controllerID, box)
	return SetString(key, atype)
}

func GetBoxEnabled(controllerID string, box int) (bool, error) {
	key := fmt.Sprintf("%s.KV.BOX_%d_ENABLED", controllerID, box)
	return GetBool(key)
}
