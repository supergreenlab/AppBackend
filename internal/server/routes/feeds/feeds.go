package feeds

import (
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
	anon.Use(createDBSession)
	return anon
}

func authStack(withUserEndID bool) middleware.Stack {
	auth := middleware.NewStack()
	auth.Use(wares.Logging)
	auth.Use(jwtToken)
	auth.Use(createDBSession)

	if withUserEndID == true {
		auth.Use(userEndIDRequired)
	}

	return auth
}

// InitFeeds -
func InitFeeds(router *httprouter.Router) {
	initDB()
	initStorage()

	anon := anonStack()
	auth := authStack(false)
	authWithUserEndID := authStack(true)

	router.POST("/login", anon.Wrap(loginHandler()))

	router.POST("/user", anon.Wrap(createUserHandler))

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
}
