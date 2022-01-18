package models

import (
	"database/sql"
	"fmt"
	"time"
	models "yeric-blog/models/db"
)

// CREATE TABLE IF NOT EXISTS "comments" (
//     comment_id character varying(45) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
//     comment_content text NOT NULL,
//     comment_post_id character varying(45),
//     comment_comment_id character varying(45),
//     comment_user_id character varying(45) NOT NULL,
//     comment_type character varying(45) NOT NULL,
//     comment_created_at timestamp with time zone NOT NULL DEFAULT now(),
//     comment_updated_at timestamp with time zone NOT NULL DEFAULT now(),
//     CONSTRAINT comment_type_check CHECK (comment_type IN ('post', 'comment'))
// )
// WITH (
//     OIDS = FALSE
// );
type Comment struct {
	ID              string    `json:"id"`
	AuthorID        string    `json:"author_id"`
	Content         string    `json:"content"`
	PostID          string    `json:"post_id"`
	ParentCommentID string    `json:"parent_id"`
	EntityType      string    `json:"entity_type"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}

type CommentResponse struct {
	ID              string        `json:"id"`
	Author          Author        `json:"author"`
	Content         string        `json:"content"`
	PostID          string        `json:"post_id"`
	Comments        ChildComments `json:"comments"`
	ParentCommentID string        `json:"parent_id"`
	// Likes      []LikeResponse `json:"likes"`
	EntityType string    `json:"entity_type"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type ChildComments struct {
	Count int    `json:"count"`
	Link  string `json:"link"`
}

// CREATE TABLE IF NOT EXISTS "comments" (
//     comment_id character varying(45) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
//     comment_content text NOT NULL,
//     comment_post_id character varying(45) NO NULL,
//     comment_user_id character varying(45) NOT NULL,
//     comment_type character varying(45) NOT NULL,
//     comment_created_at timestamp with time zone NOT NULL DEFAULT now(),
//     comment_updated_at timestamp with time zone NOT NULL DEFAULT now(),
//     CONSTRAINT comment_type_check CHECK (comment_type IN ('post', 'comment'))
// )
// WITH (
//     OIDS = FALSE
// );

// CREATE TABLE IF NOT EXISTS "parent_child_comments" (
//     parent_child_comments_id character varying(45) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
//     parent_child_comments_parent_id character varying(45) NOT NULL,
//     parent_child_comments_child_id character varying(45) NOT NULL,
//     CONSTRAINT parent_child_comments_parent_id_unique UNIQUE (parent_child_comments_parent_id, parent_child_comments_child_id)
// )
// WITH (
//     OIDS = FALSE
// );

func (c *Comment) Validate() error {
	if c.Content == "" {
		return fmt.Errorf("comment content is required")
	}

	if c.AuthorID == "" {
		return fmt.Errorf("comment author is required")
	}

	if c.PostID == "" {
		return fmt.Errorf("comment post is required")
	}

	if c.EntityType == "" {
		return fmt.Errorf("comment entity type is required")
	}

	if c.EntityType == "comment" && c.ParentCommentID == "" {
		return fmt.Errorf("comment parent is required")
	}

	return nil
}

func (c *Comment) Create() error {

	if err := c.Validate(); err != nil {
		return fmt.Errorf("error validating user: %s", err)
	}

	db := models.Connection

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error star transaction create comment: %s", err)
	}

	query := `INSERT INTO comments (comment_content, comment_post_id, comment_user_id, comment_type) VALUES ($1, $2, $3, $4) RETURNING comment_id, comment_created_at, comment_updated_at`

	err = tx.QueryRow(query, c.Content, c.PostID, c.AuthorID, c.EntityType).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating comment: %s", err)
	}

	if c.EntityType == "comment" {
		queryPath := `INSERT INTO parent_child_comments (parent_child_comments_parent_id, parent_child_comments_child_id) VALUES ($1, $2)`

		fmt.Printf("comment parent id: %s\n", c.ParentCommentID)
		fmt.Printf("comment id: %s\n", c.ID)

		_, err = tx.Exec(queryPath, c.ParentCommentID, c.ID)

		if err != nil {
			return fmt.Errorf("error creating path comment: %s", err)
		}
	}

	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("error creating comment: %s", err)
	}

	return nil
}

