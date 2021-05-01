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

package prometheus

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/sirupsen/logrus"
)

func LoadTimeSeries(device appbackend.Device, from, to time.Time, module, metric string, i int) (appbackend.TimeSeries, error) {
	rr, err := QueryProm(fmt.Sprintf("g_%s{id=\"%s\"}", fmt.Sprintf("%s_%d_%s", module, i, metric), device.Identifier), from.Unix(), to.Unix(), 50)

	if err != nil {
		logrus.Errorf("QueryProm in loadTimeSeries %q - device: %+v from: %d to: %d module: %s metric: %s i: %d", err, device, from.Unix(), to.Unix(), module, metric, i)
		return appbackend.TimeSeries{}, err
	}

	if rr.Status != "success" {
		err := errors.New(fmt.Sprintf("cid parameter error: %s", rr.Status))
		logrus.Errorf("QueryProm in loadTimeSeries %q - device: %+v from: %d to: %d module: %s metric: %s i: %d", err, device, from.Unix(), to.Unix(), module, metric, i)
		return appbackend.TimeSeries{}, err
	}
	return rr.ToFloat64(float64(math.MinInt32), float64(math.MaxInt32)), nil
}

// QueryProm fetches metrics from prometheus
func QueryProm(query string, start, end int64, n int) (RangeResult, error) {
	res := RangeResult{}
	c := http.DefaultClient

	v := url.Values{}
	v.Set("query", query)
	v.Set("start", fmt.Sprintf("%d", start))
	v.Set("end", fmt.Sprintf("%d", end))
	v.Set("step", fmt.Sprintf("%d", (end-start)/int64(n)))
	u, err := url.Parse(fmt.Sprintf("http://prometheus:9090/api/v1/query_range?%s", v.Encode()))
	if err != nil {
		return res, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return res, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&res)
	return res, err
}
