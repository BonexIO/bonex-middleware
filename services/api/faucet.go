package api

import (
    "bonex-middleware/log"
    "bonex-middleware/services/api/response"
    "net/http"
    "bonex-middleware/types"
    "github.com/Nargott/goutils"
    "github.com/gorilla/mux"
    "github.com/wedancedalot/decimal"
)

func (this *api) requestMoney(w http.ResponseWriter, r *http.Request) {
    var (
        err error
        qi types.QueueItem
    )

    params := mux.Vars(r)
    qi.Address = params["address"]
    if len(qi.Address) != addressLen {
        response.JsonError(w, types.NewError(types.ErrBadParam, "address"))
        return
    }

    amountStr := r.URL.Query().Get("amount")
    if amountStr == "" {
        qi.Amount = this.config.Faucet.MaxAllowed24HoursValue
    } else {
        qi.Amount, err = decimal.NewFromString(amountStr)
        if err != nil {
            response.JsonError(w, types.NewError(types.ErrBadParam, "amount"))
            return
        }
    }

    err = this.faucet.AddToQueue(&qi, goutils.GetClearIpAddress(r))
    if err != nil {
        log.Error(err.Error())
        response.JsonError(w, err)
        return
    }

    response.Json(w, map[string]interface{}{
        "ok":  true,
    })
    return
}
