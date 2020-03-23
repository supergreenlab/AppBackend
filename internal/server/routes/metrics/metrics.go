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

package metrics

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/prometheus"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

var (
	queryFilter = regexp.MustCompile("[^a-zA-Z0-9_]*")
	cidFilter   = regexp.MustCompile("[^a-f0-9]*")
)

func ServeMetricsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	q := r.URL.Query().Get("q")
	if q == "" {
		log.Error("q parameter error missing")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	timeFrom := time.Now().Unix() - 60*60*72
	timeTo := time.Now().Unix()
	q = queryFilter.ReplaceAllString(q, "")
	if r.URL.Query().Get("t") != "" {
		t, err := strconv.Atoi(r.URL.Query().Get("t"))
		if err != nil {
			log.Errorf("t parameter error: %s\n", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		timeFrom = time.Now().Unix() - 60*60*int64(t)
	} else if r.URL.Query().Get("t1") != "" && r.URL.Query().Get("t2") != "" {
		t1, err := strconv.Atoi(r.URL.Query().Get("t1"))
		if err != nil {
			log.Errorf("t parameter error: %s\n", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		t2, err := strconv.Atoi(r.URL.Query().Get("t2"))
		if err != nil {
			log.Errorf("t parameter error: %s\n", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		timeFrom = time.Now().Unix() - 60*60*int64(t1)
		timeTo = time.Now().Unix() - 60*60*int64(t2)
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
		log.Errorf("cid parameter error: %s\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	cid = cidFilter.ReplaceAllString(cid, "")

	rr, err := prometheus.QueryProm(fmt.Sprintf("g_%s{id=\"%s\"}", q, cid), timeFrom, timeTo, n)

	if err != nil {
		log.Errorf("prometheus query failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rr.Status != "success" {
		log.Errorf("cid parameter error: %s\n", rr.Status)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sr := newServedResult(rr, float64(min), float64(max))

	js, err := json.Marshal(sr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
