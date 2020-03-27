package feeds

var createUserHandler = simpleInsert("users", func() interface{} { return User{} })
var createPlantHandler = simpleInsert("plants", func() interface{} { return Plant{} })
var createTimelapseHandler = simpleInsert("timelapses", func() interface{} { return Timelapse{} })
var createDeviceHandler = simpleInsert("devices", func() interface{} { return Device{} })
var createFeedHandler = simpleInsert("feeds", func() interface{} { return Feed{} })
var createFeedEntryHandler = simpleInsert("feedEntries", func() interface{} { return FeedEntry{} })
