package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/haniifac/simplebank/db/sqlc"
	"github.com/haniifac/simplebank/util"
	"github.com/lib/pq"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username" binding:"required,alphanum"`
	Fullname string `json:"fullname" binding:"required"`
}

type CreateUserResponse struct {
	Email             string    `json:"email"`
	Username          string    `json:"username"`
	Fullname          string    `json:"fullname"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordUpdatedAt time.Time `json:"password_updated_at"`
}

type GetUserRequest struct {
	Username string `uri:"username" binding:"required"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		if errName, ok := err.(*pq.Error); ok {
			switch errName.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
		Username:       req.Username,
		Fullname:       req.Fullname,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	res := CreateUserResponse{
		Email:             user.Email,
		Username:          user.Username,
		Fullname:          user.Fullname,
		CreatedAt:         user.CreatedAt,
		PasswordUpdatedAt: user.PasswordUpdatedAt,
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetUser(ctx *gin.Context) {
	var req GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(http.StatusNotFound, errResponse(err))
		default:
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, user)
}
