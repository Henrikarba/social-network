package models

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func GetLastUnreadGroupMessages(db *sqlx.DB, recipientID int) ([]Message, error) {
	var groupIDs []int

	groupQuery := `
        SELECT gm.group_id
        FROM group_members gm
        WHERE gm.user_id = ? AND gm.status = 'joined'
    `

	if err := db.Select(&groupIDs, groupQuery, recipientID); err != nil {
		return nil, err
	}

	if len(groupIDs) == 0 {
		return nil, nil
	}

	var messages []Message

	for _, groupID := range groupIDs {
		chatroomID, err := GetGroupChatRoomID(db, groupID)
		if err != nil {
			return nil, err
		}

		messageQuery := `
			SELECT messages.*
			FROM messages
			WHERE recipient_id = 0 AND chatroom_id = ? AND id = (
				SELECT MAX(id) 
				FROM messages 
				WHERE recipient_id = 0 AND chatroom_id = ?
			)
		`

		var groupMessages []Message
		if err := db.Select(&groupMessages, messageQuery, chatroomID, chatroomID); err != nil {
			return nil, err
		}

		messages = append(messages, groupMessages...)
	}

	var user *User

	for i := range messages {
		user, _ = GetPrivateProfile(db, messages[i].SenderID)
		messages[i].Sender = user
	}

	return messages, nil
}

func InsertGroupMessage(db *sqlx.DB, fromid int, groupid int, content string) error {
	chatroomID, err := GetGroupChatRoomID(db, groupid)
	if err != nil {
		return err
	}

	if chatroomID == 0 {
		id, err := CreateChatRoomForGroups(db, groupid)
		if err != nil {
			return err
		}
		chatroomID = id
	}

	tx, err := db.Beginx()
	if err != nil {
		// return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	message := Message{
		SenderID:    fromid,
		RecipientID: 0,
		Content:     content,
		ChatroomID:  chatroomID,
	}

	_, err = tx.NamedExec(`
        INSERT INTO messages (sender_id, recipient_id, content, chatroom_id)
        VALUES (:sender_id, :recipient_id, :content, :chatroom_id)
    `, message)

	// Check if user is in chatroom_participants table
	var count int
	err = tx.Get(&count, `
        SELECT COUNT(*) FROM chatroom_participants 
        WHERE chatroom_id = ? AND user_id = ?
    `, chatroomID, fromid)

	if err != nil {
		return err
	}

	if count == 0 {
		_, err = tx.Exec(`
            INSERT INTO chatroom_participants (chatroom_id, user_id)
            VALUES (?, ?)
        `, chatroomID, fromid)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetGroupMessageHistory(db *sqlx.DB, groupid int) ([]Message, error) {
	chatroomID, err := GetGroupChatRoomID(db, groupid)
	if err != nil {
		return nil, err
	}

	var messages []Message
	query := `
        SELECT * FROM messages
        WHERE chatroom_id = ?
        ORDER BY created_at ASC
    `

	err = db.Select(&messages, query, chatroomID)
	if err != nil {
		return nil, err
	}

	for i := range messages {
		profile, _ := GetPrivateProfile(db, messages[i].SenderID)
		messages[i].CreatedBy = profile
	}
	return messages, nil
}

func GetGroupChatRoomID(db *sqlx.DB, groupid int) (int, error) {
	var chatroomID int
	query := `
        SELECT id FROM chatrooms
        WHERE name = ?
        LIMIT 1
    `

	err := db.Get(&chatroomID, query, fmt.Sprintf("g-%d", groupid))
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return chatroomID, nil
}

func CreateChatRoomForGroups(db *sqlx.DB, groupid int) (int, error) {
	tx, err := db.Beginx()
	if err != nil {
		return 0, err
	}
	defer func() {
		err = tx.Commit()
		if err != nil {
			_ = tx.Rollback()
			return
		}
	}()

	cr := Chatroom{
		Name: fmt.Sprintf("g-%d", groupid),
	}
	result, err := tx.NamedExec(`INSERT INTO chatrooms (name) VALUES (:name)`, cr)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}
