package feeds

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func outputObjectID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := r.Context().Value(insertedIDContextKey{}).(uuid.UUID)
	response := struct {
		ID string `json:"id"`
	}{id.String()}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
