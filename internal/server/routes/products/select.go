package products

import (
	"encoding/json"
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

type searchProductsResult struct {
	Products []db.Products `json:"products"`
}

func searchProducts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	terms := r.URL.Query().Get("terms")
	if terms == "" {
		http.Error(w, "missing 'terms' parameter", http.StatusInternalServerError)
		return
	}
	sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
	selector := sess.Select("p.*")
	selector = selector.From("products p")
	selector = selector.Where("p.name ilike '%' || ? || '%'", terms)
	products := []db.Products{}
	if err := selector.All(&products); err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(searchProductsResult{products}); err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
