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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "appbackend_http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
	notificationsCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "appbackend_notifications",
		Help: "Number of notifications",
	}, []string{"type"})
	notificationErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "appbackend_notification_errors",
		Help: "Number of notification errors",
	}, []string{"type"})
	requestsCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "appbackend_requests",
		Help: "Number of http requests",
	}, []string{"path", "method"})
	alertsCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "appbackend_alerts",
		Help: "Number of alerts",
	}, []string{"metric", "type"})
)