func (c *ChildComments) GetCommentChildComments(parentCommentID string) error {
	db := models.Connection

	var query string
	var err error

	query = `SELECT COUNT(*) FROM comments 
	LEFT join parent_child_comments ON parent_child_comments_parent_id = comments.comment_id
	WHERE parent_child_comments_parent_id = $1`

	err = db.QueryRow(query, parentCommentID).Scan(&c.Count)

	if err != nil {
		return fmt.Errorf("error getting child comments: %s", err)
	}

	if c.Count > 0 {
		c.Link = fmt.Sprintf("http://localhost:7070/comments?entity_type=comment&parent_id=%s", parentCommentID)
	}
	return nil
}

func (c *ChildComments) GetPostChildComments(postID string) error {
	db := models.Connection

	var query string
	var err error

	query = `SELECT COUNT(*) FROM comments WHERE comment_post_id = $1`
	err = db.QueryRow(query, postID).Scan(&c.Count)

	if err != nil {
		return fmt.Errorf("error getting child comments: %s", err)
	}

	if c.Count > 0 {
		c.Link = fmt.Sprintf("http://localhost:7070/comments?entity_type=post&parent_id=%s", postID)
	}
	return nil

}

func (c *CommentResponse) GetComments(entityType string, parentID string) ([]CommentResponse, error) {
	db := models.Connection

	var query string
	var err error
	var rows *sql.Rows

	if entityType == "comment" {
		query = `SELECT comment_id, comment_content, comment_post_id, comment_type, comment_created_at, comment_user_id, user_name, user_picture FROM comments
		LEFT JOIN parent_child_comments ON comment_id = parent_child_comments_child_id
		LEFT JOIN users ON user_id = comment_user_id
		WHERE parent_child_comments_parent_id = $1`

		rows, err = db.Query(query, parentID)
	}
	if entityType == "post" {
		query = `SELECT comment_id, comment_content, comment_post_id, comment_type, comment_created_at, comment_user_id, user_name, user_picture FROM comments
		LEFT JOIN users ON user_id = comment_user_id
		WHERE comment_post_id = $1 AND comment_type = 'post'`

		rows, err = db.Query(query, parentID)
	}

	if err != nil {
		return nil, fmt.Errorf("error getting comments: %s", err)
	}

	defer rows.Close()

	var comments []CommentResponse

	for rows.Next() {
		var comment CommentResponse

		err := rows.Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.EntityType, &comment.CreatedAt, &comment.Author.ID, &comment.Author.Name, &comment.Author.Picture)

		if err != nil {
			return nil, fmt.Errorf("error scanning comment: %s", err)
		}

		childComments := &ChildComments{}

		err = childComments.GetCommentChildComments(comment.ID)

		if err != nil {
			return nil, fmt.Errorf("error getting comments: %s", err)
		}

		comment.Comments = *childComments

		comments = append(comments, comment)
	}

	return comments, nil
}

func (c *CommentResponse) GetCommentByID(commentID string) error {
	db := models.Connection

	query := `SELECT comment_id, comment_content, comment_user_id, user_name, user_picture, comment_post_id, comment_type, comment_created_at, comment_updated_at FROM comments 
	LEFT JOIN users ON comments.comment_user_id = users.user_id
	WHERE comment_id = $1`

	err := db.QueryRow(query, commentID).Scan(&c.ID, &c.Content, &c.Author.ID, &c.Author.Name, &c.Author.Picture, &c.PostID, &c.EntityType, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error getting comment: %s", err)
	}

	return nil
}
