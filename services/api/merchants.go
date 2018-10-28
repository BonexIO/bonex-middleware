package api

import (
    dmodels "bonex-middleware/dao/models"
    "bonex-middleware/log"
    "bonex-middleware/models"
    "bonex-middleware/services/api/response"
    "encoding/base64"
    "encoding/json"
    "github.com/gorilla/mux"
    "io/ioutil"
    "net/http"
)

const addressLen = 56

// Index returns the service name in plaintext.
func (this *api) getMerchant(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    a := params["address"]
    if len(a) != addressLen {
        response.JsonError(w, models.NewError(models.ErrBadParam, "address"))
        return
    }

    merchant, err := this.dao.GetMerchantByPubkeyOrNil(a, nil)
    if err != nil {
        log.Errorf("GetMerchantByPubkey: %s", err.Error())
        response.JsonError(w, models.NewError(models.ErrService))
        return
    }

    if merchant == nil {
        response.JsonError(w, models.NewError(models.ErrNotFound))
        return
    }

    response.Json(w, map[string]interface{}{
        "title":      merchant.Title,
        "pubkey":     merchant.Pubkey,
        "asset_code": merchant.AssetCode,
        "logo":       merchant.Logo,
        "created_at": merchant.CreatedAt.Time.Unix(),
    })
}

func (this *api) listMerchants(w http.ResponseWriter, r *http.Request) {
    merchants, err := this.dao.GetMerchants()
    if err != nil {
        log.Errorf("GetMerchantByPubkey: %s", err.Error())
        response.JsonError(w, models.NewError(models.ErrService))
        return
    }

    var result []map[string]interface{}
    if merchants != nil {
        for _, merchant := range merchants {
            result = append(result, map[string]interface{}{
                "title":      merchant.Title,
                "pubkey":     merchant.Pubkey,
                "asset_code": merchant.AssetCode,
                "logo":       merchant.Logo,
            })
        }
    }

    response.Json(w, result)
}

// Index returns the service name in plaintext.
func (this *api) createMerchant(w http.ResponseWriter, r *http.Request) {
    type reqParams struct {
        Title     string `json:"title"`
        Pubkey    string `json:"pubkey"`
        AssetCode string `json:"asset_code"`
        Logo      string `json:"logo"`
    }

    var params reqParams
    err := json.NewDecoder(r.Body).Decode(&params)
    if err != nil {
        response.JsonError(w, models.NewError(models.ErrBadRequest))
        return
    }

    if params.Title == "" {
        response.JsonError(w, models.NewError(models.ErrBadParam, "title"))
        return
    }

    if len(params.Pubkey) != addressLen {
        response.JsonError(w, models.NewError(models.ErrBadParam, "pubkey"))
        return
    }

    if params.AssetCode == "" {
        response.JsonError(w, models.NewError(models.ErrBadParam, "asset_code"))
        return
    }

    if params.Logo == "" {
        response.JsonError(w, models.NewError(models.ErrBadParam, "logo"))
        return
    }

    tmpfile, err := ioutil.TempFile("/opt/images", "logo")
    if err != nil {
        log.Errorf("Cannot create temp file: %s", err.Error())
        response.JsonError(w, models.NewError(models.ErrService))
        return
    }

    decodedLogo, err := base64.StdEncoding.DecodeString(params.Logo)
    if err != nil {
        log.Errorf("Cannot DecodeString: %s", err.Error())
        response.JsonError(w, models.NewError(models.ErrService))
        return
    }

    if _, err := tmpfile.Write(decodedLogo); err != nil {
        log.Errorf("Cannot write temp file: %s", err.Error())
        response.JsonError(w, models.NewError(models.ErrService))
        return
    }

    if err := tmpfile.Close(); err != nil {
        log.Errorf("Cannot close temp file: %s", err.Error())
        response.JsonError(w, models.NewError(models.ErrService))
        return
    }

    merchant := &dmodels.Merchant{
        Title:     params.Title,
        Pubkey:    params.Pubkey,
        AssetCode: params.AssetCode,
        Logo:      tmpfile.Name(),
    }

    err = this.dao.CreateMerchant(merchant, nil)
    if err != nil {
        log.Errorf("CreateMerchant: %s", err.Error())
        response.JsonError(w, models.NewError(models.ErrService))
        return
    }

    response.Json(w, "ok")
}
