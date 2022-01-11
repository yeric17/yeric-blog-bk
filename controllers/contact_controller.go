package controllers

import (
	"fmt"
	"net/http"
	"yeric-blog/models"
	"yeric-blog/utils"

	"github.com/gin-gonic/gin"
)

func ContactEmail(g *gin.Context) {
	contact := models.ContactMessage{}

	if err := g.BindJSON(&contact); err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error parsing contact: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusBadRequest, resp)
		fmt.Println(err)
		return
	}

	if err := contact.Create(); err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error creating contact: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}
	resp := utils.JSONResponse{
		Success: true,
		Message: "Contact send email",
		Data:    contact,
	}

	g.JSON(http.StatusCreated, resp)
}
