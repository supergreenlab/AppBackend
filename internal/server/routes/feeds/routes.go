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
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/rileyr/middleware/wares"
	"github.com/spf13/pflag"
)

var (
	jwtSecret = pflag.String("jwtsecret", "", "JWT secret")
)

func anonStack() middleware.Stack {
	anon := middleware.NewStack()
	anon.Use(wares.Logging)
	anon.Use(middlewares.CreateDBSession)
	return anon
}

func authStack(withUserEndID bool) middleware.Stack {
	auth := middleware.NewStack()
	auth.Use(wares.Logging)
	auth.Use(middlewares.JwtToken)
	auth.Use(middlewares.CreateDBSession)

	if withUserEndID == true {
		auth.Use(middlewares.UserEndIDRequired)
	}

	return auth
}

// InitFeeds -
func InitFeeds(router *httprouter.Router) {
	anon := anonStack()
	auth := authStack(false)
	authWithUserEndID := authStack(true)

	router.POST("/userend", auth.Wrap(createUserEndHandler))
	router.POST("/plantsharing", auth.Wrap(createPlantSharingHandler))

	router.POST("/box", authWithUserEndID.Wrap(createBoxHandler))
	router.POST("/plant", authWithUserEndID.Wrap(createPlantHandler))
	router.POST("/timelapse", authWithUserEndID.Wrap(createTimelapseHandler))
	router.POST("/device", authWithUserEndID.Wrap(createDeviceHandler))
	router.POST("/feed", authWithUserEndID.Wrap(createFeedHandler))
	router.POST("/feedEntry", authWithUserEndID.Wrap(createFeedEntryHandler))
	router.POST("/feedMedia", authWithUserEndID.Wrap(createFeedMediaHandler))

	router.PUT("/box", authWithUserEndID.Wrap(updateBoxHandler))
	router.PUT("/plant", authWithUserEndID.Wrap(updatePlantHandler))
	router.PUT("/timelapse", authWithUserEndID.Wrap(updateTimelapseHandler))
	router.PUT("/device", authWithUserEndID.Wrap(updateDeviceHandler))
	router.PUT("/feed", authWithUserEndID.Wrap(updateFeedHandler))
	router.PUT("/feedEntry", authWithUserEndID.Wrap(updateFeedEntryHandler))
	router.PUT("/feedMedia", authWithUserEndID.Wrap(updateFeedMediaHandler))

	router.POST("/deletes", authWithUserEndID.Wrap(deletesHandler))

	router.POST("/feedMediaUploadURL", authWithUserEndID.Wrap(feedMediaUploadURLHandler))

	router.GET("/syncBoxes", authWithUserEndID.Wrap(syncBoxesHandler))
	router.GET("/syncPlants", authWithUserEndID.Wrap(syncPlantsHandler))
	router.GET("/syncTimelapses", authWithUserEndID.Wrap(syncTimelapsesHandler))
	router.GET("/syncDevices", authWithUserEndID.Wrap(syncDevicesHandler))
	router.GET("/syncFeeds", authWithUserEndID.Wrap(syncFeedsHandler))
	router.GET("/syncFeedEntries", authWithUserEndID.Wrap(syncFeedEntriesHandler))
	router.GET("/syncFeedMedias", authWithUserEndID.Wrap(syncFeedMediasHandler))

	router.POST("/box/:id/sync", authWithUserEndID.Wrap(syncedBoxHandler))
	router.POST("/plant/:id/sync", authWithUserEndID.Wrap(syncedPlantHandler))
	router.POST("/timelapse/:id/sync", authWithUserEndID.Wrap(syncedTimelapseHandler))
	router.POST("/device/:id/sync", authWithUserEndID.Wrap(syncedDeviceHandler))
	router.POST("/feed/:id/sync", authWithUserEndID.Wrap(syncedFeedHandler))
	router.POST("/feedEntry/:id/sync", authWithUserEndID.Wrap(syncedFeedEntryHandler))
	router.POST("/feedMedia/:id/sync", authWithUserEndID.Wrap(syncedFeedMediaHandler))

	router.GET("/public/plants", anon.Wrap(fetchPublicPlants))
	router.GET("/public/plant/:id", anon.Wrap(fetchPublicPlant))
	router.GET("/public/plant/:id/feedEntries", anon.Wrap(fetchPublicFeedEntries))
	router.GET("/public/feedEntry/:id/feedMedias", anon.Wrap(fetchPublicFeedMedias))
	router.GET("/public/feedMedia/:id", anon.Wrap(fetchPublicFeedMedia))
}
