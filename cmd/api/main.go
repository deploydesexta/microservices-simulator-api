package main

import (
	"github.com/labstack/echo/v4"
	"github.com/microservices-simulator-api/internal/config"
	"github.com/microservices-simulator-api/internal/health"
	"github.com/microservices-simulator-api/internal/utils/logger"
	"github.com/rs/zerolog/log"
)

type Server struct {
	cfg config.Config
	di  *Container
}

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	logger.Setup(logger.Level(cfg.LogLevel))

	log.Info().Msg("Starting server...")
	NewServer(cfg).Start()
}

func NewServer(cfg config.Config) *Server {
	return &Server{
		cfg: cfg,
		di:  NewContainer(cfg),
	}
}

func (s *Server) Start() {
	middlewares := NewAuthenticationMiddleware(s.di)
	blacklistMid := middlewares.BlacklistMiddleware()
	jwtMid := middlewares.JwtMiddleware()

	e := echo.New()
	e.HideBanner = true
	e.Use(blacklistMid)
	e.GET("/ping", health.Ping)

	ag := e.Group("/auth")
	ag.POST("/login", s.di.authenticationController.Login)
	ag.GET("/me", s.di.authenticationController.Me, jwtMid)

	bg := e.Group("/boards")
	bg.Use(jwtMid)
	bg.POST("", s.di.boardController.CreateBoard)
	bg.GET("", s.di.boardController.BoardsOfUser)
	bg.GET("/:id", s.di.boardController.BoardOfId)
	bg.GET("/:id/ws", s.di.liveBoardService.NewBoardConnection)

	sh := e.Group("/shapes")
	sh.Use(jwtMid)
	sh.POST("", s.di.shapeController.CreateShape)
	sh.GET("", s.di.shapeController.AllShapes)
	sh.GET("/:id", s.di.shapeController.ShapeOfId)

	ug := e.Group("/users")
	ug.POST("/register", s.di.userController.Register)
	ug.GET("/me", s.di.authenticationController.Me, jwtMid)

	e.Logger.Fatal(e.Start(":8081"))
}
