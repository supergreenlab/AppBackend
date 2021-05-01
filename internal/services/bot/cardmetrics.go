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
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/data/kv"
	"github.com/SuperGreenLab/AppBackend/internal/data/prometheus"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/sirupsen/logrus"
)

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

			from := t.Add(-36 * time.Hour)
			to := t.Add(36 * time.Hour)
			meta := appbackend.FeedEntryMeta{
				MetricsMeta: appbackend.LoadMetricsMeta(device, box, from, to, prometheus.LoadTimeSeries),
			}

			j, err := json.Marshal(meta)
			if err != nil {
				logrus.Errorf("json.Marshal in cardMetricsProcess %q - box: %+v device: %+v", err, box, device)
				time.Sleep(1 * time.Second)
				continue
			}
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
