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

func GetContacts(g *gin.Context) {
	contact := &models.ContactResponse{}

	contacts, err := contact.GetContacts()

	if err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error getting contacts: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}

	resp := utils.JSONResponse{
		Success: true,
		Message: "Contacts",
		Data:    contacts,
	}

	g.JSON(http.StatusOK, resp)
}

func DeleteContact(g *gin.Context) {
	contact := models.ContactResponse{}
	contact.ID = g.Param("id")
	if(contact.ID == ""){
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error id is required: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusBadRequest, resp)
		fmt.Println(err)
		return		
	}
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

	if err := contact.Delete(); err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error deleting contact: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}
	resp := utils.JSONResponse{
		Success: true,
		Message: "Contact deleted",
		Data:    contact,
	}

	g.JSON(http.StatusOK, resp)
}