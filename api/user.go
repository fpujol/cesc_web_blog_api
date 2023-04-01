package api

import (
	"database/sql"
	"net/http"

	"blogapi/api/request"
	"blogapi/api/response"
	db "blogapi/db/sqlc"
	"blogapi/pkg/password"

	"github.com/gin-gonic/gin"
)

func (s *Server) Me(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"name": "Cesc",
	})
}

func (server *Server) loginUser(ctx *gin.Context) {
    var req request.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = password.CheckPassword(req.Email, req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Email,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Email,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Email:     user.Email,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("gin_cookie",accessToken, 60*60*24, "/","localhost",true, true)

	rsp := response.LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  mapUserToResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) logout(ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("gin_cookie","", -1, "/","localhost",true, true)
	ctx.JSON(http.StatusOK, gin.H{
		"logout": true,
	})
}
