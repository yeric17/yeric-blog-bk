package models

import (
	"fmt"
	"os"
	"strings"
	"time"
	models "yeric-blog/models/db"
	"yeric-blog/utils"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Picture  string `json:"picture,omitempty"`
	Status   string `json:"status"`
	RoleID   int    `json:"role_id"`
	Token    string `json:"token,omitempty"`
}

type UserResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture,omitempty"`
	Status  string `json:"status"`
	Role    string `json:"role"`
	Token   string `json:"token"`
}

type Roles struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type UserClaims struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
	Status  string `json:"status"`
	RoleID  int    `json:"role_id"`
	jwt.StandardClaims
}

func (u *User) GetUserByID() error {
	db := models.Connection

	query := `SELECT user_id, user_name, email, password, user_picture, user_status, role_name
	FROM users 
	LEFT JOIN roles ON role_id = user_role_id
	WHERE user_id = $1`

	err := db.QueryRow(query, u.ID).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Picture, &u.Status, &u.RoleID)

	if err != nil {
		return fmt.Errorf("error getting user: %s", err)
	}

	return nil
}

func (u *UserResponse) GetUsers() ([]UserResponse, error) {
	db := models.Connection

	query := `SELECT user_id, user_name, email, user_picture, user_status, role_name
	FROM users
	LEFT JOIN roles ON role_id = user_role_id`

	rows, err := db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("error getting users: %s", err)
	}

	defer rows.Close()

	var users []UserResponse

	for rows.Next() {
		var user UserResponse

		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Picture, &user.Status, &user.Role)

		if err != nil {
			return nil, fmt.Errorf("error getting users: %s", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (u *User) GetUserByEmail(email string) error {
	db := models.Connection

	query := `SELECT user_id, user_name, email, password, user_picture, user_status, user_role_id 
	FROM users
	WHERE email = $1`

	err := db.QueryRow(query, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Picture, &u.Status, &u.RoleID)

	if err != nil {
		return fmt.Errorf("error getting user: %s", err)
	}

	return nil
}

func (u *User) Create() error {
	db := models.Connection

	query := `INSERT INTO users (user_name, email, password, user_role_id) VALUES ($1, $2, $3, $4) RETURNING user_id`

	passByte, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)

	if err != nil {
		return fmt.Errorf("error creating user: %s", err)
	}

	u.Password = string(passByte)

	err = db.QueryRow(query, u.Name, u.Email, u.Password, u.RoleID).Scan(&u.ID)

	if err != nil {
		return fmt.Errorf("error creating user: %s", err)
	}

	return nil
}

func (u *User) Update() error {
	db := models.Connection

	var args []interface{}
	var instructions []string

	if u.Name != "" {
		instructions = append(instructions, fmt.Sprintf("user_name = $%d", len(args)+1))
		args = append(args, u.Name)
	}

	if u.Email != "" {
		instructions = append(instructions, fmt.Sprintf("email = $%d", len(args)+1))
		args = append(args, u.Email)
	}

	if u.Password != "" {
		instructions = append(instructions, fmt.Sprintf("password = $%d", len(args)+1))
		passByte, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)

		if err != nil {
			return fmt.Errorf("error updating user: %s", err)
		}
		u.Password = string(passByte)
		args = append(args, u.Password)
	}

	if u.Status != "" {
		instructions = append(instructions, fmt.Sprintf("user_status = $%d", len(args)+1))
		args = append(args, u.Status)
	}

	if u.RoleID != 0 {
		instructions = append(instructions, fmt.Sprintf("user_role_id = $%d", len(args)+1))
		args = append(args, u.RoleID)
	}

	if u.Picture != "" {
		instructions = append(instructions, fmt.Sprintf("user_picture = $%d", len(args)+1))
		args = append(args, u.Picture)
	}

	args = append(args, u.ID)

	query := fmt.Sprintf(`UPDATE "users" SET %s WHERE user_id = $%d`, strings.Join(instructions[:], ", "), len(args))

	_, err := db.Exec(query, args...)

	if err != nil {
		return fmt.Errorf("error updating user: %s", err)
	}

	return nil
}

func (u *User) Login() error {
	db := models.Connection

	query := `SELECT user_id, user_name, email, password, user_picture, user_status, user_role_id
	FROM users 
	WHERE email = $1`

	prevPass := u.Password

	err := db.QueryRow(query, u.Email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Picture, &u.Status, &u.RoleID)

	if err != nil {
		return utils.NewCustomError("credentials", fmt.Sprintf("error not found user: %s", err))
	}
	println(u.Password)
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(prevPass))

	if err != nil {
		return utils.NewCustomError("credentials", fmt.Sprintf("error comparing password: %s", err))
	}

	if u.Status == "email_not_verified" {
		return utils.NewCustomError("confirm_email", fmt.Sprintf("error user email not verified: %s", err))
	}
	if u.Status == "inactive" {
		return fmt.Errorf("error user is inactive")
	}

	u.Password = ""

	claim := UserClaims{
		u.ID,
		u.Name,
		u.Email,
		u.Picture,
		u.Status,
		u.RoleID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		return fmt.Errorf("error logging in: %s", err)
	}

	u.Token = tokenString

	return nil
}

func Authenticate(token string) (User, error) {
	var claim UserClaims

	_, err := jwt.ParseWithClaims(token, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return User{}, fmt.Errorf("error authenticating: %s", err)
	}

	updateUser := User{}

	err = updateUser.GetUserByEmail(claim.Email)

	if err != nil {
		fmt.Println(err)
		return User{}, fmt.Errorf("error authenticating: %s", err)
	}

	updateClaim := UserClaims{
		updateUser.ID,
		updateUser.Name,
		updateUser.Email,
		updateUser.Picture,
		updateUser.Status,
		updateUser.RoleID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, updateClaim)

	newTokenString, err := newToken.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		return User{}, fmt.Errorf("error authenticating: %s", err)
	}

	updateUser.Token = newTokenString
	updateUser.Password = ""
	return updateUser, nil
}

func SaveEmailToConfirm(email string) (string, error) {
	db := models.Connection

	query := `INSERT INTO email_confirms (email_confirms_email) VALUES ($1) RETURNING email_confirms_id`

	var emailConfirmID string

	err := db.QueryRow(query, email).Scan(&emailConfirmID)

	if err != nil {
		return "", fmt.Errorf("error saving email to confirm: %s", err)
	}

	return emailConfirmID, nil
}

func ConfirmEmail(id string) error {
	db := models.Connection

	var email string
	query := `SELECT email_confirms_email FROM email_confirms WHERE email_confirms_id = $1`

	err := db.QueryRow(query, id).Scan(&email)

	if err != nil {
		return fmt.Errorf("error confirming email: %s", err)
	}

	query = `UPDATE users SET user_status = $1 WHERE email = $2`

	//transaction
	tx, err := db.Begin()

	if err != nil {
		return fmt.Errorf("error confirming email transaction: %s", err)
	}

	_, err = tx.Exec(query, "active", email)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error confirming email set active: %s", err)
	}

	query = `DELETE FROM email_confirms WHERE email_confirms_id = $1`

	_, err = tx.Exec(query, id)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error confirming email delete confirm: %s", err)
	}

	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("error confirming email commit transaction: %s", err)
	}

	return nil
}

func (u *User) Delete() error {
	db := models.Connection

	query := `DELETE FROM users WHERE user_id = $1`

	_, err := db.Exec(query, u.ID)

	if err != nil {
		return fmt.Errorf("error deleting user: %s", err)
	}

	return nil
}
