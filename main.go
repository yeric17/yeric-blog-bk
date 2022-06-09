package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"yeric-blog/controllers"

	corsgin "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func repeatHandler(r int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var buffer bytes.Buffer
		for i := 0; i < r; i++ {
			buffer.WriteString("Hello from Go!\n")
		}
		c.String(http.StatusOK, buffer.String())
	}
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	tStr := os.Getenv("REPEAT")
	repeat, err := strconv.Atoi(tStr)

	if err != nil {
		log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
		repeat = 5
	}

	router := gin.New()
	router.Use(gin.Logger())
	//router.Use(TestMiddleWare())
	gin.SetMode(gin.ReleaseMode)

	//cors policy gin config
	router.SetTrustedProxies([]string{"http://localhost:3000", "https://yeric-blog-web.herokuapp.com"})
	router.Use(corsgin.New(corsgin.Config{
		AllowOrigins:     []string{"https://yericdev.herokuapp.com", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	//router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Yeric Blog API")
	})

	router.Static("/images", "./images")
	//image size 1MB
	router.MaxMultipartMemory = 1 << 20

	router.POST("/users", controllers.CreateUser)
	router.GET("/users", Authenticate(), controllers.GetUsers)
	router.GET("/users/id/:id", Authenticate(), controllers.GetUserByID)
	router.GET("/users/email/:email", Authenticate(), controllers.GetUserByEmail)
	router.PUT("/users", Authenticate(), controllers.UpdateUser)
	router.POST("/users/login", controllers.UserLogin)
	router.GET("/users/auth", controllers.Authenticate)
	router.POST("/users/register", controllers.Register)
	router.GET("/confirm/:id", controllers.ConfirmEmail)
	router.POST("/contact", controllers.ContactEmail)
	router.GET("/contact", Authenticate(), controllers.GetContacts)
	router.DELETE("/contact/id/:id", Authenticate(), controllers.DeleteContact)
	router.POST("/users/upload", Authenticate(), controllers.UploadUserPicture)

	router.POST("/posts", Authenticate(), controllers.CreatePost)
	router.GET("/posts", controllers.GetPosts)
	router.POST("/posts/comment", Authenticate(), controllers.CreateComment)
	router.GET("/posts/id/:id", controllers.GetPostByID)
	router.GET("/comments", controllers.GetComments)
	router.GET("/comments/id/:id", controllers.GetCommentByID)
	router.POST("/posts/upload", Authenticate(), controllers.UploadPostImage)
	router.GET("/posts/categories", controllers.GetPostsCategories)

	router.GET("/repeat", repeatHandler(repeat))

	router.Run(":" + port)

}

func Authenticate() gin.HandlerFunc {
	return controllers.Authenticate
}
