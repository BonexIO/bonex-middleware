package api

import (
    "bonex-middleware/log"
    "bonex-middleware/services/api/response"
    "net/http"
    "encoding/json"
    "bonex-middleware/types"
    "github.com/Nargott/goutils"
)

func (this *api) requestMoney(w http.ResponseWriter, r *http.Request) {
    var qi types.QueueItem
    err := json.NewDecoder(r.Body).Decode(&qi)
    if err != nil {
        log.Warn(err.Error())
        response.JsonError(w, types.NewError(types.ErrBadRequest))
        return
    }

    err = this.faucet.AddToQueue(&qi, goutils.GetClearIpAddress(r))
    if err != nil {
        log.Error(err.Error())
        response.JsonError(w, types.NewError(types.ErrService))
        return
    }
}
