package api

import (
    "bonex-middleware/dao/models"
    "bonex-middleware/log"
    "bonex-middleware/services/api/response"
    "encoding/json"
    "github.com/gorilla/mux"
    "net/http"
    "bonex-middleware/types"
)

func (this *api) subscribe(w http.ResponseWriter, r *http.Request) {
    type reqParams struct {
        Account  string `json:"account"`
        Merchant string `json:"merchant"`
    }

    var params reqParams
    err := json.NewDecoder(r.Body).Decode(&params)
    if err != nil {
        response.JsonError(w, types.NewError(types.ErrBadRequest))
        return
    }

    if len(params.Account) != addressLen {
        response.JsonError(w, types.NewError(types.ErrBadParam, "account"))
        return
    }

    if len(params.Merchant) != addressLen {
        response.JsonError(w, types.NewError(types.ErrBadParam, "merchant"))
        return
    }

    account, err := this.dao.GetAccountByPubkeyOrNil(params.Account)
    if err != nil {
        response.JsonError(w, types.NewError(types.ErrService))
        return
    }

    if account == nil {
        // Create acc in DB
        account = &models.Account{
            Pubkey: params.Account,
        }

        err = this.dao.CreateAccount(account)
        if err != nil {
            response.JsonError(w, types.NewError(types.ErrService))
            return
        }
    }

    merchant, err := this.dao.GetMerchantByPubkeyOrNil(params.Merchant, nil)
    if err != nil {
        response.JsonError(w, types.NewError(types.ErrService))
        return
    }

    if merchant == nil {
        response.JsonError(w, types.NewError(types.ErrNotFound))
        return
    }

    err = this.dao.Subscribe(account.Id, merchant.Id)
    if err != nil {
        response.JsonError(w, types.NewError(types.ErrService))
        return
    }

    response.Json(w, map[string]interface{}{
        "account":  account.Pubkey,
        "merchant": merchant.Pubkey,
    })
}

// Index returns the service name in plaintext.
func (this *api) unsubscribe(w http.ResponseWriter, r *http.Request) {
    type reqParams struct {
        Account  string `json:"account"`
        Merchant string `json:"merchant"`
    }

    var params reqParams
    err := json.NewDecoder(r.Body).Decode(&params)
    if err != nil {
        response.JsonError(w, types.NewError(types.ErrBadRequest))
        return
    }

    if len(params.Account) != addressLen {
        response.JsonError(w, types.NewError(types.ErrBadParam, "account"))
        return
    }

    if len(params.Merchant) != addressLen {
        response.JsonError(w, types.NewError(types.ErrBadParam, "merchant"))
        return
    }

    account, err := this.dao.GetAccountByPubkeyOrNil(params.Account)
    if err != nil {
        response.JsonError(w, types.NewError(types.ErrService))
        return
    }

    if account == nil {
        response.JsonError(w, types.NewError(types.ErrNotFound))
        return
    }

    merchant, err := this.dao.GetMerchantByPubkeyOrNil(params.Merchant, nil)
    if err != nil {
        response.JsonError(w, types.NewError(types.ErrService))
        return
    }

    if merchant == nil {
        response.JsonError(w, types.NewError(types.ErrNotFound))
        return
    }

    err = this.dao.Unsubscribe(account.Id, merchant.Id)
    if err != nil {
        response.JsonError(w, types.NewError(types.ErrService))
        return
    }

    response.Json(w, map[string]interface{}{
        "account":  account.Pubkey,
        "merchant": merchant.Pubkey,
    })
}

func (this *api) getSubscriptions(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    a := params["address"]
    if len(a) != addressLen {
        response.JsonError(w, types.NewError(types.ErrBadParam, "address"))
        return
    }

    account, err := this.dao.GetAccountByPubkeyOrNil(a)
    if err != nil {
        response.JsonError(w, types.NewError(types.ErrService))
        return
    }

    var result []map[string]interface{}

    if account != nil {
        merchants, err := this.dao.GetSubscriptions(account.Id)
        if err != nil {
            log.Errorf("GetSubscriptions: %s", err.Error())
            response.JsonError(w, types.NewError(types.ErrService))
            return
        }

        if merchants != nil {
            for _, merchant := range merchants {
                result = append(result, map[string]interface{}{
                    "title":      merchant.Title,
                    "pubkey":     merchant.Pubkey,
                    "asset_code": merchant.AssetCode,
                    "created_at": merchant.CreatedAt.Time.Unix(),
                })
            }
        }
    }

    response.Json(w, result)
}

func (this *api) getSubscribers(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    a := params["address"]
    if len(a) != addressLen {
        response.JsonError(w, types.NewError(types.ErrBadParam, "address"))
        return
    }

    merchant, err := this.dao.GetMerchantByPubkeyOrNil(a, nil)
    if err != nil {
        response.JsonError(w, types.NewError(types.ErrService))
        return
    }

    var result []map[string]interface{}
    if merchant != nil {
        accounts, err := this.dao.GetSubscribers(merchant.Id)
        if err != nil {
            log.Errorf("GetSubscribers: %s", err.Error())
            response.JsonError(w, types.NewError(types.ErrService))
            return
        }

        if accounts != nil {
            for _, account := range accounts {
                result = append(result, map[string]interface{}{
                    "pubkey": account.Pubkey,
                })
            }
        }
    }

    response.Json(w, result)
}
