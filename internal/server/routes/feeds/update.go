/*
 * Copyright (C) 2020  SuperGreenLab <towelie@supergreenlab.com>
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

package feeds

import (
	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/rileyr/middleware"
)

var updateBoxHandler = middlewares.UpdateEndpoint(
	"boxes",
	func() interface{} { return &db.Box{} },
	[]middleware.Middleware{
		middlewares.ObjectIDRequired,
		middlewares.SetUserID,
		middlewares.CheckAccessRight("boxes", "ID", false, func() db.UserObject { return &db.Box{} }),
		middlewares.CheckAccessRight("devices", "DeviceID", true, func() db.UserObject { return &db.Device{} }),
	},
	[]middleware.Middleware{
		middlewares.UpdateUserEndObjects("userend_devices", "deviceid"),
	},
)

var updatePlantHandler = middlewares.UpdateEndpoint(
	"plants",
	func() interface{} { return &db.Plant{} },
	[]middleware.Middleware{
		middlewares.ObjectIDRequired,
		middlewares.SetUserID,
		middlewares.CheckAccessRight("plants", "ID", false, func() db.UserObject { return &db.Plant{} }),
		middlewares.CheckAccessRight("boxes", "BoxID", false, func() db.UserObject { return &db.Box{} }),
	},
	[]middleware.Middleware{
		middlewares.UpdateUserEndObjects("userend_plants", "plantid"),
	},
)

var updateTimelapseHandler = middlewares.UpdateEndpoint(
	"timelapses",
	func() interface{} { return &db.Timelapse{} },
	[]middleware.Middleware{
		middlewares.ObjectIDRequired,
		middlewares.SetUserID,
		middlewares.CheckAccessRight("timelapses", "ID", false, func() db.UserObject { return &db.Timelapse{} }),
		middlewares.CheckAccessRight("plants", "PlantID", false, func() db.UserObject { return &db.Plant{} }),
	},
	[]middleware.Middleware{
		middlewares.UpdateUserEndObjects("userend_timelapses", "timelapseid"),
	},
)

var updateDeviceHandler = middlewares.UpdateEndpoint(
	"devices",
	func() interface{} { return &db.Device{} },
	[]middleware.Middleware{
		middlewares.ObjectIDRequired,
		middlewares.SetUserID,
		middlewares.CheckAccessRight("devices", "ID", false, func() db.UserObject { return &db.Device{} }),
	},
	[]middleware.Middleware{
		middlewares.UpdateUserEndObjects("userend_devices", "deviceid"),
	},
)

var updateFeedHandler = middlewares.UpdateEndpoint(
	"feeds",
	func() interface{} { return &db.Feed{} },
	[]middleware.Middleware{
		middlewares.ObjectIDRequired,
		middlewares.SetUserID,
		middlewares.CheckAccessRight("feeds", "ID", false, func() db.UserObject { return &db.Feed{} }),
	},
	[]middleware.Middleware{
		middlewares.UpdateUserEndObjects("userend_feeds", "feedid"),
	},
)

var updateFeedEntryHandler = middlewares.UpdateEndpoint(
	"feedentries",
	func() interface{} { return &db.FeedEntry{} },
	[]middleware.Middleware{
		middlewares.ObjectIDRequired,
		middlewares.SetUserID,
		middlewares.CheckAccessRight("feedentries", "ID", false, func() db.UserObject { return &db.FeedEntry{} }),
		middlewares.CheckAccessRight("feeds", "FeedID", false, func() db.UserObject { return &db.Feed{} }),
	},
	[]middleware.Middleware{
		middlewares.UpdateUserEndObjects("userend_feedentries", "feedentryid"),
	},
)

var updateFeedMediaHandler = middlewares.UpdateEndpoint(
	"feedmedias",
	func() interface{} { return &db.FeedMedia{} },
	[]middleware.Middleware{
		middlewares.ObjectIDRequired,
		middlewares.SetUserID,
		middlewares.CheckAccessRight("feedmedias", "ID", false, func() db.UserObject { return &db.FeedMedia{} }),
		middlewares.CheckAccessRight("feedentries", "FeedEntryID", false, func() db.UserObject { return &db.FeedEntry{} }),
	},
	[]middleware.Middleware{
		middlewares.UpdateUserEndObjects("userend_feedmedias", "feedmediaid"),
	},
)
