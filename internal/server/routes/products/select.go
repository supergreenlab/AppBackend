package products

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func searchProducts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
}
