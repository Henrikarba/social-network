package models

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type Follower struct {
	FollowerID int    `db:"follower_id"`
	FolloweeID int    `db:"followee_id"`
	Status     string `db:"status"`
}

func GetFollowers(db *sqlx.DB, id int) ([]User, error) {
	var followers []User
	err := db.Select(&followers, `
		SELECT u.id, u.first_name, u.last_name, u.avatar, f.status
		FROM users u
		INNER JOIN followers f ON u.id = f.follower_id
		WHERE f.followee_id = ?
	`, id)
	if err != nil {
		return nil, err
	}

	return followers, nil
}

func GetFollowing(db *sqlx.DB, id int) ([]User, error) {
	var following []User
	err := db.Select(&following, `
		SELECT u.id, u.first_name, u.last_name, u.avatar, f.status
		FROM users u
		INNER JOIN followers f ON u.id = f.followee_id
		WHERE f.follower_id = ?
	`, id)
	if err != nil {
		return nil, err
	}

	return following, nil
}
func CreateFollowRequest(db *sqlx.DB, followerID, followeeID int) error {
	tx, err := db.Beginx()
	if err != nil {
		return nil
	}

	var privacy int
	query := `SELECT privacy FROM users WHERE id = ?`
	err = db.Get(&privacy, query, followeeID)
	if err != nil {
		fmt.Println("?", followeeID, followerID)
		return err
	}

	var status string
	if privacy == 0 {
		status = "accepted"
	} else if privacy == 1 {
		status = "pending"
	}

	follower := &Follower{
		Status:     status,
		FolloweeID: followeeID,
		FollowerID: followerID,
	}

	_, err = tx.NamedExec(`INSERT INTO followers (follower_id, followee_id, status) VALUES (:follower_id, :followee_id, :status)`, follower)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}

	err = NewNotification(db, followeeID, followerID, status, 0)
	if err != nil {
		log.Printf("newnotif: %v", err)
		return nil
	}
	fmt.Println("asd")
	return nil
}

func AcceptFollowRequest(db *sqlx.DB, userid, senderid int) {
	_, err := db.Exec("UPDATE followers SET status = 'accepted', updated_at = CURRENT_TIMESTAMP WHERE followee_id = ? AND follower_id = ?", userid, senderid)
	if err != nil {
		log.Printf("can't accept follow request: %v", err)
		return
	}
	err = NewNotification(db, senderid, userid, "accepted", 0)
	log.Printf("Follow request accepted successfully.")
}

func RejectFollowRequest(db *sqlx.DB, userid, senderid int) {
	_, err := db.Exec("UPDATE followers SET status = 'rejected', updated_at = CURRENT_TIMESTAMP WHERE followee_id = ? AND follower_id = ?", userid, senderid)
	if err != nil {
		log.Printf("can't accept follow request: %v", err)
		return
	}
	err = NewNotification(db, senderid, userid, "rejected", 0)
	log.Printf("Follow request accepted successfully.")
}
