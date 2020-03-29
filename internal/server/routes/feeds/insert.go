package feeds

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"golang.org/x/crypto/bcrypt"
	"upper.io/db.v3/lib/sqlbuilder"
)

var createUserHandler = insertEndpoint(
	"users",
	func() interface{} { return &User{} },
	[]middleware.Middleware{
		func(fn httprouter.Handle) httprouter.Handle {
			return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				u := r.Context().Value(objectContextKey{}).(*User)
				sess := r.Context().Value(sessContextKey{}).(sqlbuilder.Database)
				n, err := sess.Collection("users").Find().Where("nickname = ?", u.Nickname).Count()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if n > 0 {
					http.Error(w, "User already exists", http.StatusBadRequest)
					return
				}

				bc, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
				u.Password = string(bc)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				fn(w, r, p)
			}
		},
	},
	nil,
)

var createUserEndHandler = insertEndpoint(
	"userends",
	func() interface{} { return &UserEnd{} },
	[]middleware.Middleware{setUserID},
	[]middleware.Middleware{
		func(fn httprouter.Handle) httprouter.Handle {
			return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				fn(w, r, p)
			}
		},
	},
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

var createFeedMediaHandler = insertEndpoint(
	"feedmedias",
	func() interface{} { return &FeedMedia{} },
	nil,
	nil,
)

var createPlantSharingHandler = insertEndpoint(
	"plantsharings",
	func() interface{} { return &PlantSharing{} },
	[]middleware.Middleware{setUserID},
	nil,
)
