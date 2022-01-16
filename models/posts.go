package models

import (
	"fmt"
	"strings"
	"time"
	models "yeric-blog/models/db"
)

// CREATE TABLE IF NOT EXISTS "posts" (
//     post_id character varying(45) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
//     post_title character varying(200) NOT NULL,
//     post_content text NOT NULL,
//     post_image text NOT NULL DEFAULT 'http://localhost:7070/images/default_image.png',
//     post_author_id character varying(45) NOT NULL,
//     post_created_at timestamp with time zone NOT NULL DEFAULT now(),
//     post_updated_at timestamp with time zone NOT NULL DEFAULT now()
// )
// WITH (
//     OIDS = FALSE
// );

type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type PostResponse struct {
	ID        string         `json:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Author    Author         `json:"author"`
	Comments  ChildComments  `json:"comments"`
	Likes     []LikeResponse `json:"likes"`
	Image     string         `json:"image"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Author struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (p *Post) Create() error {
	db := models.Connection

	query := `INSERT INTO posts (post_title, post_content, post_author_id) VALUES ($1, $2, $3) RETURNING post_id, post_created_at, post_updated_at`

	err := db.QueryRow(query, p.Title, p.Content, p.AuthorID).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating post: %v", err)
	}

	//fmt.Println(p.CreatedAt.Hour())
	return nil
}

func (p *Post) GetPosts() ([]PostResponse, error) {
	db := models.Connection

	query := `SELECT post_id, post_title, post_content, post_image, post_created_at, post_updated_at, post_author_id, user_name, user_picture
	FROM posts
	INNER JOIN users ON posts.post_author_id = users.user_id`

	rows, err := db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("error getting posts: %v", err)
	}

	defer rows.Close()

	var posts []PostResponse

	for rows.Next() {
		var post PostResponse

		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.CreatedAt, &post.UpdatedAt, &post.Author.ID, &post.Author.Name, &post.Author.Picture)

		if err != nil {
			return nil, fmt.Errorf("error scanning posts: %v", err)
		}

		comment := ChildComments{}

		err = comment.GetChildComments("post", post.ID, "")

		if err != nil {
			return nil, fmt.Errorf("error getting comments: %v", err)
		}

		post.Comments = comment

		like := &LikeResponse{}
		post.Likes, err = like.GetLikes(post.ID, "post")

		if err != nil {
			return nil, fmt.Errorf("error getting like count: %v", err)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (p *Post) AddPostLike(like Like) error {

	err := like.Create()

	if err != nil {
		return fmt.Errorf("error creating like: %v", err)
	}

	return nil
}

func (p *PostResponse) GetPostByID(id string) error {
	db := models.Connection

	query := `SELECT post_id, post_title, post_content, post_image, post_created_at, post_updated_at, post_author_id, user_name, user_picture
	FROM posts
	INNER JOIN users ON posts.post_author_id = users.user_id
	WHERE post_id = $1`

	err := db.QueryRow(query, id).Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.CreatedAt, &p.UpdatedAt, &p.Author.ID, &p.Author.Name, &p.Author.Picture)

	if err != nil {
		return fmt.Errorf("error getting post: %v", err)
	}

	comment := &ChildComments{}

	err = comment.GetChildComments("post", p.ID, "")

	if err != nil {
		return fmt.Errorf("error getting comments: %v", err)
	}

	p.Comments = *comment

	like := &LikeResponse{}
	p.Likes, err = like.GetLikes(p.ID, "post")

	if err != nil {
		return fmt.Errorf("error getting like count: %v", err)
	}

	return nil
}

func (p *Post) Update() error {
	db := models.Connection

	var args []interface{}
	var instructions []string

	if p.Title != "" {
		instructions = append(instructions, fmt.Sprintf("post_title = $%d", len(args)+1))
		args = append(args, p.Title)
	}

	if p.Content != "" {
		instructions = append(instructions, fmt.Sprintf("post_content = $%d", len(args)+1))
		args = append(args, p.Content)
	}

	if p.Image != "" {
		instructions = append(instructions, fmt.Sprintf("post_image = $%d", len(args)+1))
		args = append(args, p.Image)
	}

	if len(instructions) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE posts SET %s WHERE post_id = $%d", strings.Join(instructions, ", "), len(args)+1)

	args = append(args, p.ID)

	_, err := db.Exec(query, args...)

	if err != nil {
		return fmt.Errorf("error updating post: %v", err)
	}

	return nil
}
