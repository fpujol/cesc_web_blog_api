package api

import (
	"context"
	"fmt"
	"time"

	db "blogapi/db/sqlc"
	"blogapi/pkg/token"
	"blogapi/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	log *logrus.Logger
	context context.Context
	config util.Config
	store db.Store
	router *gin.Engine
	tokenMaker token.Maker
}

func NewServer(log *logrus.Logger, ctx context.Context, config util.Config, store db.Store) (*Server, error) {
	
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		log: log,
		context: ctx,
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("currency", validCurrency)
	// }

	server.Init()
	return server, nil
}

func (s *Server) Init() {
	r := gin.Default()
	s.ConfigCors(r)
	s.ConfigMaxFileUpload(r)
	s.ConfigStaticAssets(r)
	s.ConfigRoutes(r)
	
}

func (s *Server) ConfigCors(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5010","https://webapi.indiketa.net","https://localhost:7050", "http://localhost:3000", "http://127.0.0.1:5500"},
		AllowMethods:     []string{"PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Origin","Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
}

func (s *Server) ConfigMaxFileUpload(r *gin.Engine) {
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
}

func (s *Server) ConfigStaticAssets(r *gin.Engine) {
	r.Use(static.Serve("/", static.LocalFile("./public", true)))
}

func (s *Server) ConfigRoutes(r *gin.Engine) {
	
	r.POST("/user/login", s.loginUser)
	r.POST("/user/logout", s.logout)
	
	v1 := r.Group("api/v1").Use(authMiddleware(s.tokenMaker))
	v1.GET("/user/me", s.Me)
	v1.GET("/post/:slug", s.GetBySlug)
	v1.GET("/post/", s.GetBySlug)
	v1.GET("/post", s.ListPosts)
	v1.POST("/post", s.RegisterNewPost)
	v1.PUT("/post", s.EditPost)
	v1.POST("/upload/single", s.UploadFile)
	v1.PUT("/upload/main-image/:id", s.UploadMainImageFile)
	v1.PUT("/upload/thumbnail-image/:id", s.UploadThumbnailImageFile)

	s.router = r
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
