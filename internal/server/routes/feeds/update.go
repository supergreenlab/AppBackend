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
	"context"
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	fmiddlewares "github.com/SuperGreenLab/AppBackend/internal/server/routes/feeds/middlewares"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
)

var updateBoxHandler = middlewares.UpdateEndpoint(
	"boxes",
	func() interface{} { return &appbackend.Box{} },
	[]middleware.Middleware{
		middlewares.ObjectIDRequired,
		middlewares.SetUserID,
		middlewares.CheckAccessRight("boxes", "ID", false, func() appbackend.UserObject { return &appbackend.Box{} }),
		middlewares.CheckAccessRight("devices", "DeviceID", true, func() appbackend.UserObject { return &appbackend.Device{} }),
	},
	[]middleware.Middleware{
		fmiddlewares.UpdateUserEndObjects("userend_devices", "deviceid"),
	},
)

var updatePlantHandler = middlewares.UpdateEndpoint(
	"plants",
	func() interface{} { return &appbackend.Plant{} },
	[]middleware.Middleware{
		middlewares.ObjectIDRequired,
		middlewares.SetUserID,
		middlewares.CheckAccessRight("plants", "ID", false, func() appbackend.UserObject { return &appbackend.Plant{} }),
		middlewares.CheckAccessRight("boxes", "BoxID", false, func() appbackend.UserObject { return &appbackend.Box{} }),
	},
	[]middleware.Middleware{
		fmiddlewares.CheckPlantArchivedForPlant,
		fmiddlewares.UpdateUserEndObjects("userend_plants", "plantid"),
	},
)

var updateTimelapseHandler = middlewares.UpdateEndpoint(
	"timelapses",
	func() interface{} { return &appbackend.Timelapse{} },
	[]middleware.Middleware{
		middlewares.ObjectIDRequired,
		middlewares.SetUserID,
		middlewares.CheckAccessRight("timelapses", "ID", false, func() appbackend.UserObject { return &appbackend.Timelapse{} }),
		middlewares.CheckAccessRight("plants", "PlantID", false, func() appbackend.UserObject { return &appbackend.Plant{} }),
	},
	[]middleware.Middleware{
		fmiddlewares.CheckPlantArchivedForTimelapse,
		fmiddlewares.UpdateUserEndObjects("userend_timelapses", "timelapseid"),
	},
)

var updateDeviceHandler = middlewares.UpdateEndpoint(
	"devices",
	func() interface{} { return &appbackend.Device{} },
	[]middleware.Middleware{
		middlewares.ObjectIDRequired,
		middlewares.SetUserID,
		middlewares.CheckAccessRight("devices", "ID", false, func() appbackend.UserObject { return &appbackend.Device{} }),
	},
	[]middleware.Middleware{
		fmiddlewares.UpdateUserEndObjects("userend_devices", "deviceid"),
	},
)

var updateFeedHandler = middlewares.UpdateEndpoint(
	"feeds",
	func() interface{} { return &appbackend.Feed{} },
	[]middleware.Middleware{
		middlewares.ObjectIDRequired,
		middlewares.SetUserID,
		middlewares.CheckAccessRight("feeds", "ID", false, func() appbackend.UserObject { return &appbackend.Feed{} }),
	},
	[]middleware.Middleware{
		fmiddlewares.CheckPlantArchivedForFeed,
		fmiddlewares.UpdateUserEndObjects("userend_feeds", "feedid"),
	},
)

var updateFeedEntryHandler = middlewares.UpdateEndpoint(
	"feedentries",
	func() interface{} { return &appbackend.FeedEntry{} },
	[]middleware.Middleware{
		middlewares.ObjectIDRequired,
		middlewares.SetUserID,
		middlewares.CheckAccessRight("feedentries", "ID", false, func() appbackend.UserObject { return &appbackend.FeedEntry{} }),
		middlewares.CheckAccessRight("feeds", "FeedID", false, func() appbackend.UserObject { return &appbackend.Feed{} }),
	},
	[]middleware.Middleware{
		fmiddlewares.CheckPlantArchivedForFeedEntry,
		fmiddlewares.UpdateUserEndObjects("userend_feedentries", "feedentryid"),
	},
)

var updateFeedMediaHandler = middlewares.UpdateEndpoint(
	"feedmedias",
	func() interface{} { return &appbackend.FeedMedia{} },
	[]middleware.Middleware{
		middlewares.ObjectIDRequired,
		middlewares.SetUserID,
		middlewares.CheckAccessRight("feedmedias", "ID", false, func() appbackend.UserObject { return &appbackend.FeedMedia{} }),
		middlewares.CheckAccessRight("feedentries", "FeedEntryID", false, func() appbackend.UserObject { return &appbackend.FeedEntry{} }),
	},
	[]middleware.Middleware{
		fmiddlewares.CheckPlantArchivedForFeedMedia,
		fmiddlewares.UpdateUserEndObjects("userend_feedmedias", "feedmediaid"),
	},
)

func setUserEndID(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ueid := r.Context().Value(fmiddlewares.UserEndIDContextKey{}).(uuid.UUID)
		ue := r.Context().Value(middlewares.ObjectContextKey{}).(*db.UserEnd)

		ue.ID = uuid.NullUUID{UUID: ueid, Valid: true}

		ctx := context.WithValue(r.Context(), middlewares.ObjectContextKey{}, ue)
		fn(w, r.WithContext(ctx), p)
	}
}

var updateUserEndHandler = middlewares.UpdateEndpoint(
	"userends",
	func() interface{} { return &db.UserEnd{} },
	[]middleware.Middleware{
		middlewares.SetUserID,
		setUserEndID,
		middlewares.CheckAccessRight("userends", "ID", false, func() appbackend.UserObject { return &db.UserEnd{} }),
	},
	[]middleware.Middleware{},
)
