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
	ID         string    `json:"id"`
	AuthorID   string    `json:"author_id"`
	Content    string    `json:"content"`
	PostID     string    `json:"post_id"`
	CommentID  string    `json:"comment_id"`
	EntityType string    `json:"entity_type"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type CommentResponse struct {
	ID         string         `json:"id"`
	Author     Author         `json:"author"`
	Content    string         `json:"content"`
	PostID     string         `json:"post_id"`
	CommentID  string         `json:"comment_id"`
	Comments   ChildComments  `json:"comments"`
	Likes      []LikeResponse `json:"likes"`
	EntityType string         `json:"entity_type"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
}

type ChildComments struct {
	Count int    `json:"count"`
	Link  string `json:"link"`
}

// CREATE TABLE IF NOT EXISTS "comments" (
//     comment_id character varying(45) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
//     comment_content text NOT NULL,
//     comment_post_id character varying(45) NO NULL,
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

func (c *Comment) Create() error {
	db := models.Connection

	query := `INSERT INTO comments (comment_content, comment_post_id, comment_comment_id, comment_user_id, comment_type) VALUES ($1, $2, $3, $4, $5) RETURNING comment_id, comment_created_at, comment_updated_at`

	var err error
	if c.EntityType == "" {
		return fmt.Errorf("entity type is required")
	}

	nullCommentID := sql.NullString{}

	if c.CommentID != "" {
		nullCommentID.String = c.CommentID
		nullCommentID.Valid = true
	}

	err = db.QueryRow(query, c.Content, c.PostID, nullCommentID, c.AuthorID, c.EntityType).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating comment: %s", err)
	}

	return nil
}

func (c *ChildComments) GetChildComments(entityType string, postID string, commentID string) error {
	db := models.Connection
	if entityType == "" {
		return fmt.Errorf("entity type is required")
	}

	var query string
	var err error
	if entityType == "post" {
		query = `SELECT COUNT(*) FROM comments WHERE comment_post_id = $1`
		err = db.QueryRow(query, postID).Scan(&c.Count)
	} else {
		query = `SELECT COUNT(*) FROM comments WHERE comment_comment_id = $1 AND comment_post_id = $2`
		err = db.QueryRow(query, commentID, postID).Scan(&c.Count)
	}

	if err != nil {
		return fmt.Errorf("error getting child comments: %s", err)
	}

	if c.Count > 0 {
		if entityType == "post" {
			c.Link = fmt.Sprintf("http://localhost:7070/comments?entity_type=%s&post_id=%s", entityType, postID)
		}
		if entityType == "comment" {
			c.Link = fmt.Sprintf("http://localhost:7070/comments?entity_type=%s&post_id=%s&comment_id=%s", entityType, postID, commentID)
		}
	}
	return nil
}

func (c *CommentResponse) GetComments(entityType string, postID string, commentID string) ([]CommentResponse, error) {
	db := models.Connection

	var query string
	var err error
	var rows *sql.Rows

	if entityType == "post" {
		query = `SELECT comment_id, comment_content, comment_user_id, user_name, user_picture, comment_post_id, comment_comment_id, comment_type, comment_created_at, comment_updated_at FROM comments 
		LEFT JOIN users ON comments.comment_user_id = users.user_id
		WHERE comment_post_id = $1 AND comment_type = 'post' ORDER BY comment_created_at DESC`
		rows, err = db.Query(query, postID)

	} else {
		query = `SELECT comment_id, comment_content, comment_user_id, user_name, user_picture, comment_post_id, comment_comment_id, comment_type, comment_created_at, comment_updated_at FROM comments 
		LEFT JOIN users ON comments.comment_user_id = users.user_id
		WHERE comment_comment_id = $1 AND comment_post_id = $2 AND comment_type = 'comment' ORDER BY comment_created_at DESC`

		rows, err = db.Query(query, commentID, postID)
	}

	if err != nil {
		return nil, fmt.Errorf("error getting comments: %s", err)
	}

	defer rows.Close()

	var comments []CommentResponse

	for rows.Next() {
		var comment CommentResponse

		var nullCommentID sql.NullString

		err := rows.Scan(&comment.ID, &comment.Content, &comment.Author.ID, &comment.Author.Name, &comment.Author.Picture, &comment.PostID, &nullCommentID, &comment.EntityType, &comment.CreatedAt, &comment.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("error scanning comment: %s", err)
		}

		if nullCommentID.Valid {
			comment.CommentID = nullCommentID.String
		}

		fmt.Printf("comment: %+v\n", comment)

		//get likes
		like := &LikeResponse{}

		likes, err := like.GetLikes(comment.ID, "comment")

		if err != nil {
			return nil, fmt.Errorf("error getting likes: %s", err)
		}

		comment.Likes = likes

		//get comments

		childComments := &ChildComments{}
		err = childComments.GetChildComments("comment", comment.PostID, comment.ID)

		if err != nil {
			return nil, fmt.Errorf("error getting comments: %s", err)
		}

		comment.Comments = *childComments

		comments = append(comments, comment)
	}

	return comments, nil
}

func (c *Comment) AddCommentLike() error {
	like := &Like{
		EntityID:   c.ID,
		EntityType: "comment",
		AuthorID:   c.AuthorID,
	}

	err := like.Create()

	if err != nil {
		return fmt.Errorf("error creating like: %s", err)
	}

	return nil
}
