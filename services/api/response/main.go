package response

import (
	"encoding/json"
	"bonex-middleware/models"
	"net/http"
)

// Json writes to ResponseWriter a single JSON-object
func Json(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// JsonError writes to ResponseWriter error
func JsonError(w http.ResponseWriter, err error) {
	var e *models.Error
	var ok bool

	if e, ok = err.(*models.Error); !ok {
		e = models.FromError(err)
	}

	js, _ := json.Marshal(e.ToMap())
	w.WriteHeader(e.GetHttpCode())
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
