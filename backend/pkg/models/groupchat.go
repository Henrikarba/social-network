package models

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

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
