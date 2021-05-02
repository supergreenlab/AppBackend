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

package appbackend

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type MetricsMeta struct {
	Temperature *TimeSeries   `json:"temperature,omitempty"`
	Humidity    *TimeSeries   `json:"humidity,omitempty"`
	VPD         *TimeSeries   `json:"vpd,omitempty"`
	Timer       *TimeSeries   `json:"timer,omitempty"`
	Dimming     *[]TimeSeries `json:"dimming,omitempty"`
	Ventilation *TimeSeries   `json:"ventilation,omitempty"`
}

type MetricsLoader func(device Device, from, to time.Time, module, metric string, i int) (TimeSeries, error)
type GetLedBox func(i int) (int, error)

func LoadMetricsMeta(device Device, box Box, from, to time.Time, loader MetricsLoader, getLedBox GetLedBox) MetricsMeta {
	meta := MetricsMeta{}
	if temp, err := loader(device, from, to, "BOX", "TEMP", int(*box.DeviceBox)); err == nil {
		meta.Temperature = &temp
	}
	if humi, err := loader(device, from, to, "BOX", "HUMI", int(*box.DeviceBox)); err == nil {
		meta.Humidity = &humi
	}
	if vpd, err := loader(device, from, to, "BOX", "VPD", int(*box.DeviceBox)); err == nil {
		meta.VPD = &vpd
	}
	if timer, err := loader(device, from, to, "BOX", "TIMER_OUTPUT", int(*box.DeviceBox)); err == nil {
		meta.Timer = &timer
	}
	dimmings := []TimeSeries{}
	for i := 0; ; i += 1 {
		if ledBox, err := getLedBox(i); err != nil || ledBox != int(*box.DeviceBox) {
			if err != nil {
				if !errors.Is(err, redis.Nil) {
					logrus.Errorf("kv.GetLedBox in cardMetricsProcess %q - box: %+v device: %+v i: %d", err, box, device, i)
				}
				break
			}
			continue
		}
		if dimming, err := loader(device, from, to, "LED", "DIM", i); err == nil {
			dimmings = append(dimmings, dimming)
		}
	}
	meta.Dimming = &dimmings
	if ventilation, err := loader(device, from, to, "BOX", "BLOWER_DUTY", int(*box.DeviceBox)); err == nil {
		meta.Ventilation = &ventilation
	}
	return meta
}
