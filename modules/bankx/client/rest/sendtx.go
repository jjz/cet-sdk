package rest

import (
	"fmt"
	"github.com/coinexchain/dex/modules/bankx"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// SendReq defines the properties of a send request's body.
type SendReq struct {
	BaseReq    rest.BaseReq `json:"base_req"`
	Amount     sdk.Coins    `json:"amount"`
	UnlockTime int64        `json:"unlock_time"`
}

// SendRequestHandlerFn - http request handler to send coins to a address.
func SendTxRequestHandlerFn(cdc *codec.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		toAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req SendReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		currentTime := time.Now().Unix()
		if req.UnlockTime < currentTime {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("unlock time should be later than the current time").Error())
			return
		}
		msg := bankx.NewMsgSend(fromAddr, toAddr, req.Amount, req.UnlockTime)
		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
