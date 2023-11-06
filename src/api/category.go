package api

import (
	"errors"
	"fmt"
	"husir_blog/src/common"
	"husir_blog/src/db"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createCategoryReq struct {
	Name string `json:"name" binding:"required"`
}

type createCategoryRes struct {
	ID   primitive.ObjectID `json:"_id"`
	Name string             `json:"name"`
	Slug string             `json:"slug"`
}

func (server *Server) createCategory(ctx *gin.Context) {
	var req createCategoryReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		return
	}

	arg := db.CreateCategoryParams{
		Name:      req.Name,
		Slug:      slug.Make(req.Name),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := server.store.CreateCategory(ctx, &arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	result, err := server.store.GetCategoryById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	server.rdb.DeleteKey(ctx, "all_categories")
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(createCategoryRes{
		ID: result.ID, Name: result.Name, Slug: result.Slug}))
}

func (server *Server) getCategoryById(ctx *gin.Context) {
	id := ctx.Param("id")

	categoryId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
		return
	}

	category, err := server.store.GetCategoryById(ctx, &categoryId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrCannotGetEntity(db.CategoryCollection, err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (server *Server) getCategoriesName(ctx *gin.Context) {

	res, err := server.rdb.GetCategories(ctx, "all_categories")
	if err == nil {
		ctx.JSON(http.StatusOK, res)
		return
	}

	categories, err := server.store.GetCategoriesName(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrCannotGetEntity(db.CategoryCollection, err))
		return
	}

	server.rdb.SetCategories(ctx, "all_categories", categories)

	ctx.JSON(http.StatusOK, categories)
}

func (server *Server) DeleteCategoryBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(errors.New("category slug is empty")))
		return
	}

	err := server.store.DeleteCategoryBySlug(ctx, slug)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrCannotDeleteEntity("category", err))
		return
	}
	message := fmt.Sprintf("category deleted. Slug: %s", slug)
	log.Println(message)
	ctx.JSON(http.StatusOK, gin.H{"message": message})
}

func (server *Server) SearchCategoriesByName(ctx *gin.Context) {
	text := ctx.Query("value")
	var response []*db.CategoryResponse

	categories, err := server.rdb.GetCategories(ctx, "all_categories")
	if err == nil {

		for _, c := range categories {
			if strings.Contains(strings.ToLower(c.Name), strings.ToLower(text)) {
				response = append(response, c)
			}
		}

		ctx.JSON(http.StatusOK, response)
		return
	}

	res, err := server.store.GetCategoriesName(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrCannotGetEntity(db.CategoryCollection, err))
		return
	}
	server.rdb.SetCategories(ctx, "all_categories", res)

	for _, c := range res {
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(text)) {
			response = append(response, c)
		}
	}

	ctx.JSON(http.StatusOK, response)

}
