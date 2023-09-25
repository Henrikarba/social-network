package models

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type Group struct {
	ID          int    `json:"id,omitempty" db:"id"`
	CreatorID   int    `json:"creator_id,omitempty" db:"creator_id"`
	Title       string `json:"title,omitempty" db:"title"`
	Description string `json:"description,omitempty" db:"description"`
	CreatedAt   string `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   string `json:"updated_at,omitempty" db:"updated_at"`

	Status     string         `json:"status"`
	ChatroomID int            `json:"chatroom_id"`
	Members    []User         `json:"members,omitempty"`
	MemberIDS  []int          `json:"member_ids,omitempty"`
	Posts      []PostResponse `json:"posts,omitempty"`
}

type GroupMember struct {
	ID        int    `db:"id"`
	GroupID   int    `json:"id" db:"group_id"`
	UserID    int    `json:"user_id" db:"user_id"`
	Status    string `db:"status"`
	InvitedBy int    `json:"invited_by" db:"invited_by"`
	JoinedAt  string `db:"joined_at"`
}

func GetAllGroups(db *sqlx.DB) ([]Group, error) {
	query := "SELECT id, title, creator_id, created_at FROM groups"
	var groups []Group

	err := db.Select(&groups, query)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func InviteToGroup(db *sqlx.DB, groupID int, userID int, invitedby int) {
	_, err := db.Exec(`
		INSERT INTO group_members (group_id, user_id, status, invited_by)
		VALUES (?, ?, "invited", ?)`,
		groupID, userID, invitedby)
	if err != nil {
		log.Printf("can't add group_join_invite: %v", err)
		return
	}

	// Create a notification for the group join request
	err = NewNotification(db, userID, invitedby, "group_join_invite", groupID)
	if err != nil {
		log.Printf("can't add group_join_request notification: %v", err)
	}
}

func AcceptGroupJoinInvite(db *sqlx.DB, userid, senderid, groupid int) {

	_, err := db.Exec(`UPDATE group_members SET status = "joined", updated_at = CURRENT_TIMESTAMP where group_id = ? and user_id = ?`,
		groupid, userid)
	if err != nil {
		log.Printf("can't accept group join request: %v", err)
		return
	}
	err = NewNotification(db, senderid, userid, "group_join_invite_accept", groupid)
	if err != nil {
		log.Printf("can't add group_accept notif: %v", err)
		return
	}
	log.Printf("group join invite accepted successfully.")
}

func DeclineGroupJoinInvite(db *sqlx.DB, userid, senderid, groupid int) {

	_, err := db.Exec(`UPDATE group_members SET status = "rejected", updated_at = CURRENT_TIMESTAMP where group_id = ? and user_id = ?`,
		groupid, userid)
	if err != nil {
		log.Printf("can't reject group join request: %v", err)
		return
	}
	err = NewNotification(db, senderid, userid, "group_join_invite_reject", groupid)
	if err != nil {
		log.Printf("can't add group_reject notif: %v", err)
		return
	}
	log.Printf("group join request accepted successfully.")
}

func JoinGroup(db *sqlx.DB, userID, creatorID, groupID int) {
	_, err := db.Exec(`
		INSERT INTO group_members (group_id, user_id, status)
		VALUES (?, ?, "requested")`,
		groupID, userID)
	if err != nil {
		log.Printf("can't add group_join_request: %v", err)
		return
	}

	// Create a notification for the group join request
	err = NewNotification(db, creatorID, userID, "group_join_request", groupID)
	if err != nil {
		log.Printf("can't add group_join_request notification: %v", err)
	}
	log.Printf("group join request added successfully.")
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
		return
	}
	log.Printf("group join invited accepted successfully.")
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
	log.Printf("group join invite rejected successfully.")
}

func GetGroup(db *sqlx.DB, groupID int, userID int) *Group {
	var group Group
	err := db.Get(&group, "SELECT * FROM groups WHERE id = ?", groupID)
	if err != nil {
		log.Printf("selecting group: %v", err)
		return nil
	}
	group.Members = []User{}
	var members []User
	err = db.Select(&members, `
    SELECT u.id, u.first_name, u.last_name 
    FROM users u
    JOIN group_members gm ON u.id = gm.user_id
    WHERE gm.group_id = ? AND gm.status = 'joined'`, groupID)
	if err != nil {
		log.Printf("selecting members: %v", err)
		return nil
	}
	group.Members = members

	// Check if the requesting user is a member of the group
	isMember := false
	for _, member := range members {
		if member.ID == userID {
			isMember = true
			break
		}
	}

	// If the user is a member, retrieve posts
	if isMember {
		group.Posts = []PostResponse{}
		var posts []PostResponse
		err = db.Select(&posts, `
            SELECT id, title, created_at 
            FROM group_posts 
            WHERE group_id = ?
            ORDER BY created_at DESC`, groupID)
		if err != nil {
			log.Printf("selecting posts: %v", err)
			return nil
		}
		group.Posts = posts
	} else {
		group.Posts = nil // User is not a member, so no posts are returned
	}

	return &group
}

func GetGroupMembers(db *sqlx.DB, groupid int) []int {
	var members []int

	err := db.Select(&members, `
	SELECT user_id FROM group_members WHERE group_id = ? AND status = 'joined'
	`, groupid)
	if err != nil {
		log.Printf("selecting group members: %v", err)
		return nil
	}

	return members
}

func GetAnyGroupMembers(db *sqlx.DB, groupid int) []int {
	var members []int

	err := db.Select(&members, `
	SELECT user_id FROM group_members WHERE group_id = ? AND status = 'joined' OR status = 'requested' OR status = 'invited'
	`, groupid)
	if err != nil {
		log.Printf("selecting group members: %v", err)
		return nil
	}

	return members
}
