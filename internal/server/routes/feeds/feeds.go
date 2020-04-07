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
	//router.POST("/feedMedia/:id/uploadURL", authWithUserEndID.Wrap(uploadURLHandler))

	router.GET("/box/sync", authWithUserEndID.Wrap(syncBoxesHandler))
	router.GET("/plant/sync", authWithUserEndID.Wrap(syncPlantsHandler))
	router.GET("/timelapse/sync", authWithUserEndID.Wrap(syncTimelapsesHandler))
	router.GET("/device/sync", authWithUserEndID.Wrap(syncDevicesHandler))
	router.GET("/feed/sync", authWithUserEndID.Wrap(syncFeedsHandler))
	router.GET("/feedEntry/sync", authWithUserEndID.Wrap(syncFeedEntriesHandler))
	router.GET("/feedMedia/sync", authWithUserEndID.Wrap(syncFeedMediasHandler))

	router.POST("/box/sync/:id", authWithUserEndID.Wrap(syncedBoxHandler))
	router.POST("/plant/sync/:id", authWithUserEndID.Wrap(syncedPlantHandler))
	router.POST("/timelapse/sync/:id", authWithUserEndID.Wrap(syncedTimelapseHandler))
	router.POST("/device/sync/:id", authWithUserEndID.Wrap(syncedDeviceHandler))
	router.POST("/feed/sync/:id", authWithUserEndID.Wrap(syncedFeedHandler))
	router.POST("/feedEntry/sync/:id", authWithUserEndID.Wrap(syncedFeedEntryHandler))
	router.POST("/feedMedia/sync/:id", authWithUserEndID.Wrap(syncedFeedMediaHandler))
}
