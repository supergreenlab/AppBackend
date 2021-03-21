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
	"fmt"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/data/prometheus"
	"github.com/sirupsen/logrus"
)

type metricsParam struct {
	Metrics [][]float64 `json:"metrics"`
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
		logrus.Infof("%+v", feedEntries)
		for _, fe := range feedEntries {
			box, err := db.GetBoxFromPlantFeed(fe.FeedID)
			if err != nil {
				logrus.Errorf("db.GetBoxFromPlantFeed in cardMetricsProcess %q - fe: %+v", err, fe)
				continue
			}

			device, err := db.GetDeviceFromPlantFeed(fe.FeedID)
			if err != nil {
				logrus.Errorf("db.GetDeviceFromFeed in cardMetricsProcess %q - fe: %+v", err, fe)
				continue
			}
			rr, err := prometheus.QueryProm(fmt.Sprintf("g_%s{id=\"%s\"}", fmt.Sprintf("BOX_%d_TEMP", *box.DeviceBox), device.Identifier), t.Add(-18*time.Hour).Unix(), t.Add(18*time.Hour).Unix(), 40)

			if err != nil {
				logrus.Errorf("prometheus query failed: %s\n", err)
				continue
			}

			if rr.Status != "success" {
				logrus.Errorf("cid parameter error: %s\n", rr.Status)
				continue
			}

			time.Sleep(1 * time.Second)
		}
		time.Sleep(15 * time.Minute)
	}
}

func initCardMetrics() {
	go cardMetricsProcess()
}
