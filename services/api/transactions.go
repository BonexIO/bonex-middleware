package api

import (
	dmodels "bonex-middleware/dao/models"
	"bonex-middleware/log"
	"bonex-middleware/models"
	"bonex-middleware/services/api/handler"
	"bonex-middleware/services/api/response"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

func (this *api) createTransaction(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		IssuerPubKey string `json:"pubkey"`
		AssetCode    string `json:"asset_code"`
		Secret       string `json:"secret"`
		Amount       uint64 `json:"amount"`
	}

	var params reqParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		response.JsonError(w, models.NewError(models.ErrBadRequest))
		return
	}

	if len(params.IssuerPubKey) != addressLen {
		response.JsonError(w, models.NewError(models.ErrBadParam, "pubkey"))
		return
	}

	if params.AssetCode == "" {
		response.JsonError(w, models.NewError(models.ErrBadParam, "asset_code"))
		return
	}

	if len(params.Secret) < 5 {
		response.JsonError(w, models.NewError(models.ErrBadParam, "secret"))
		return
	}

	transaction := &dmodels.Transaction{
		Pubkey:    params.IssuerPubKey,
		AssetCode: params.AssetCode,
		Amount:    params.Amount,
		Secret:    params.Secret,
	}

	err = this.dao.CreateTransaction(transaction, nil)
	if err != nil {
		log.Errorf("CreateMerchant: %s", err.Error())
		response.JsonError(w, models.NewError(models.ErrService))
		return
	}

	response.Json(w, "ok")
}

func (this *api) redeemTransaction(w http.ResponseWriter, r *http.Request) {

	middlewarePrivateKey := os.Getenv("MIDDLEWARE_PRIVATE_KEY")

	type reqParams struct {
		IssuerPubKey string `json:"issuer_pubkey"`
		AssetCode    string `json:"asset_code"`
		Secret       string `json:"secret"`
		DestPubKey   string `json:"dest_pubkey"`
	}

	var params reqParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		response.JsonError(w, models.NewError(models.ErrBadRequest))
		return
	}

	if len(params.IssuerPubKey) != addressLen {
		response.JsonError(w, models.NewError(models.ErrBadParam, "issuer_pubkey"))
		return
	}

	if len(params.DestPubKey) != addressLen {
		response.JsonError(w, models.NewError(models.ErrBadParam, "dest_pubkey"))
		return
	}

	if params.AssetCode == "" {
		response.JsonError(w, models.NewError(models.ErrBadParam, "asset_code"))
		return
	}

	if len(params.Secret) < 16 {
		response.JsonError(w, models.NewError(models.ErrBadParam, "secret"))
		return
	}

	tx, err := this.dao.RedeemTransaction(params.IssuerPubKey, params.AssetCode, params.Secret, nil)

	if err != nil {
		log.Errorf("RedeemTransaction: %s", err.Error())
		response.JsonError(w, models.NewError(models.ErrService))
		return
	}

	if tx == nil {
		response.JsonError(w, models.NewError(models.ErrNotFound))
		return
	}

	err = handler.SendAssetTx(middlewarePrivateKey, params.DestPubKey, tx.AssetCode, strconv.FormatUint(tx.Amount, 10))

	if err != nil {
		log.Errorf("RedeemTransaction: %s", err.Error())
		return
	}

	defer response.Json(w, "ok")

}
