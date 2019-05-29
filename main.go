/*
 * Copyright (C) 2019  SuperGreenLab <towelie@supergreenlab.com>
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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// curl 'http://localhost:9090/api/v1/query_range?query=g_BOX_0_SHT1X_TEMP_C%7Bid%3D%22a5a524ceee3a7d80%22%7D&start=1552927230&end=1552934445&step=15'

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

func (r RangeResult) toFloat64(min, max float64) [][]float64 {
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

func queryProm(query string, start, end int64, n int, res *RangeResult) error {
	c := http.DefaultClient

	v := url.Values{}
	v.Set("query", query)
	v.Set("start", fmt.Sprintf("%d", start))
	v.Set("end", fmt.Sprintf("%d", end))
	v.Set("step", fmt.Sprintf("%d", (end-start)/int64(n)))
	u, err := url.Parse(fmt.Sprintf("http://prometheus:9090/api/v1/query_range?%s", v.Encode()))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	/*b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", b)*/

	err = json.NewDecoder(resp.Body).Decode(&res)
	return err
}

type ServedResult struct {
	Metrics [][]float64 `json:"metrics"`
}

func serveRange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Access-Control-Allow-Origin")
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	} else if r.Method != "GET" {
		w.WriteHeader(404)
		return
	}

	q := r.URL.Query().Get("q")
	if q == "" {
		w.WriteHeader(404)
		return
	}

	t, err := strconv.Atoi(r.URL.Query().Get("t"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}

	n, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil {
		n = 200
	}

	min, err := strconv.Atoi(r.URL.Query().Get("min"))
	if err != nil {
		min = math.MinInt32
	}

	max, err := strconv.Atoi(r.URL.Query().Get("max"))
	if err != nil {
		max = math.MaxInt32
	}

	cid := r.URL.Query().Get("cid")
	if cid == "" {
		w.WriteHeader(404)
		return
	}

	rr := RangeResult{}
	queryProm(fmt.Sprintf(`g_%s{id="%s"}`, q, cid), time.Now().Unix()-60*60*int64(t), time.Now().Unix(), n, &rr)

	if rr.Status != "success" {
		w.WriteHeader(404)
		return
	}

	sr := ServedResult{Metrics: rr.toFloat64(float64(min), float64(max))}

	js, err := json.Marshal(sr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func main() {
	http.HandleFunc("/", serveRange)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
