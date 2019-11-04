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

package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/prometheus"
)

func serveRangeHandler(w http.ResponseWriter, r *http.Request) {
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

	rr, err := prometheus.QueryProm(fmt.Sprintf(`g_%s{id="%s"}`, q, cid), time.Now().Unix()-60*60*int64(t), time.Now().Unix(), n)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	if rr.Status != "success" {
		w.WriteHeader(404)
		return
	}

	sr := newServedResult(rr, float64(min), float64(max))

	js, err := json.Marshal(sr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}
