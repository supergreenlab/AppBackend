package feeds

import "github.com/rileyr/middleware"

var createUserHandler = insertEndpoint(
	"users",
	func() interface{} { return &User{} },
	nil,
	nil,
)

var createUserEndHandler = insertEndpoint(
	"userends",
	func() interface{} { return &UserEnd{} },
	[]middleware.Middleware{setUserID},
	nil,
)

var createPlantHandler = insertEndpoint(
	"plants",
	func() interface{} { return &Plant{} },
	[]middleware.Middleware{setUserID},
	nil,
)

var createTimelapseHandler = insertEndpoint(
	"timelapses",
	func() interface{} { return &Timelapse{} },
	[]middleware.Middleware{setUserID},
	nil,
)

var createDeviceHandler = insertEndpoint(
	"devices",
	func() interface{} { return &Device{} },
	nil,
	nil,
)

var createFeedHandler = insertEndpoint(
	"feeds",
	func() interface{} { return &Feed{} },
	[]middleware.Middleware{setUserID},
	nil,
)

var createFeedEntryHandler = insertEndpoint(
	"feedentries",
	func() interface{} { return &FeedEntry{} },
	nil,
	nil,
)

var createPlantSharingHandler = insertEndpoint(
	"plantsharings",
	func() interface{} { return &PlantSharing{} },
	[]middleware.Middleware{setUserID},
	nil,
)
