package feeds

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
)

func insertEndpoint(
	collection string,
	factory func() interface{},
	preInsert []middleware.Middleware,
	postInsert []middleware.Middleware,
) httprouter.Handle {
	s := middleware.NewStack()

	s.Use(decodeJSON(factory))
	if preInsert != nil {
		for _, m := range preInsert {
			s.Use(m)
		}
	}
	s.Use(insertObject(collection))

	if postInsert != nil {
		for _, m := range postInsert {
			s.Use(m)
		}
	}

	return s.Wrap(outputObjectID)
}
