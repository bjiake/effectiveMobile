package api

import (
	"effectiveMobile/pkg/api/handler"
	"github.com/gin-gonic/gin"

	_ "effectiveMobile/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.Handler) *ServerHTTP {
	engine := gin.New()

	// Use logger from Gin
	engine.Use(gin.Logger())

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	engine.POST("/registration", userHandler.Registration)
	engine.POST("/login", userHandler.Login)
	engine.GET("/info", userHandler.InfoPeople)
	engine.GET("/people", userHandler.GetPeople)
	engine.GET("/tasks", userHandler.GetAllTask)

	// Use middleware from Gin
	engine.Use(userHandler.AuthMiddleware())

	//Peoples
	engine.PUT("/people", userHandler.PutPeople)
	engine.DELETE("/people", userHandler.DeletePeople)

	//Task
	engine.GET("/people/task/", userHandler.GetTask)
	engine.POST("/people/task/start", userHandler.StartTask)
	engine.POST("/people/task/finish/:taskId", userHandler.FinishTask)
	engine.DELETE("/people/task/:taskId", userHandler.DeleteTask)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run("127.0.0.1:8001")
}
