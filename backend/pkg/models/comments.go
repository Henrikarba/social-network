package models

import (
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
