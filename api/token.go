package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req RenewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	fmt.Println("RenewAccessToken payload:", payload.ID)

	session, err := server.store.GetSession(ctx, payload.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(http.StatusNotFound, errResponse(err))
		default:
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
		}
		return
	}

	if session.IsBlocked {
		err = errors.New("session is blocked")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	if session.Username != payload.Username {
		err = errors.New("session username does not match payload username")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	newAccessToken, _, err := server.tokenMaker.CreateToken(payload.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, RenewAccessTokenResponse{
		AccessToken: newAccessToken,
	})
}
