package api

import (
	"husir_blog/src/common"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	// RefreshToken string `json:"refresh_token" binding:"required"`
	SessionId string `json:"session_id" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

type checkTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		return
	}

	session, err := server.store.GetSessionById(ctx, req.SessionId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, common.ErrUnauthorized(err))
		return
	}

	now := time.Now()
	if now.After(session.ExpireAt) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "session expired"})
		return
	}

	refreshPayload, err := server.token.Verify(session.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	accessToken, accessPayload, err := server.token.CreateToken(
		refreshPayload.Username,
		server.cfg.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	res := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, res)

}
func (server *Server) checkRefreshToken(ctx *gin.Context) {
	var req checkTokenReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	_, err := server.token.Verify(req.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "session expired"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "ok"})
}
