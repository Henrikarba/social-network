package models

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type Notification struct {
	ID       int    `json:"id" db:"id"`
	UserID   int    `json:"user_id" db:"user_id"`
	SenderID int    `json:"sender_id" db:"sender_id"`
	GroupID  int    `json:"group_id" db:"group_id"`
	Type     string `json:"type" db:"type"`
	Message  string `json:"message" db:"message"`
	IsRead   bool   `json:"is_read" db:"is_read"`

	Sender *User  `json:"sender"`
	Group  *Group `json:"group,omitempty"`
}

func NewNotification(db *sqlx.DB, userid, senderid int, status string, groupid int) error {
	tx, err := db.Beginx()
	if err != nil {
		return nil
	}
	var msg string
	var notif_type string
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if groupid == 0 {

		if status == "accepted" {
			msg = "accepted your follow request"
			notif_type = "follow_accept"

		} else if status == "rejected" {
			msg = "rejected your follow request"
			notif_type = "follow_accept"
		} else {
			msg = "wants to follow you"
			notif_type = "follow_request"
		}

		query := `
    INSERT INTO user_notifications (user_id, sender_id, type, message)
    VALUES (?, ?, ?, ?)`

		_, err = tx.Exec(query, userid, senderid, notif_type, msg)
		if err != nil {
			return err
		}
	} else if groupid >= 1 {
		if status == "accepted" {
			msg = "accepted your request to join group"
			notif_type = "group_join_accept"
		} else if status == "rejected" {
			msg = "rejected your request to join group"
			notif_type = "group_join_accept"

		} else if status == "group_join_request" {
			msg = "wants to join your group"
			notif_type = "group_join_request"
		} else if status == "group_join_invite" {
			msg = "invited you to join group"
			notif_type = "group_join_invite"
		} else if status == "group_join_invite_accept" {
			msg = "accepted your invitation to join"
			notif_type = "group_join_invite_accept"
		} else if status == "group_join_invite_reject" {
			msg = "rejected your invitation to join"
			notif_type = "group_join_invite_reject"
		}

		query := `INSERT INTO user_notifications (user_id, sender_id, type, message, group_id)
		VALUES (?, ?, ?, ?, ?)`

		_, err := tx.Exec(query, userid, senderid, notif_type, msg, groupid)
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

func GetNotifications(db *sqlx.DB, userid int) ([]Notification, error) {
	var notifs []Notification

	err := db.Select(&notifs, "SELECT id, user_id, sender_id, group_id, type, message FROM user_notifications WHERE is_read = false AND user_id = ?", userid)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return notifs, nil
}

func MarkNotificationAsRead(db *sqlx.DB, id int) {
	_, err := db.Exec("UPDATE user_notifications SET is_read = true WHERE id = ?", id)
	if err != nil {
		log.Printf("can't mark as read: %v", err)
	}
}
