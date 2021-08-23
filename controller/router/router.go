package router

import (
	"thinhtd4/controller/handler"
	"thinhtd4/controller/middleware"
	"thinhtd4/customlog"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func Init() {
	customlog.Info("Initialize router")

	router = gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.LoadHTMLGlob("./view/*")

	authorized := router.Group("/")
	authorized.GET("/login", handler.HandlleLoginGET)
	authorized.GET("/register", handler.HandleRegisterGET)
	authorized.POST("/login", handler.HandleLoginPOST)
	authorized.POST("/register", handler.HandleRegisterPOST)

	authorized.Use(middleware.AuthMiddleWare())
	{
		authorized.GET("", func(c *gin.Context) {})
		authorized.GET("/survey", handler.HandleSurvey)
		authorized.GET("/question", handler.HandleGetListQuestion)
		authorized.PUT("/submitone", handler.HandleUpdateUserQuestionDone)
		authorized.PUT("/submitall", handler.HandleSubmit)
	}
}

func Router() *gin.Engine {
	return router
}
