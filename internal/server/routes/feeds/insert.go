package feeds

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	uuid "github.com/satori/go.uuid"
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
				ue := r.Context().Value(objectContextKey{}).(*UserEnd)
				uid := r.Context().Value(userIDContextKey{}).(uuid.UUID)

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"userID":    uid.String(),
					"userEndID": ue.ID.UUID.String(),
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
	nil,
)

var createPlantHandler = insertEndpoint(
	"plants",
	func() interface{} { return &Plant{} },
	[]middleware.Middleware{
		setUserID,
		checkAccessRight("boxes", "BoxID", false, func() interface{} { return &Box{} }),
	},
	nil,
)

var createTimelapseHandler = insertEndpoint(
	"timelapses",
	func() interface{} { return &Timelapse{} },
	[]middleware.Middleware{
		setUserID,
		checkAccessRight("plants", "PlantID", false, func() interface{} { return &Plant{} }),
	},
	nil,
)

var createDeviceHandler = insertEndpoint(
	"devices",
	func() interface{} { return &Device{} },
	[]middleware.Middleware{setUserID},
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
	[]middleware.Middleware{
		setUserID,
		checkAccessRight("feeds", "FeedID", false, func() interface{} { return &Feed{} }),
	},
	nil,
)

var createFeedMediaHandler = insertEndpoint(
	"feedmedias",
	func() interface{} { return &FeedMedia{} },
	[]middleware.Middleware{
		setUserID,
		checkAccessRight("feedentries", "FeedEntryID", false, func() interface{} { return &FeedEntry{} }),
	},
	nil,
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
