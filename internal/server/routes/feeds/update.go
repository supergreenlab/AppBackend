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

import "github.com/rileyr/middleware"

var updateBoxHandler = updateEndpoint(
	"boxes",
	func() interface{} { return &Box{} },
	[]middleware.Middleware{
		objectIDRequired,
		setUserID,
		checkAccessRight("boxes", "ID", false, func() interface{} { return &Box{} }),
		checkAccessRight("devices", "DeviceID", true, func() interface{} { return &Device{} }),
	},
	[]middleware.Middleware{
		updateUserEndObjects("userend_devices", "deviceid"),
	},
)

var updatePlantHandler = updateEndpoint(
	"plants",
	func() interface{} { return &Plant{} },
	[]middleware.Middleware{
		objectIDRequired,
		setUserID,
		checkAccessRight("plants", "ID", false, func() interface{} { return &Plant{} }),
		checkAccessRight("boxes", "BoxID", false, func() interface{} { return &Box{} }),
	},
	[]middleware.Middleware{
		updateUserEndObjects("userend_plants", "plantid"),
	},
)

var updateTimelapseHandler = updateEndpoint(
	"timelapses",
	func() interface{} { return &Timelapse{} },
	[]middleware.Middleware{
		objectIDRequired,
		setUserID,
		checkAccessRight("timelapses", "ID", false, func() interface{} { return &Timelapse{} }),
		checkAccessRight("plants", "PlantID", false, func() interface{} { return &Plant{} }),
	},
	[]middleware.Middleware{
		updateUserEndObjects("userend_timelapses", "timelapseid"),
	},
)

var updateDeviceHandler = updateEndpoint(
	"devices",
	func() interface{} { return &Device{} },
	[]middleware.Middleware{
		objectIDRequired,
		setUserID,
		checkAccessRight("devices", "ID", false, func() interface{} { return &Device{} }),
	},
	[]middleware.Middleware{
		updateUserEndObjects("userend_devices", "deviceid"),
	},
)

var updateFeedHandler = updateEndpoint(
	"feeds",
	func() interface{} { return &Feed{} },
	[]middleware.Middleware{
		objectIDRequired,
		setUserID,
		checkAccessRight("feeds", "ID", false, func() interface{} { return &Feed{} }),
	},
	[]middleware.Middleware{
		updateUserEndObjects("userend_feeds", "feedid"),
	},
)

var updateFeedEntryHandler = updateEndpoint(
	"feedentries",
	func() interface{} { return &FeedEntry{} },
	[]middleware.Middleware{
		objectIDRequired,
		setUserID,
		checkAccessRight("feedentries", "ID", false, func() interface{} { return &FeedEntry{} }),
		checkAccessRight("feeds", "FeedID", false, func() interface{} { return &Feed{} }),
	},
	[]middleware.Middleware{
		updateUserEndObjects("userend_feedentries", "feedentryid"),
	},
)

var updateFeedMediaHandler = updateEndpoint(
	"feedmedias",
	func() interface{} { return &FeedMedia{} },
	[]middleware.Middleware{
		objectIDRequired,
		setUserID,
		checkAccessRight("feedmedias", "ID", false, func() interface{} { return &FeedMedia{} }),
		checkAccessRight("feedentries", "FeedEntryID", false, func() interface{} { return &FeedEntry{} }),
	},
	[]middleware.Middleware{
		updateUserEndObjects("userend_feedmedias", "feedmediaid"),
	},
)
