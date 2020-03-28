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

	router.POST("/plant", authWithUserEndID.Wrap(createPlantHandler))
	router.POST("/timelapse", authWithUserEndID.Wrap(createTimelapseHandler))
	router.POST("/device", authWithUserEndID.Wrap(createDeviceHandler))
	router.POST("/feed", authWithUserEndID.Wrap(createFeedHandler))
	router.POST("/feedEntry", authWithUserEndID.Wrap(createFeedEntryHandler))
}
