package models

import (
	"fmt"
	models "yeric-blog/models/db"
)

type ContactMessage struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

type ContactResponse struct {
	ID      string `json:"id"`
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

	return nil
}

func (c *ContactResponse) GetContacts() ([]ContactResponse, error) {
	db := models.Connection

	query := `SELECT contacts_id, contacts_name, contacts_email, contacts_message FROM contacts ORDER BY contacts_create_at DESC`

	rows, err := db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("error getting contacts: %s", err.Error())
	}

	defer rows.Close()

	var contacts []ContactResponse

	for rows.Next() {
		var contact ContactResponse

		err := rows.Scan(&contact.ID, &contact.Name, &contact.Email, &contact.Message)

		if err != nil {
			return nil, fmt.Errorf("error getting contacts: %s", err.Error())
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func (c *ContactResponse) Delete() error {
	db := models.Connection

	query := `DELETE FROM contacts WHERE contacts_id = $1`

	_, err := db.Exec(query, c.ID)

	if err != nil {
		return fmt.Errorf("error deleting contact: %s", err.Error())
	}

	return nil
}
