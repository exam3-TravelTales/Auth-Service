package api

import (
	_ "auth/api/docs"
	"auth/api/handler"
	"auth/api/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @title User
// @version 1.0
// @description API Gateway of Authorazation
// @host localhost:8085
// BasePath: /

func Router(hand *handler.Handler) *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", hand.Register)
		auth.POST("/login", hand.Login)
		auth.POST("/refresh", hand.Refresh)
		auth.POST("/logout", hand.Logout)
	}

	userAuth := router.Group("/api/v1/auth")
	userAuth.Use(middleware.Check)
	{
		userAuth.POST("/reset-password", hand.ResetPassword)
	}

	user := router.Group("/api/v1/users")
	user.Use(middleware.Check)
	{
		user.GET("/profile", hand.Profile)
		user.PUT("/profile", hand.UserProfileUpdate)
		user.GET("", hand.GetAllUsers)
		user.DELETE("/:user_id", hand.Delete)
		user.GET("/:user_id/activity", hand.Activity)
		user.POST("/:user_id/follow", hand.Follow)
		user.GET("/:user_id/followers", hand.GetFollowers)
	}

	return router
}
