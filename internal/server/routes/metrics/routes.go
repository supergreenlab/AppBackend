package metrics

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/rileyr/middleware/wares"
	"github.com/spf13/viper"
)

// Init -
func Init(router *httprouter.Router) {
	s := middleware.NewStack()

	if viper.GetString("LogRequests") == "true" {
		s.Use(wares.Logging)
	}

	router.GET("/metrics", s.Wrap(ServeMetricsHandler))
}
