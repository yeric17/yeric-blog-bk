package main

import (
	"time"
	"yeric-blog/config"
	"yeric-blog/controllers"

	corsgin "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	//cors policy gin config
	router.SetTrustedProxies([]string{"http://localhost:3000", "https://yeric-blog-web.herokuapp.com"})
	router.Use(corsgin.New(corsgin.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	//router.Use(gin.Logger())

	router.Static("/images", "./images")
	//image size 1MB
	router.MaxMultipartMemory = 1 << 20
	router.POST("/users", controllers.CreateUser)
	router.GET("/users", controllers.GetUsers)
	router.GET("/users/id/:id", controllers.GetUserByID)
	router.GET("/users/email/:email", controllers.GetUserByEmail)
	router.PUT("/users", controllers.UpdateUser)
	router.POST("/users/login", controllers.UserLogin)
	router.GET("/users/auth", controllers.Authenticate)
	router.POST("/users/register", controllers.Register)
	router.GET("/confirm/:id", controllers.ConfirmEmail)
	router.POST("/contact", controllers.ContactEmail)
	router.POST("/users/upload", controllers.UploadUserPicture)

	router.POST("/posts", controllers.CreatePost)
	router.GET("/posts", controllers.GetPosts)
	router.POST("/posts/comment", controllers.CreateComment)
	router.GET("/posts/id/:id", controllers.GetPostByID)
	router.GET("/comments", controllers.GetComments)
	router.GET("/comments/id/:id", controllers.GetCommentByID)
	router.POST("/posts/upload", controllers.UploadPostImage)
	router.GET("/posts/categories", controllers.GetPostsCategories)

	router.Run(":" + config.APP_PORT)

}
