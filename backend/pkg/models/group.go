package models

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type Group struct {
	ID          int            `json:"id,omitempty" db:"id"`
	CreatorID   int            `json:"creator_id,omitempty" db:"creator_id"`
	Title       string         `json:"title,omitempty" db:"title"`
	Description string         `json:"description,omitempty" db:"description"`
	CreatedAt   string         `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   string         `json:"updated_at,omitempty" db:"updated_at"`
	Members     []User         `json:"members,omitempty"`
	Posts       []PostResponse `json:"posts,omitempty"`
}
type GroupMember struct {
	ID        int    `db:"id"`
	GroupID   int    `json:"id" db:"group_id"`
	UserID    int    `json:"user_id" db:"user_id"`
	Status    string `db:"status"`
	InvitedBy int    `json:"invited_by" db:"invited_by"`
	JoinedAt  string `db:"joined_at"`
}

func AcceptGroupJoinRequest(db *sqlx.DB, userid, senderid, groupid int) {
	_, err := db.Exec(`UPDATE group_members SET status = "joined", updated_at = CURRENT_TIMESTAMP where group_id = ? and user_id = ?`,
		groupid, senderid)
	if err != nil {
		log.Printf("can't accept group join request: %v", err)
		return
	}
	err = NewNotification(db, senderid, userid, "accepted", groupid)
	if err != nil {
		log.Printf("can't add group_accept notif: %v", err)
	}
	log.Printf("group join request accepted successfully.")
}

func RejectGroupJoinRequest(db *sqlx.DB, userid, senderid, groupid int) {
	_, err := db.Exec(`UPDATE group_members SET status = "rejected", updated_at = CURRENT_TIMESTAMP where group_id = ? and user_id = ?`,
		groupid, senderid)
	if err != nil {
		log.Printf("can't reject group join request: %v", err)
		return
	}
	err = NewNotification(db, senderid, userid, "rejected", groupid)
	if err != nil {
		log.Printf("can't add group_accept notif: %v", err)
	}
	log.Printf("group join request accepted successfully.")
}

func GetGroup(db *sqlx.DB, groupid int) *Group {
	var group Group
	err := db.Get(&group, "SELECT id, title FROM groups WHERE id = ?", groupid)
	if err != nil {
		log.Printf("can't get group: %v", err)
		return nil
	}
	return &group
}
