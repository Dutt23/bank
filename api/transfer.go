package api

import (
	"database/sql"
	"errors"
	"fmt"
	db "github/dutt23/bank/db/sqlc"
	"github/dutt23/bank/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transferAccountRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required",min=1`
	ToAccountID   int64  `json:"to_account_id" binding:"required",min=1`
	Amount        int64  `json:"amount" binding:"required",gt=0`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, isValid := server.validAccount(ctx, req.FromAccountID, req.Currency)

	if !isValid {
		return
	}

	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)

	if authPayload.Username != fromAccount.Owner {
		err := errors.New("from account doesn't belong to the user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, isValidTo := server.validAccount(ctx, req.ToAccountID, req.Currency)

	if !isValidTo {
		return
	}

	arg := db.TransferTxParams{
		FromAccountId: req.FromAccountID,
		ToAccountId:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)

	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] current mismatch: %s vs %s", account.ID, currency, account.Currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
