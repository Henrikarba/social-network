package models

import (
	"log"
	"social-network/pkg/utils"

	"github.com/jmoiron/sqlx"
)

type PostRequest struct {
	ID         int    `json:"id"`
	CreatedBy  int    `schema:":-"`
	PostTarget string `json:"post_target" schema:"post_target"`
	Title      string `schema:"title"`
	Content    string `schema:"content"`
	Privacy    int    `schema:"privacy"`
	ImageData  []byte `schema:"image"`
	GroupID    int    `schema:"groupid"`
	MimeType   string `schema:"mimeType"`
	Followers  []int  `schema:"followers"`
}

type PostResponse struct {
	PostID  int `json:"post_id,omitempty" db:"post_id"`
	UserID  int `json:"user_id,omitempty" db:"user_id"`
	Privacy int `json:"privacy,omitempty" db:"privacy"`

	Title    string    `json:"title,omitempty" db:"title"`
	Content  string    `json:"content,omitempty" db:"content"`
	GroupID  int       `json:"group_id,omitempty" db:"group_id"`
	ImageUrl *string   `json:"image_url,omitempty" db:"image_url"`
	Comments []Comment `json:"comments,omitempty"`

	CreatedAt string `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty" db:"updated_at"`

	PostTarget string `json:"post_target" schema:"post_target"`
	Group      *Group `json:"group,omitempty"`
	CreatedBy  *User  `json:"created_by,omitempty"`
}

func GetPost(userid int, postID int, db *sqlx.DB) (*PostResponse, error) {
	var post PostResponse

	query := `
    SELECT
        p.id as post_id,
        p.user_id,
        p.title,
        p.content,
        p.image_url,
        p.privacy,
        p.created_at
    FROM
        posts p
    WHERE
        p.id = ? AND (
            p.privacy = 0
            OR (p.privacy = 1 AND EXISTS (SELECT 1 FROM followers f WHERE f.followee_id = p.user_id AND f.follower_id = ? AND f.status = 'accepted'))
            OR (p.privacy = 2 AND EXISTS (SELECT 1 FROM followers f WHERE f.followee_id = p.user_id AND f.follower_id = ? AND f.status = 'accepted') AND EXISTS (SELECT 1 FROM post_followers pf WHERE pf.post_id = p.id AND pf.user_id = ?))
        )
`

	err := db.Get(&post, query, postID, userid, userid, userid)
	if err != nil {
		return nil, err
	}
	user, _ := GetPrivateProfile(db, post.UserID)
	post.CreatedBy = user

	comments, _ := GetCommentForRegularPost(db, post.PostID)
	for i := range comments {
		user, _ := GetPrivateProfile(db, comments[i].UserID)
		comments[i].CreatedBy = user
	}

	post.Comments = comments

	return &post, nil
}

func GetSingleGroupPost(db *sqlx.DB, userid, postID int) (*PostResponse, error) {
	var post PostResponse
	query := `
        SELECT
            gp.id AS post_id,
            gp.user_id AS user_id,
            gp.group_id,
            gp.title,
            gp.content,
            gp.created_at,
            gp.image_url
        FROM
            group_posts gp
        JOIN
            group_members gm ON gp.group_id = gm.group_id
        WHERE
            gm.user_id = ? AND gm.status = 'joined' AND gp.id = ?
    `

	err := db.Get(&post, query, userid, postID)
	if err != nil {
		return nil, err
	}

	user, err := GetPrivateProfile(db, post.UserID)
	if err != nil {
		log.Printf("Error retrieving user profile: %v", err)
	}

	group := GetGroup(db, post.GroupID)
	if err != nil {
		log.Printf("Error retrieving groups: %v", err)
	}
	comments, _ := GetCommentForGroupPost(db, post.PostID)
	for i := range comments {
		user, _ := GetPrivateProfile(db, comments[i].UserID)
		comments[i].CreatedBy = user
	}

	post.Comments = comments
	post.Group = group
	post.CreatedBy = user

	return &post, nil
}

func NewPost(p PostRequest, db *sqlx.DB) (*PostResponse, error) {
	post := &PostResponse{}

	switch p.PostTarget {
	case "regular_post":
		newPost, err := handleNewRegularPost(p, db)
		if err != nil {
			return nil, err
		}

		if len(p.Followers) > 0 && p.Privacy == 2 {
			err = addPostFollowers(newPost.PostID, p.Followers, db)
			if err != nil {
				log.Printf("can't add post followers: %v", err)
			}
		}

		post = newPost
	case "group_post":
		newPost, err := handleNewGroupPost(p, db)
		if err != nil {
			return nil, err
		}
		post = newPost
	}

	return post, nil
}

