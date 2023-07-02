package main

import (
	"github.com/microservices-simulator-api/internal/authentication"
	"github.com/microservices-simulator-api/internal/boards"
	"github.com/microservices-simulator-api/internal/config"
	"github.com/microservices-simulator-api/internal/lives"
	"github.com/microservices-simulator-api/internal/shapes"
	"github.com/microservices-simulator-api/internal/users"
	"github.com/microservices-simulator-api/internal/utils/hashutil"
	"github.com/microservices-simulator-api/internal/utils/jwtutil"
	"github.com/microservices-simulator-api/internal/utils/redis"
	"log"
)

type Container struct {
	jwt                      *jwtutil.Authentication
	authenticationProvider   authentication.Provider
	authenticationManager    authentication.Manager
	authenticationController *authentication.Controller
	boardController          *boards.Controller
	liveBoardService         *lives.LiveBoardService
	shapeController          *shapes.Controller
	userController           *users.Controller
}

func NewContainer(cfg config.Config) *Container {
	jwtAuthentication, err := jwtutil.NewAuthentication(cfg.Security)
	if err != nil {
		log.Fatal(err)
	}

	pwdHashing := hashutil.NewPasswordHashing(cfg.Security.Salts)
	redisClient := redis.NewClient()
	// users
	userRepository := users.NewRepository(redisClient)
	userService := users.NewService(pwdHashing, userRepository)
	userController := users.NewController(userService)
	// boards
	boardRepository := boards.NewRepository(redisClient)
	boardService := boards.NewService(boardRepository)
	boardController := boards.NewController(boardService)
	// lives
	liveBoardManager := lives.NewLiveBoardManager()
	liveBoardService := lives.NewLiveBoardService(liveBoardManager, boardService, userService)
	// Shapes
	shapeRepository := shapes.NewRepository(redisClient)
	shapeService := shapes.NewService(shapeRepository)
	shapeController := shapes.NewController(shapeService)
	// security
	authenticationProvider := authentication.NewProvider(pwdHashing, jwtAuthentication, userService)
	authenticationManager := authentication.NewManager(authenticationProvider, redisClient)
	authenticationController := authentication.NewController(cfg.Security, authenticationManager)

	return &Container{
		jwt:                      jwtAuthentication,
		authenticationProvider:   authenticationProvider,
		authenticationManager:    authenticationManager,
		authenticationController: authenticationController,
		boardController:          boardController,
		liveBoardService:         liveBoardService,
		shapeController:          shapeController,
		userController:           userController,
	}
}
