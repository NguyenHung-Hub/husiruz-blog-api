package api

import (
	"fmt"
	"husir_blog/src/db"
	"husir_blog/src/rdb"
	"husir_blog/src/token"
	"husir_blog/src/util"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

type Server struct {
	cfg    util.Config
	store  db.Store
	rdb    *rdb.RedisClient
	token  token.Maker
	router *gin.Engine
	cld    *cloudinary.Cloudinary
}

func NewServer(config util.Config, store db.Store, rdb *rdb.RedisClient) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %s", err)
	}

	cld, err := cloudinary.NewFromParams(config.CloudinaryName, config.CloudinaryApiKey, config.CloudinarySecret)
	if err != nil {
		return nil, fmt.Errorf("cannot create cloudinary: %s", err)
	}

	server := &Server{cfg: config, store: store, rdb: rdb, token: tokenMaker, cld: cld}
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.MaxMultipartMemory = 8 << 20
	router.Static("/assets", "./assets")

	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://localhost:8080"},
	// 	AllowCredentials: true,
	// 	// Enable Debugging for testing, consider disabling in production
	// 	Debug: true,
	// })
	router.Use(cors.AllowAll())

	router.POST("/user", server.createUser)
	router.GET("/user/:id", server.getUserById)
	router.POST("/login", server.login)
	router.POST("/tokens/refresh", server.renewAccessToken)
	router.POST("/tokens/check", server.checkRefreshToken)
	router.GET("/post/:slug", server.getPostBySlug)
	router.GET("/post", server.ListPost)
	router.GET("/post/random", server.ListPostRandom)
	router.GET("/post/cate", server.getPostByCategory)

	router.GET("/categories", server.getCategoriesName)
	router.DELETE("/categories/:slug", server.DeleteCategoryBySlug)
	router.GET("/categories/search", server.SearchCategoriesByName)

	router.POST("/upload", server.uploadFile)
	router.POST("/upload2", server.uploadImageToCloud)
	router.GET("/session/check/:id", server.checkSession)

	authRoutes := router.Group("/").Use(authMiddleware(server.token))
	// authRoutes.POST("/upload", server.uploadFile)
	authRoutes.POST("/category", server.createCategory)
	authRoutes.GET("/category/:id", server.getCategoryById)
	authRoutes.POST("/post", server.createPost)

	server.router = router
}

func (server *Server) Start(port string) error {
	return server.router.Run(port)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}

}
