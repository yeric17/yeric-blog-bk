package models

import (
	"bytes"
	"fmt"
	"html/template"
	"yeric-blog/config"
	models "yeric-blog/models/db"
)

type ContactMessage struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// CREATE TABLE IF NOT EXISTS "contacts" (
//     contacts_id character varying(45) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
//     contacts_name character varying(200) NOT NULL,
//     contacts_email character varying(200) NOT NULL,
//     contacts_message text NOT NULL
// )
// WITH (
//     OIDS = FALSE
// );

func (c *ContactMessage) Create() error {
	db := models.Connection

	query := `INSERT INTO contacts (contacts_name, contacts_email, contacts_message) VALUES ($1, $2, $3)`

	_, err := db.Exec(query, c.Name, c.Email, c.Message)

	if err != nil {
		return fmt.Errorf("error creating contact: %s", err.Error())
	}

	t, err := template.ParseFiles("email/contact_template.html")

	if err != nil {
		return fmt.Errorf("error parsing email template: %s", err.Error())
	}

	buff := new(bytes.Buffer)

	if err := t.Execute(buff, c); err != nil {
		return fmt.Errorf("error executing email template: %s", err.Error())
	}

	err = SendMail(fmt.Sprintf("%s contact you, email: %s", c.Name, c.Email), buff.String(), config.Mail.From)

	if err != nil {
		return fmt.Errorf("error sending email: %s", err.Error())
	}

	return nil
}
