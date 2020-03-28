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

// InitFeeds -
func InitFeeds(router *httprouter.Router) {
	initDB()

	anon := middleware.NewStack()
	anon.Use(wares.Logging)
	anon.Use(createDBSession)

	auth := middleware.NewStack()
	auth.Use(wares.Logging)
	auth.Use(jwtToken)
	auth.Use(createDBSession)

	router.POST("/user", anon.Wrap(createUserHandler()))

	router.POST("/userend", auth.Wrap(createUserEndHandler()))

	router.POST("/plantsharing", auth.Wrap(createPlantSharingHandler()))
	router.POST("/plant", auth.Wrap(createPlantHandler()))
	router.POST("/timelapse", auth.Wrap(createTimelapseHandler()))
	router.POST("/device", auth.Wrap(createDeviceHandler()))
	router.POST("/feed", auth.Wrap(createFeedHandler()))
	router.POST("/feedEntry", auth.Wrap(createFeedEntryHandler()))
}
