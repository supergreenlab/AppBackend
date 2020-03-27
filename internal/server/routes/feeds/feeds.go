package feeds

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/rileyr/middleware/wares"
)

// InitFeeds -
func InitFeeds(router *httprouter.Router) {
	initDB()

	anon := middleware.NewStack()
	anon.Use(wares.Logging)
	anon.Use(createDBSession)

	auth := middleware.NewStack()
	auth.Use(wares.Logging)
	auth.Use(createDBSession)

	router.POST("/user", anon.Wrap(createUserHandler()))
	router.POST("/plant", auth.Wrap(createPlantHandler()))
	router.POST("/timelapse", auth.Wrap(createTimelapseHandler()))
	router.POST("/device", auth.Wrap(createDeviceHandler()))
	router.POST("/feed", auth.Wrap(createFeedHandler()))
	router.POST("/feedEntry", auth.Wrap(createFeedEntryHandler()))
}
