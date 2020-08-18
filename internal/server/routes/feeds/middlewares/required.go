package middlewares

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// UserEndIDRequired - Checks if the request has a userEndID
func UserEndIDRequired(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ueid := r.Context().Value(UserEndIDContextKey{})
		if ueid == nil {
			logrus.Errorln("Missing userEndID")
			http.Error(w, "Missing userEndID", http.StatusBadRequest)
			return
		}
		fn(w, r, p)
	}
}
