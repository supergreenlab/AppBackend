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
	"strconv"

	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
)

// RangeResult query result from prometheus
type RangeResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name     string `json:"__name__"`
				ID       string `json:"id"`
				Instance string `json:"instance"`
				Job      string `json:"job"`
				Module   string `json:"module"`
			} `json:"metric"`
			Values [][]interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// ToFloat64 returns the RangeResult as an array of [timestamp, value]
func (r RangeResult) ToFloat64(min, max float64) appbackend.TimeSeries {
	res := [][]float64{}
	if len(r.Data.Result) < 1 {
		return res
	}
	var lasti float64
	for _, v := range r.Data.Result[0].Values {
		i, err := strconv.ParseFloat(v[1].(string), 64)
		if err != nil || i < min || i > max {
			i = lasti
		} else {
			lasti = i
		}
		res = append(res, []float64{
			v[0].(float64),
			i,
		})
	}
	return res
}
