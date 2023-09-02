package models

import (
	"social-network/pkg/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type Comment struct {
	ID        int     `db:"id"`
	PostID    int     `json:"post_id" db:"post_id"`
	UserID    int     `json:"user_id" db:"user_id"`
	Content   string  `json:"content" db:"content"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	UpdatedAt string  `json:"updated_at" db:"updated_at"`
	CreatedBy *User   `json:"created_by" db:"created_by"`
	ImageUrl  *string `json:"image_url,omitempty" db:"image_url"`
}

func GetCommentForRegularPost(db *sqlx.DB, postid int) ([]Comment, error) {
	var comments []Comment
	query := `
		SELECT
			id,
			post_id,
			user_id,
			content,
			created_at,
			updated_at,
			image_url
		FROM
			comments
		WHERE
			post_id = ?
		ORDER BY
			created_at;
	`

	err := db.Select(&comments, query, postid)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func GetCommentForGroupPost(db *sqlx.DB, postid int) ([]Comment, error) {
	var comments []Comment
	query := `
		SELECT
			id,
			post_id,
			user_id,
			content,
			created_at,
			updated_at,
			image_url
		FROM
			group_comments
		WHERE
			post_id = ?
		ORDER BY
			created_at;
	`
	err := db.Select(&comments, query, postid)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func NewComment(p PostRequest, db *sqlx.DB) (*Comment, error) {
	now := time.Now()
	comment := &Comment{
		PostID:    p.ID,
		UserID:    p.CreatedBy,
		Content:   p.Content,
		CreatedAt: utils.FormatTime(&now),
	}

	var stmt string
	switch p.PostTarget {
	case "regular_post":
		if p.ImageData != nil {
			imagePath, err := utils.SaveImage(p.ImageData, p.MimeType, "comment")
			if err != nil {
				return nil, err
			}
			comment.ImageUrl = &imagePath
		}

		stmt = `INSERT INTO comments (post_id, user_id, content, image_url)
			VALUES (:post_id, :user_id, :content, :image_url)`

	case "group_post":
		if p.ImageData != nil {
			imagePath, err := utils.SaveImage(p.ImageData, p.MimeType, "comment")
			if err != nil {
				return nil, err
			}
			comment.ImageUrl = &imagePath
		}

		stmt = `INSERT INTO group_comments (post_id, user_id, content, image_url)
			VALUES (:post_id, :user_id, :content, :image_url)`
	}

	result, err := db.NamedExec(stmt, comment)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	comment.ID = int(id)
	return comment, nil
}
