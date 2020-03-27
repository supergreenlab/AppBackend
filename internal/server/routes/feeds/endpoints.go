package feeds

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func outputObjectID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := r.Context().Value(insertedIDContextKey{})
	response := struct {
		ID string `json:"id"`
	}{id.(string)}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
