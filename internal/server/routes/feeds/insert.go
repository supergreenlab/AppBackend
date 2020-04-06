package feeds

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/spf13/viper"
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
				hmacSampleSecret := []byte(viper.GetString("JWTSecret"))
				id := r.Context().Value(insertedIDContextKey{}).(uuid.UUID)
				uid := r.Context().Value(userIDContextKey{}).(uuid.UUID)

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"userID":    uid.String(),
					"userEndID": id.String(),
				})
				tokenString, err := token.SignedString(hmacSampleSecret)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("x-sgl-token", tokenString)

				fn(w, r, p)
			}
		},
	},
)

var createBoxHandler = insertEndpoint(
	"boxes",
	func() interface{} { return &Box{} },
	[]middleware.Middleware{
		setUserID,
		checkAccessRight("devices", "DeviceID", true, func() interface{} { return &Device{} }),
	},
	[]middleware.Middleware{
		createUserEndObjects("userend_devices", "DeviceID", func() interface{} { return &UserEndDevice{} }),
	},
)

var createPlantHandler = insertEndpoint(
	"plants",
	func() interface{} { return &Plant{} },
	[]middleware.Middleware{
		setUserID,
		checkAccessRight("boxes", "BoxID", false, func() interface{} { return &Box{} }),
	},
	[]middleware.Middleware{
		createUserEndObjects("userend_plants", "PlantID", func() interface{} { return &UserEndPlant{} }),
	},
)

var createTimelapseHandler = insertEndpoint(
	"timelapses",
	func() interface{} { return &Timelapse{} },
	[]middleware.Middleware{
		setUserID,
		checkAccessRight("plants", "PlantID", false, func() interface{} { return &Plant{} }),
	},
	[]middleware.Middleware{
		createUserEndObjects("userend_timelapses", "TimelapseID", func() interface{} { return &UserEndTimelapse{} }),
	},
)

var createDeviceHandler = insertEndpoint(
	"devices",
	func() interface{} { return &Device{} },
	[]middleware.Middleware{setUserID},
	[]middleware.Middleware{
		createUserEndObjects("userend_devices", "DeviceID", func() interface{} { return &UserEndDevice{} }),
	},
)

var createFeedHandler = insertEndpoint(
	"feeds",
	func() interface{} { return &Feed{} },
	[]middleware.Middleware{setUserID},
	[]middleware.Middleware{
		createUserEndObjects("userend_feeds", "FeedID", func() interface{} { return &UserEndFeed{} }),
	},
)

var createFeedEntryHandler = insertEndpoint(
	"feedentries",
	func() interface{} { return &FeedEntry{} },
	[]middleware.Middleware{
		setUserID,
		checkAccessRight("feeds", "FeedID", false, func() interface{} { return &Feed{} }),
	},
	[]middleware.Middleware{
		createUserEndObjects("userend_feedentries", "FeedEntryID", func() interface{} { return &UserEndFeedEntry{} }),
	},
)

var createFeedMediaHandler = insertEndpoint(
	"feedmedias",
	func() interface{} { return &FeedMedia{} },
	[]middleware.Middleware{
		setUserID,
		checkAccessRight("feedentries", "FeedEntryID", false, func() interface{} { return &FeedEntry{} }),
	},
	[]middleware.Middleware{
		createUserEndObjects("userend_feedmedias", "FeedMediaID", func() interface{} { return &UserEndFeedMedia{} }),
	},
)

var createPlantSharingHandler = insertEndpoint(
	"plantsharings",
	func() interface{} { return &PlantSharing{} },
	[]middleware.Middleware{
		setUserID,
		checkAccessRight("feedentries", "FeedEntryID", false, func() interface{} { return &FeedEntry{} }),
	},
	nil,
)
