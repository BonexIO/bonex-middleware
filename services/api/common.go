package api

import (
    "bonex-middleware/services/api/response"
    "net/http"
)

// Index returns the service name in plaintext.
func (this *api) index(w http.ResponseWriter, r *http.Request) {
    response.Json(w, "bonex-middleware")
}
