package feeds

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
)

func insertEndpoint(
	collection string,
	factory func() interface{},
	pre []middleware.Middleware,
	post []middleware.Middleware,
) httprouter.Handle {
	s := middleware.NewStack()

	s.Use(decodeJSON(factory))
	if pre != nil {
		for _, m := range pre {
			s.Use(m)
		}
	}
	s.Use(insertObject(collection))

	if post != nil {
		for _, m := range post {
			s.Use(m)
		}
	}

	return s.Wrap(outputObjectID)
}

func updateEndpoint(
	collection string,
	factory func() interface{},
	pre []middleware.Middleware,
	post []middleware.Middleware,
) httprouter.Handle {
	s := middleware.NewStack()

	s.Use(decodeJSON(factory))
	if pre != nil {
		for _, m := range pre {
			s.Use(m)
		}
	}
	s.Use(updateObject(collection))

	if post != nil {
		for _, m := range post {
			s.Use(m)
		}
	}

	return s.Wrap(outputOK)
}
