package api

import (
	"context"
	"errors"
	"fmt"
	"husir_blog/src/common"
	"husir_blog/src/db"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/maps"
)

type createPostReq struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Photo       string   `json:"photo" binding:"required"`
	Author      string   `json:"author" binding:"required"`
	Categories  []string `json:"categories" binding:"required"`
	Status      string   `json:"status" binding:"required"`
}

func (server *Server) createPost(ctx *gin.Context) {
	var req createPostReq

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		return
	}

	authorId, err := primitive.ObjectIDFromHex(req.Author)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidObjectId(err))
		return
	}

	postStatus, ok := db.AllPostStatus[req.Status]
	if !ok {
		err = errors.New("status must be: " + strings.Join(maps.Values(db.AllPostStatus), " or "))
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidPostStatus(err))
		return
	}

	categoriesId := make([]*primitive.ObjectID, len(req.Categories))

	for index, v := range req.Categories {
		id, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, common.ErrInvalidObjectId(err))
			return
		} else {
			categoriesId[index] = &id
		}
	}

	arg := db.CreatePostParams{
		Title:       req.Title,
		Description: req.Description,
		Photo:       req.Photo,
		Author:      authorId,
		Categories:  categoriesId,
		Status:      postStatus,
		Slug:        slug.Make(req.Title),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := server.store.CreatePost(ctx, &arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	result, err := server.store.GetPostById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) getPostBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(errors.New("slug is empty")))
		return
	}

	post, err := server.rdb.GetPostBySlug(ctx, slug)
	if err != nil {
		log.Printf("can not get post: %s from redist: %s", slug, err)
	} else {
		ctx.JSON(http.StatusOK, post)
		return
	}

	res, err := server.store.GetPostBySlug(ctx, slug)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrCannotGetEntity(db.PostCollection, err))
		return
	}

	server.rdb.SetPostBySlug(ctx, res)

	ctx.JSON(http.StatusOK, res)
}
func (server *Server) getPostByCategory(ctx *gin.Context) {
	categorySlug := ctx.Query("category_slug")
	if categorySlug == "" {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(errors.New("category slug id is empty")))
		return
	}

	category, err := server.store.GetCategoryBySlug(ctx, categorySlug)
	if err != nil {
		ctx.JSON(http.StatusNotFound, common.ErrEntityNotFound("Category", err))
		return
	}

	res, err := server.store.GetPostByCategory(ctx, category.ID.Hex())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrCannotGetEntity(db.PostCollection, err))
		return
	}
	ctx.JSON(http.StatusOK, res)
}
func (server *Server) ListPost(ctx *gin.Context) {
	var filter db.PostFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		return
	}

	filter.Default()

	_, ok := db.AllPostStatus[filter.Status]
	if !ok {
		err := errors.New("status must be: " + strings.Join(maps.Values(db.AllPostStatus), " or "))
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidPostStatus(err))
		return
	}

	var paging db.Paging
	if err := ctx.ShouldBindQuery(&paging); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		return
	}

	paging.Default()

	key := fmt.Sprintf("page%d%s%s", paging.Page, filter.Status, filter.CategoryId)
	fmt.Println(key)

	list := server.rdb.GetListPost(ctx, key, paging)
	if list != nil {
		ctx.JSON(http.StatusOK, common.NewSuccessResponse(list, paging, filter))
		return
	}

	result, err := server.store.ListPost(ctx, &filter, &paging)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	server.rdb.SetListPost(ctx, key, result)

	paging.Total = len(result)

	ctx.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))
}
func (server *Server) ListPostRandom(ctx *gin.Context) {

	n, err := strconv.Atoi(ctx.Query("n"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		return
	}
	if n <= 0 {
		n = 4
	}

	val, err := server.rdb.Client.HRandFieldWithValues(ctx, server.cfg.KeyPostRecommend, n).Result()
	if err != nil {
		log.Printf("key:%s does not exists. Error: %s", server.cfg.KeyPostRecommend, err)
	} else {

		var list []*db.PostRecommendResponse

		for _, v := range val {
			list = append(list, &db.PostRecommendResponse{Title: v.Value, Slug: v.Key})
		}

		if len(list) == n {
			log.Println(">> Post recommend get from redis")
			ctx.JSON(http.StatusOK, list)
			return
		}

	}

	var filter db.PostFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		return
	}

	filter.Default()

	result, err := server.store.ListPostRandom(ctx, &filter, n)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) CachePostRecommend(ctx context.Context) error {
	post, err := server.store.ListPostRandom(ctx, &db.PostFilter{Status: string(db.PostStatusVisibility)}, 100)
	if err != nil {
		return err
	}

	map2 := make(map[string]string)
	for _, value := range post {
		map2[value.Slug] = value.Title
	}

	res, err := server.rdb.Client.HSet(ctx, server.cfg.KeyPostRecommend, map2).Result()
	if err != nil {
		return err
	}
	log.Print(res)

	return nil
}

func (server *Server) ListenPostChange(ctx context.Context) error {

	chanelName := "post_recom_ch"
	res := server.rdb.Client.Subscribe(ctx, chanelName)
	defer res.Close()

	message, err := res.ReceiveMessage(ctx)
	if err != nil {
		log.Printf("Receive Message of chanel %s error: %s", chanelName, err)
		return err
	}

	switch message.Payload {
	case "product:create":

		err = server.CachePostRecommend(ctx)
		if err != nil {
			log.Println(err)
		}

	}

	return nil
}
