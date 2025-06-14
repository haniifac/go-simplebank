package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/haniifac/simplebank/db/sqlc"
	"github.com/haniifac/simplebank/token"
)

type TransferRequestParams struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req TransferRequestParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if fromAccount.Owner != authPayload.Username {
		err := fmt.Errorf("from account %d does not belong to the authenticated user %s", req.FromAccountID, authPayload.Username)
		ctx.JSON(http.StatusForbidden, errResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			err := fmt.Errorf("account %d not found", accountID)
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account %d currency mismatch: %s with %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))

		return account, false
	}

	return account, true
}
