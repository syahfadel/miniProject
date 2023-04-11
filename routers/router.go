package routers

import (
	controller "miniProject/controllers"
	"miniProject/middlewares"
	"miniProject/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StartService(db *gorm.DB) *gin.Engine {
	sosmedService := services.SosmedService{
		DB: db,
	}

	sosmedController := controller.SosmedController{
		DB:            db,
		SosmedService: &sosmedService,
	}

	app := gin.Default()
	app.POST("/register", sosmedController.Register)
	app.POST("/login", sosmedController.Login)

	app.Use(middlewares.Authentication())
	app.POST("/follow/:id", sosmedController.Follow)
	app.POST("/post", sosmedController.CreatePost)
	app.GET("/users", sosmedController.GetAllUser)
	app.GET("/following", sosmedController.GetFollowing)
	app.GET("/followed", sosmedController.GetFollowed)
	app.GET("/posts", sosmedController.GetPostByFollowing)
	app.GET("/posts/:id", sosmedController.GetPostByUserId)
	app.DELETE("/unfollow/:id", sosmedController.Unfollow)
	app.DELETE("/post/:id", sosmedController.DeletePost)

	return app

}
