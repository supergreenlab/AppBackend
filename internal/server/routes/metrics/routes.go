package metrics

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/rileyr/middleware/wares"
)

// Init -
func Init(router *httprouter.Router) {
	s := middleware.NewStack()

	s.Use(wares.Logging)

	router.GET("/metrics", s.Wrap(ServeMetricsHandler))
}
