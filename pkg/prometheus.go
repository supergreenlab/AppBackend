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
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

type TimeSeries [][]float64

func (mv TimeSeries) minMax() (float64, float64) {
	min := math.MaxFloat64
	max := math.SmallestNonzeroFloat32

	for _, v := range mv {
		min = math.Min(min, v[1])
		max = math.Max(max, v[1])
	}

	return min, max
}

func (mv TimeSeries) current() float64 {
	if len(mv) < 1 {
		return 0
	}
	return mv[len(mv)-1][1]
}

func LoadGraphValue(device Device, from, to time.Time, module, metric string, i int) (TimeSeries, error) {
	m := TimeSeries{}

	name := fmt.Sprintf("%s_%d_%s", module, i, metric)
	url := fmt.Sprintf("https://api2.supergreenlab.com/metrics?cid=%s&q=%s&timeFrom=%d&timeTo=%d&n=50", device.Identifier, name, from.Unix(), to.Unix())
	r, err := http.Get(url)
	if err != nil {
		return m, err
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&m)
	return m, nil
}
