package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"yeric-blog/models"
	"yeric-blog/utils"

	"github.com/gin-gonic/gin"
)

type ConfirmLink struct {
	Link string `json:"link"`
}

func CreateUser(g *gin.Context) {
	user := models.User{}

	if err := g.BindJSON(&user); err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error parsing user: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusBadRequest, resp)
		fmt.Println(err)
		return
	}

	if err := user.Create(); err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error creating user: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}

	user.Password = ""
	resp := utils.JSONResponse{
		Success: true,
		Message: "User created",
		Data:    user,
	}

	g.JSON(http.StatusCreated, resp)

}

func GetUsers(g *gin.Context) {
	user := &models.UserResponse{}

	users, err := user.GetUsers()

	if err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error getting users: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}

	resp := utils.JSONResponse{
		Success: true,
		Message: "Users retrieved",
		Data:    users,
	}

	g.JSON(http.StatusOK, resp)

}

func GetUserByEmail(g *gin.Context) {
	user := &models.User{}

	err := user.GetUserByEmail(g.Param("email"))

	if err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error getting user: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusNotFound, resp)
		fmt.Println(err)
		return
	}

	resp := utils.JSONResponse{
		Success: true,
		Message: "User retrieved",
		Data:    user,
	}

	g.JSON(http.StatusOK, resp)

}
func GetUserByID(g *gin.Context) {
	user := &models.User{}
	user.ID = g.Param("id")
	println(user.ID)
	err := user.GetUserByID()

	if err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error getting user: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusNotFound, resp)
		fmt.Println(err)
		return
	}

	resp := utils.JSONResponse{
		Success: true,
		Message: "User retrieved",
		Data:    user,
	}

	g.JSON(http.StatusOK, resp)

}

func UpdateUser(g *gin.Context) {
	user := &models.User{}

	if err := g.BindJSON(&user); err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error parsing user: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusBadRequest, resp)
		fmt.Println(err)
		return
	}

	if err := user.Update(); err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error updating user: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}

	user.Password = ""
	resp := utils.JSONResponse{
		Success: true,
		Message: "User updated",
		Data:    user,
	}

	g.JSON(http.StatusOK, resp)

}

func UserLogin(g *gin.Context) {
	user := &models.User{}

	if err := g.BindJSON(&user); err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error parsing user: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusBadRequest, resp)
		fmt.Println(err)
		return
	}
	err := user.Login()

	if err != nil {
		var message string
		switch e := err.(type) {
		case *utils.ConfirmEmailError:
			message = fmt.Sprintf("Error email not confirmed: %s", e.Error())
		case *utils.CredentialsError:
			message = fmt.Sprintf("Error invalid email or password: %s", e.Error())
		default:
			message = fmt.Sprintf("Error logging in: %s", e.Error())
		}

		resp := utils.JSONResponse{
			Success: false,
			Message: message,
			Data:    nil,
		}

		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}

	user.Password = ""

	resp := utils.JSONResponse{
		Success: true,
		Message: "User logged",
		Data:    user,
	}

	g.JSON(http.StatusOK, resp)

}

func Authenticate(g *gin.Context) {
	//Get the token from cookie
	token := g.Request.Header.Get("Authorization")

	if token == "" {
		resp := utils.JSONResponse{
			Success: false,
			Message: "Error getting token",
			Data:    nil,
		}
		g.JSON(http.StatusUnauthorized, resp)

		return
	}

	newUser, err := models.Authenticate(token)

	if err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error authenticating token: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusUnauthorized, resp)
		fmt.Println(err)
		return
	}

	resp := utils.JSONResponse{
		Success: true,
		Message: "Token authenticated",
		Data:    newUser,
	}

	g.JSON(http.StatusOK, resp)

}

func Register(g *gin.Context) {

	user := &models.User{}

	if err := g.BindJSON(&user); err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error parsing user: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusBadRequest, resp)
		fmt.Println(err)
		return
	}

	id, err := models.SaveEmailToConfirm(user.Email)

	if err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error confirming email save email: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}

	t, err := template.ParseFiles("email/confirm_email.html")

	if err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error parsing email template: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}

	buff := new(bytes.Buffer)

	err = t.Execute(buff, ConfirmLink{
		Link: fmt.Sprintf("http://localhost:7070/confirm/%s", id),
	})

	if err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error executing email template: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}

	err = models.SendMail("Yeric Blog, confirmar correo", buff.String(), user.Email)

	if err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error sending email: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}

	if err := user.Create(); err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error creating user: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}

	user.Password = ""
	resp := utils.JSONResponse{
		Success: true,
		Message: "User created",
		Data:    user,
	}

	g.JSON(http.StatusCreated, resp)
}

func ConfirmEmail(g *gin.Context) {
	id := g.Param("id")

	if err := models.ConfirmEmail(id); err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error confirming email: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}

	g.Redirect(http.StatusMovedPermanently, "http://localhost:3000/login")
}

func UploadUserPicture(g *gin.Context) {
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

	fileName := fmt.Sprintf("%s.%s", utils.RandomString(6), file.Filename[strings.LastIndex(file.Filename, ".")+1:])

	folderPath := fmt.Sprintf("%s/%s", "images/users", name)

	_, err = os.Stat(folderPath)

	//If exists folder remove it and create it again
	if !os.IsNotExist(err) {
		os.RemoveAll(folderPath)
	}
	os.Mkdir(folderPath, 0755)

	fullPath := fmt.Sprintf("%s/%s", folderPath, fileName)

	err = g.SaveUploadedFile(file, fullPath)
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

	user := &models.User{
		ID:      name,
		Picture: fmt.Sprintf("http://localhost:7070/%s", fullPath),
	}

	if err := user.Update(); err != nil {
		resp := utils.JSONResponse{
			Success: false,
			Message: fmt.Sprintf("Error updating user: %s", err.Error()),
			Data:    nil,
		}
		g.JSON(http.StatusInternalServerError, resp)
		fmt.Println(err)
		return
	}

	resp := utils.JSONResponse{
		Success: true,
		Message: "File uploaded",
		Data:    nil,
	}

	g.JSON(http.StatusOK, resp)
}
