package feeds

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"upper.io/db.v3/lib/sqlbuilder"
)

type loginParams struct {
	Handle   string `json:"handle"`
	Password string `json:"password"`
}

func loginHandler() httprouter.Handle {
	s := middleware.NewStack()

	s.Use(decodeJSON(func() interface{} { return &loginParams{} }))

	return s.Wrap(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		hmacSampleSecret := []byte(viper.GetString("JWTSecret"))
		lp := r.Context().Value(objectContextKey{}).(*loginParams)
		sess := r.Context().Value(sessContextKey{}).(sqlbuilder.Database)

		u := User{}
		err := sess.Select("id", "password").From("users").Where("nickname = ?", lp.Handle).One(&u)
		if err != nil {
			logrus.Errorln(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(lp.Password))
		if err != nil {
			logrus.Errorln(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID": u.ID.UUID.String(),
		})
		tokenString, err := token.SignedString(hmacSampleSecret)
		if err != nil {
			logrus.Errorln(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("x-sgl-token", tokenString)

		w.WriteHeader(http.StatusOK)
	})
}
