package controllers

import (
	"fmt"
	"net/http"
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

func AddLike(g *gin.Context) {
	like := models.Like{}

	if err := g.BindJSON(&like); err != nil {
		reps := &utils.JSONResponse{
			Success: false,
			Message: "Invalid like data; " + err.Error(),
			Data:    nil,
		}
		g.JSON(http.StatusBadRequest, reps)
		return
	}

	if err := like.Create(); err != nil {
		reps := &utils.JSONResponse{
			Success: false,
			Message: "Error creating like; " + err.Error(),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, reps)
		return
	}

	reps := &utils.JSONResponse{
		Success: true,
		Message: "Like created successfully",
		Data:    like,
	}

	g.JSON(http.StatusCreated, reps)
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

	fmt.Printf("%+v\n", comment)

	if err := comment.Create(); err != nil {
		reps := &utils.JSONResponse{
			Success: false,
			Message: "Error creating comment; " + err.Error(),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, reps)
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
	post_id := g.Query("post_id")
	comment_id := g.Query("comment_id")
	entity_type := g.Query("entity_type")

	fmt.Printf("post_id: %s; comment_id: %s; entity_type: %s\n", post_id, comment_id, entity_type)

	if entity_type == "" {
		resp := &utils.JSONResponse{
			Success: false,
			Message: "Invalid entity type",
			Data:    nil,
		}
		g.JSON(http.StatusBadRequest, resp)
		return
	}

	comment := models.CommentResponse{}

	comments, err := comment.GetComments(entity_type, post_id, comment_id)

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
