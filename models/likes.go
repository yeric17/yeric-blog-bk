package models

import (
	"fmt"
	models "yeric-blog/models/db"
)

type Like struct {
	ID         string `json:"id"`
	AuthorID   string `json:"author_id"`
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"`
}

type LikeResponse struct {
	ID         string `json:"id"`
	AuthorID   string `json:"author_id"`
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"`
}

func (l *Like) Create() error {
	db := models.Connection

	//fmt.Printf("%+v\n", l)
	query := `INSERT INTO likes (like_entity_id, like_user_id, like_type) VALUES ($1, $2, $3)`

	_, err := db.Exec(query, l.EntityID, l.AuthorID, l.EntityType)

	return err
}

func (l *Like) Delete() error {
	db := models.Connection

	query := `DELETE FROM likes WHERE like_id = $1`

	_, err := db.Exec(query, l.ID)

	if err != nil {
		return fmt.Errorf("error deleting like: %v", err)
	}

	return nil
}

func (l *LikeResponse) GetLikes(entityID string, entityType string) ([]LikeResponse, error) {
	db := models.Connection

	query := `SELECT like_id, like_entity_id, like_user_id, like_type FROM likes WHERE like_entity_id = $1 AND like_type = $2`

	rows, err := db.Query(query, entityID, entityType)

	if err != nil {
		return nil, fmt.Errorf("error getting likes: %v", err)
	}

	defer rows.Close()

	var likes []LikeResponse

	for rows.Next() {
		var like LikeResponse

		err := rows.Scan(&like.ID, &like.EntityID, &like.AuthorID, &like.EntityType)

		if err != nil {
			return nil, fmt.Errorf("error getting likes: %v", err)
		}

		likes = append(likes, like)
	}

	return likes, nil
}
