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

package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/data/kv"
	"github.com/SuperGreenLab/AppBackend/internal/data/prometheus"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func loadTimeSeries(device db.Device, from, to int64, module, metric string, i int) (prometheus.TimeSeries, error) {
	rr, err := prometheus.QueryProm(fmt.Sprintf("g_%s{id=\"%s\"}", fmt.Sprintf("%s_%d_%s", module, i, metric), device.Identifier), from, to, 50)

	if err != nil {
		logrus.Errorf("prometheus.QueryProm in loadTimeSeries %q - device: %+v from: %d to: %d module: %s metric: %s i: %d", err, device, from, to, module, metric, i)
		return prometheus.TimeSeries{}, err
	}

	if rr.Status != "success" {
		err := errors.New(fmt.Sprintf("cid parameter error: %s", rr.Status))
		logrus.Errorf("prometheus.QueryProm in loadTimeSeries %q - device: %+v from: %d to: %d module: %s metric: %s i: %d", err, device, from, to, module, metric, i)
		return prometheus.TimeSeries{}, err
	}
	return rr.ToFloat64(float64(math.MinInt32), float64(math.MaxInt32)), nil
}

func cardMetricsProcess() {
	for {
		t := time.Now().Add(-36 * time.Hour)
		feedEntries, err := db.GetFeedEntriesBetweenDates(t.Add(-8*time.Minute), t.Add(8*time.Minute))
		if err != nil {
			time.Sleep(60 * time.Second)
			logrus.Errorf("db.GetFeedEntriesBetweenDates in cardMetricsProcess %q", err)
			continue
		}
		for _, fe := range feedEntries {
			box, err := db.GetBoxFromPlantFeed(fe.FeedID)
			if err != nil {
				logrus.Errorf("db.GetBoxFromPlantFeed in cardMetricsProcess %q - fe: %+v", err, fe)
				time.Sleep(1 * time.Second)
				continue
			}

			if !box.DeviceID.Valid {
				time.Sleep(1 * time.Second)
				continue
			}

			device, err := db.GetDeviceFromPlantFeed(fe.FeedID)
			if err != nil {
				logrus.Errorf("db.GetDeviceFromPlantFeed in cardMetricsProcess %q - fe: %+v", err, fe)
				time.Sleep(1 * time.Second)
				continue
			}

			if sht21Present, err := kv.GetSHT21PresentForBox(device.Identifier, int(*box.DeviceBox)); !sht21Present || err != nil {
				if err != nil {
					logrus.Errorf("getSensorPresentForBox in cardMetricsProcess %q - box: %+v device: %+v", err, box, device)
				}
				time.Sleep(1 * time.Second)
				continue
			}

			meta := db.FeedEntryMeta{}
			from := t.Add(-36 * time.Hour).Unix()
			to := t.Add(36 * time.Hour).Unix()
			if temp, err := loadTimeSeries(device, from, to, "BOX", "TEMP", int(*box.DeviceBox)); err == nil {
				meta.Temperature = temp
			}
			if humi, err := loadTimeSeries(device, from, to, "BOX", "HUMI", int(*box.DeviceBox)); err == nil {
				meta.Humidity = humi
			}
			if vpd, err := loadTimeSeries(device, from, to, "BOX", "VPD", int(*box.DeviceBox)); err == nil {
				meta.VPD = vpd
			}
			if timer, err := loadTimeSeries(device, from, to, "BOX", "TIMER_OUTPUT", int(*box.DeviceBox)); err == nil {
				meta.Timer = timer
			}
			dimmings := []prometheus.TimeSeries{}
			for i := 0; ; i += 1 {
				if ledBox, err := kv.GetLedBox(device.Identifier, i); err != nil || ledBox != int(*box.DeviceBox) {
					if err != nil {
						if !errors.Is(err, redis.Nil) {
							logrus.Errorf("kv.GetLedBox in cardMetricsProcess %q - box: %+v device: %+v i: %d", err, box, device, i)
						}
						break
					}
					continue
				}
				if dimming, err := loadTimeSeries(device, from, to, "LED", "DIM", i); err == nil {
					dimmings = append(dimmings, dimming)
				}
			}
			meta.Dimming = dimmings
			if ventilation, err := loadTimeSeries(device, from, to, "BOX", "BLOWER_DUTY", int(*box.DeviceBox)); err == nil {
				meta.Ventilation = ventilation
			}

			j, err := json.Marshal(meta)
			if err != nil {
				logrus.Errorf("json.Marshal in cardMetricsProcess %q - box: %+v device: %+v", err, box, device)
				time.Sleep(1 * time.Second)
				continue
			}
			logrus.Infof("%s", string(j))
			if err := db.SetFeedEntryMeta(fe.ID.UUID, string(j)); err != nil {
				logrus.Errorf("db.SetFeedEntryMeta in cardMetricsProcess %q - fe: %+v j: %s", err, fe, string(j))
			}

			time.Sleep(1 * time.Second)
		}
		time.Sleep(15 * time.Minute)
	}
}

func initCardMetrics() {
	go cardMetricsProcess()
}
