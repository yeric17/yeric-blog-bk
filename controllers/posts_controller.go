package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"yeric-blog/config"
	"yeric-blog/models"
	"yeric-blog/utils"

	"github.com/gin-gonic/gin"
)

func CreatePost(g *gin.Context) {
	post := models.Post{}

	if err := g.BindJSON(&post); err != nil {
		reps := &utils.JSONResponse{
			Success: false,
			Message: "Invalid post data; " + err.Error(),
			Data:    nil,
		}
		g.JSON(http.StatusBadRequest, reps)
		return
	}

	if err := post.Create(); err != nil {
		reps := &utils.JSONResponse{
			Success: false,
			Message: "Error creating post; " + err.Error(),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, reps)
		fmt.Printf("%+v\n", err)
		return
	}

	reps := &utils.JSONResponse{
		Success: true,
		Message: "Post created successfully",
		Data:    post,
	}
	g.JSON(http.StatusCreated, reps)
}

func GetPosts(g *gin.Context) {
	post := models.Post{}
	posts, err := post.GetPosts()
	if err != nil {
		reps := &utils.JSONResponse{
			Success: false,
			Message: "Error getting posts; " + err.Error(),
			Data:    nil,
		}
		g.JSON(http.StatusNotFound, reps)
		return
	}

	reps := &utils.JSONResponse{
		Success: true,
		Message: "Posts retrieved successfully",
		Data:    posts,
	}
	g.JSON(http.StatusOK, reps)
}

func CreateComment(g *gin.Context) {
	comment := models.Comment{}

	if err := g.BindJSON(&comment); err != nil {
		reps := &utils.JSONResponse{
			Success: false,
			Message: "Invalid comment data; " + err.Error(),
			Data:    nil,
		}
		g.JSON(http.StatusBadRequest, reps)
		return
	}

	//fmt.Printf("%+v\n", comment)

	if err := comment.Create(); err != nil {
		reps := &utils.JSONResponse{
			Success: false,
			Message: "Error creating comment; " + err.Error(),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, reps)
		fmt.Printf("%+v\n", err)
		return
	}

	reps := &utils.JSONResponse{
		Success: true,
		Message: "Comment created successfully",
		Data:    comment,
	}

	g.JSON(http.StatusCreated, reps)
}

func GetPostByID(g *gin.Context) {
	post := models.PostResponse{}

	if err := post.GetPostByID(g.Param("id")); err != nil {
		reps := &utils.JSONResponse{
			Success: false,
			Message: "Error getting post; " + err.Error(),
			Data:    nil,
		}
		g.JSON(http.StatusNotFound, reps)
		return
	}

	reps := &utils.JSONResponse{
		Success: true,
		Message: "Post retrieved successfully",
		Data:    post,
	}
	g.JSON(http.StatusOK, reps)
}

func GetComments(g *gin.Context) {
	parent_id := g.Query("parent_id")
	entity_type := g.Query("entity_type")

	//fmt.Printf("post_id: %s; comment_id: %s; entity_type: %s\n", post_id, comment_id, entity_type)

	if entity_type == "" {
		resp := &utils.JSONResponse{
			Success: false,
			Message: "Invalid entity type",
			Data:    nil,
		}
		g.JSON(http.StatusBadRequest, resp)
		return
	}

	if parent_id == "" {
		resp := &utils.JSONResponse{
			Success: false,
			Message: "Invalid post id",
			Data:    nil,
		}
		g.JSON(http.StatusBadRequest, resp)
		return
	}

	comment := models.CommentResponse{}

	comments, err := comment.GetComments(entity_type, parent_id)

	if err != nil {
		reps := &utils.JSONResponse{
			Success: false,
			Message: "Error getting comments; " + err.Error(),
			Data:    nil,
		}
		g.JSON(http.StatusNotFound, reps)
		return
	}

	resp := &utils.JSONResponse{
		Success: true,
		Message: "Comments retrieved successfully",
		Data:    comments,
	}
	g.JSON(http.StatusOK, resp)
}

func GetCommentByID(g *gin.Context) {
	comment := models.CommentResponse{}

	if err := comment.GetCommentByID(g.Param("id")); err != nil {
		reps := &utils.JSONResponse{
			Success: false,
			Message: "Error getting comment; " + err.Error(),
			Data:    nil,
		}
		g.JSON(http.StatusNotFound, reps)
		return
	}

	resp := &utils.JSONResponse{
		Success: true,
		Message: "Comment retrieved successfully",
		Data:    comment,
	}
	g.JSON(http.StatusOK, resp)
}

func UploadPostImage(g *gin.Context) {
	file, err := g.FormFile("file")
	name := g.Query("name")

	if name == "" {
		resp := utils.JSONResponse{
			Success: false,
			Message: "Error uploading file: name is required",
			Data:    nil,
		}
		g.JSON(http.StatusBadRequest, resp)
		fmt.Println(err)
		return
	}

	if err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error getting file: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}

	fileName := fmt.Sprintf("%s.%s", name, file.Filename[strings.LastIndex(file.Filename, ".")+1:])
	err = g.SaveUploadedFile(file, "./images/posts/"+fileName)

	if err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error saving file: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}

	resp := utils.JSONResponse{
		Success: true,
		Message: "File uploaded",
		Data:    "http://localhost:" + config.APP_PORT + "/images/posts/" + fileName,
	}

	g.JSON(http.StatusOK, resp)
}

func GetPostsCategories(g *gin.Context) {

	categories, err := models.GetPostsCategories()

	if err != nil {
		reps := &utils.JSONResponse{
			Success: false,
			Message: "Error getting posts categories; " + err.Error(),
			Data:    nil,
		}
		g.JSON(http.StatusNotFound, reps)
		return
	}

	resp := &utils.JSONResponse{
		Success: true,
		Message: "Posts categories retrieved successfully",
		Data:    categories,
	}
	g.JSON(http.StatusOK, resp)
}
