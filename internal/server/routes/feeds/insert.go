package feeds

var createUserHandler = simpleInsert("users", func() interface{} { return &User{} }, false)
var createPlantHandler = simpleInsert("plants", func() interface{} { return &Plant{} }, true)
var createTimelapseHandler = simpleInsert("timelapses", func() interface{} { return &Timelapse{} }, false)
var createDeviceHandler = simpleInsert("devices", func() interface{} { return &Device{} }, false)
var createFeedHandler = simpleInsert("feeds", func() interface{} { return &Feed{} }, true)
var createFeedEntryHandler = simpleInsert("feedEntries", func() interface{} { return &FeedEntry{} }, false)
