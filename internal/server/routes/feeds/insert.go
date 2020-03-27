package feeds

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
)

func insertEndpoint(collection string, factory func() interface{}, addUserID bool) func() httprouter.Handle {
	return func() httprouter.Handle {
		s := middleware.NewStack()

		s.Use(decodeJSON(factory))
		if addUserID {
			s.Use(setUserID)
		}
		s.Use(insertObject(collection))

		return s.Wrap(outputObjectID)
	}
}

var createUserHandler = insertEndpoint("users", func() interface{} { return &User{} }, false)
var createPlantHandler = insertEndpoint("plants", func() interface{} { return &Plant{} }, true)
var createTimelapseHandler = insertEndpoint("timelapses", func() interface{} { return &Timelapse{} }, false)
var createDeviceHandler = insertEndpoint("devices", func() interface{} { return &Device{} }, false)
var createFeedHandler = insertEndpoint("feeds", func() interface{} { return &Feed{} }, true)
var createFeedEntryHandler = insertEndpoint("feedEntries", func() interface{} { return &FeedEntry{} }, false)
