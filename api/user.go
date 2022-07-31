package api

import (
	"fmt"
	"net/http"
	"time"

	db "github.com/dpsigor/cheatsheet-golang-postgres/db/sqlc"
	"github.com/dpsigor/cheatsheet-golang-postgres/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserRes struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPw, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: hashedPw,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			fmt.Printf("pgErr.Constraint = %+v\n", pgErr.Constraint)
			switch pgErr.Code.Name() {
			case "unique_violation":
				switch pgErr.Constraint {
				case "users_email_key":
					ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("email '%s' already exists", req.Email)))
					return
				case "users_pkey":
					ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("username '%s' already exists", req.Username)))
					return
				}
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createUserRes{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
