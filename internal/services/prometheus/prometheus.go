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
	"log"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	UUIDFilter = regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)
)

type HTTPTiming struct {
	router *httprouter.Router
}

func (ht *HTTPTiming) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := UUIDFilter.ReplaceAllString(r.URL.Path, "[UUID]")
	timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
	ht.router.ServeHTTP(w, r)
	timer.ObserveDuration()
	requestsCount.WithLabelValues(path, r.Method).Inc()
}

func NewHTTPTiming(router *httprouter.Router) *HTTPTiming {
	return &HTTPTiming{router}
}

func NotificationSent(notificationType string) {
	notificationsCount.WithLabelValues(notificationType).Inc()
}

func InitNotificationSent(notificationType string) {
	notificationsCount.WithLabelValues(notificationType)
}

func NotificationError(notificationType string) {
	notificationErrors.WithLabelValues(notificationType).Inc()
}

func AlertTriggered(metric, atype string) {
	alertsCount.WithLabelValues(metric, atype).Inc()
}

func InitAlertTriggered(metric, atype string) {
	alertsCount.WithLabelValues(metric, atype)
}

func Init() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()
}
