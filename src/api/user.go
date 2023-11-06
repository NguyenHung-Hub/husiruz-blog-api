package api

import (
	"fmt"
	"husir_blog/src/common"
	"husir_blog/src/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userResponse struct {
	ID          primitive.ObjectID    `json:"_id" bson:"_id"`
	Username    string                `json:"username" bson:"username"`
	Email       string                `json:"email" bson:"email"`
	Avatar      string                `json:"avatar" bson:"avatar,omitempty"`
	Description string                `json:"description" bson:"description,omitempty"`
	Posts       *[]primitive.ObjectID `json:"posts,omitempty" bson:"posts,omitempty"`
	CreatedAt   time.Time             `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt" bson:"updatedAt"`
}

type createUserReq struct {
	Username    string                `json:"username" binding:"required"`
	Email       string                `json:"email" binding:"required"`
	Password    string                `json:"password" binding:"required"`
	Avatar      string                `json:"avatar" binding:"required"`
	Description string                `json:"description" binding:"required"`
	Posts       *[]primitive.ObjectID `json:"-" bson:"posts,omitempty"`
	CreatedAt   time.Time             `json:"-" bson:"createdAt"`
	UpdatedAt   time.Time             `json:"-" bson:"updatedAt"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRes struct {
	SessionId   string       `json:"session_id"`
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

type RenewAccessTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Avatar:      user.Avatar,
		Description: user.Description,
		Posts:       user.Posts,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		return
	}

	fmt.Println(req)

	if req.Posts == nil {
		req.Posts = &[]primitive.ObjectID{}
	}

	arg := db.CreateUserParams{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		Avatar:      req.Avatar,
		Description: req.Description,
		Posts:       &[]primitive.ObjectID{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := server.store.CreateUser(ctx, &arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	result, err := server.store.GetUserById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(*result))
}

func (server *Server) getUserById(ctx *gin.Context) {
	id := ctx.Param("id")

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		return
	}

	user, err := server.store.GetUserById(ctx, &userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrCannotGetEntity(db.UserCollection, err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (server *Server) login(ctx *gin.Context) {
	var req LoginReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		return
	}
	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrEntityNotFound(db.UserCollection, err))
		return
	}

	accessToken, _, err := server.token.CreateToken(user.Username, server.cfg.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.token.CreateToken(user.Username, server.cfg.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	// fmt.Println("refreshPayload: ", refreshPayload.ID)

	session, err := server.store.CreateSession(ctx, &db.CreateSessionParams{
		SessionID:    refreshPayload.ID.String(),
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpireAt:     refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	res := LoginRes{
		SessionId:   session.SessionID,
		AccessToken: accessToken,
		User:        newUserResponse(*user),
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) checkSession(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "session id is empty"})
		return
	}

	session, err := server.store.GetSessionById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}

	now := time.Now()

	if now.After(session.ExpireAt) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "session expired"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "ok"})
}