func GetRegularPosts(db *sqlx.DB, userid int) ([]PostResponse, error) {
	var posts []PostResponse
	query := `
        SELECT
            p.id AS post_id,
            p.user_id AS user_id,
            p.title,
            p.content,
			p.created_at,
			p.image_url,
            p.privacy
        FROM
            posts p
        WHERE
            p.user_id = ?
            OR (p.privacy = 0)
            OR (p.privacy = 1 AND EXISTS (SELECT 1 FROM followers f WHERE f.followee_id = p.user_id AND f.follower_id = ?))
            OR (p.privacy = 2 AND EXISTS (SELECT 1 FROM followers f WHERE f.followee_id = p.user_id AND f.follower_id = ?) AND EXISTS (SELECT 1 FROM post_followers pf WHERE pf.post_id = p.id AND pf.user_id = ?))
        ORDER BY
            p.created_at DESC;
    `

	err := db.Select(&posts, query, userid, userid, userid, userid)
	if err != nil {
		return nil, err
	}

	// Retrieve and attach user information to each post's CreatedBy field
	for i := range posts {
		user, err := GetPrivateProfile(db, posts[i].UserID)
		if err != nil {
			log.Printf("what %v", err)

		}
		comments, _ := GetCommentForRegularPost(db, posts[i].PostID)
		posts[i].Comments = comments
		posts[i].CreatedBy = user
	}
	return posts, nil
}

func GetGroupPosts(db *sqlx.DB, userid int) ([]PostResponse, error) {
	var posts []PostResponse
	query := `
        SELECT
            gp.id AS post_id,
            gp.user_id AS user_id,
            gp.group_id,
            gp.title,
            gp.content,
            gp.created_at,
            gp.image_url
        FROM
            group_posts gp
        JOIN
            group_members gm ON gp.group_id = gm.group_id
        WHERE
            gm.user_id = ? AND gm.status = 'joined'
        ORDER BY
            gp.created_at DESC;
    `

	err := db.Select(&posts, query, userid)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		user, err := GetPrivateProfile(db, posts[i].UserID)
		if err != nil {
			log.Printf("Error retrieving user profile: %v", err)
		}
		group := GetGroup(db, posts[i].GroupID)
		if err != nil {
			log.Printf("Error retrieving groups: %v", err)
		}
		comments, _ := GetCommentForGroupPost(db, posts[i].PostID)
		posts[i].Comments = comments
		posts[i].Group = group
		posts[i].CreatedBy = user
	}
	return posts, nil
}

func handleNewRegularPost(p PostRequest, db *sqlx.DB) (*PostResponse, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	newPost := &PostResponse{
		UserID:  p.CreatedBy,
		Title:   p.Title,
		Content: p.Content,
		Privacy: p.Privacy,
	}

	if p.ImageData != nil {
		path, err := utils.SaveImage(p.ImageData, p.MimeType, "post")
		if err != nil {
			log.Println(err)
		}
		newPost.ImageUrl = &path
	}

	result, err := tx.NamedExec(`
		INSERT INTO posts (user_id, title, content, image_url, privacy, created_at, updated_at)
		VALUES (:user_id, :title, :content, :image_url, :privacy, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, newPost)

	if err != nil {
		return nil, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	newPost.PostID = int(postID)
	return newPost, nil
}

func handleNewGroupPost(p PostRequest, db *sqlx.DB) (*PostResponse, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	newPost := &PostResponse{
		UserID:  p.CreatedBy,
		Title:   p.Title,
		Content: p.Content,
		Privacy: p.Privacy,
		GroupID: p.GroupID,
	}

	if p.ImageData != nil {
		path, err := utils.SaveImage(p.ImageData, p.MimeType, "post")
		if err != nil {
			log.Println(err)
		}
		newPost.ImageUrl = &path
	}

	result, err := tx.NamedExec(`
		INSERT INTO group_posts (user_id, group_id, title, content, image_url, created_at, updated_at)
		VALUES (:user_id, :group_id, :title, :content, :image_url, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, newPost)

	if err != nil {
		return nil, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	newPost.PostID = int(postID)
	return newPost, nil
}

func addPostFollowers(postid int, followers []int, db *sqlx.DB) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	for _, followerID := range followers {
		_, err := tx.Exec(`
                INSERT INTO post_followers (post_id, user_id, created_at)
                VALUES (?, ?, CURRENT_TIMESTAMP)
            `, postid, followerID)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
