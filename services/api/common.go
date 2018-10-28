package api

import (
    "bonex-middleware/log"
    "bonex-middleware/models"
    "bonex-middleware/services/api/response"
    "github.com/gorilla/mux"
    "io"
    "net/http"
    "os"
)

// Index returns the service name in plaintext.
func (this *api) index(w http.ResponseWriter, r *http.Request) {
    response.Json(w, "bonex-middleware")
}

func (this *api) getImage(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    address := params["address"]
    if address == "" {
        response.JsonError(w, models.NewError(models.ErrBadParam, "address"))
        return
    }

    merchant, err := this.dao.GetMerchantByPubkeyOrNil(address, nil)
    if err != nil {
        log.Errorf("GetMerchantByPubkey: %s", err.Error())
        response.JsonError(w, models.NewError(models.ErrService))
        return
    }

    if merchant == nil {
        response.JsonError(w, models.NewError(models.ErrNotFound))
        return
    }

    data, err := os.Open(merchant.Logo)
    if err != nil {
        log.Errorf("Readfile: %s", err.Error())
        response.JsonError(w, models.NewError(models.ErrService))
        return
    }

    io.Copy(w, data)
}
