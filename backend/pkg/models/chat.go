package models

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type Chatroom struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type ChatroomParticipant struct {
	ChatroomID int    `db:"chatroom_id"`
	UserID     int    `db:"user_id"`
	JoinedAt   string `db:"joined_at"`
}

type Message struct {
	ID          int    `json:"id,omitempty" db:"id"`
	SenderID    int    `json:"sender_id,omitempty" db:"sender_id"`
	Sender      *User  `json:"sender,omitempty"`
	RecipientID int    `json:"recipient_id,omitempty" db:"recipient_id"`
	Content     string `json:"content,omitempty" db:"content"`
	CreatedAt   string `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   string `json:"updated_at,omitempty" db:"updated_at"`
	IsRead      bool   `json:"is_read,omitempty" db:"is_read"`
	CreatedBy   *User  `json:"created_by,omitempty"`
	ChatroomID  int    `json:"chatroom_id,omitempty" db:"chatroom_id"`
}

type ChatListEntry struct {
	ID          int    `json:"id"`
	Participant User   `json:"sender"`
	LastMessage string `json:"message"`
}

type FullChat struct {
	Messages []Message `json:"messages"`
	Partner  *User     `json:"partner"`
}

func GetChatList(db *sqlx.DB, userid int) *[]ChatListEntry {
	var senderIDs []int

	query := `
    SELECT m.sender_id
    FROM messages m
    WHERE m.recipient_id = ?
    GROUP BY m.sender_id
`

	if err := db.Select(&senderIDs, query, userid); err != nil {
		log.Printf("getting recipient user ids: %v", err)
	}
	var chatlist []ChatListEntry
	for i := range senderIDs {
		profile, _ := GetPrivateProfile(db, senderIDs[i])
		entry := ChatListEntry{
			Participant: *profile,
		}
		chatlist = append(chatlist, entry)
	}

	return &chatlist
}

func GetFullChat(db *sqlx.DB, userid, partnerid int) *FullChat {
	partner, _ := GetPrivateProfile(db, partnerid)

	var messages []Message
	query := `
    SELECT m.*
    FROM messages m
    WHERE (m.sender_id = ? AND m.recipient_id = ?) OR (m.sender_id = ? AND m.recipient_id = ?)
`
	if err := db.Select(&messages, query, userid, partnerid, partnerid, userid); err != nil {
		log.Printf("messages: %v", err)
	}

	updateQuery := `
	    UPDATE messages
	    SET is_read = true
	    WHERE recipient_id = ? AND sender_id = ?
	`
	_, err := db.Exec(updateQuery, userid, partnerid)
	if err != nil {
		log.Printf("update is_read: %v", err)
	}

	chat := FullChat{
		Partner:  partner,
		Messages: messages,
	}

	return &chat
}

func GetLastUnreadMessages(db *sqlx.DB, recipientID int) ([]Message, error) {
	var messages []Message

	query := `
		SELECT messages.*
		FROM messages
		WHERE id IN (
			SELECT MAX(id) AS max_id
			FROM messages
			WHERE recipient_id = ? AND is_read = 0
			GROUP BY sender_id
		)
	`
	if err := db.Select(&messages, query, recipientID, recipientID); err != nil {
		return nil, err
	}

	var user *User

	for i := range messages {
		user, _ = GetPrivateProfile(db, messages[i].SenderID)
		messages[i].Sender = user
	}
	return messages, nil
}

func InsertMessage(db *sqlx.DB, senderID, recipientID int, content string) (int64, error) {
	chatroomID, err := GetChatRoomID(db, senderID, recipientID)
	if err != nil {
		return 0, err
	}

	if chatroomID == 0 {
		err = CreateChatRoomForUsers(db, senderID, recipientID)
		if err != nil {
			return 0, err
		}
		chatroomID, err = GetChatRoomID(db, senderID, recipientID)
		if err != nil {
			return 0, err
		}
	}

	tx, err := db.Beginx()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	message := Message{
		SenderID:    senderID,
		RecipientID: recipientID,
		Content:     content,
		ChatroomID:  chatroomID,
	}

	res, err := tx.NamedExec(`
        INSERT INTO messages (sender_id, recipient_id, content, chatroom_id)
        VALUES (:sender_id, :recipient_id, :content, :chatroom_id)
    `, message)
	id, _ := res.LastInsertId()

	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetMessageHistory(db *sqlx.DB, senderID, recipientID int) ([]Message, error) {
	chatroomID, err := GetChatRoomID(db, senderID, recipientID)
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

	return messages, nil
}

func CreateChatRoomForUsers(db *sqlx.DB, id1, id2 int) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	name := GetChatRoomNameForUserChat(id1, id2)

	cr := Chatroom{Name: name}
	result, err := tx.NamedExec(`INSERT INTO chatrooms (name) VALUES (:name)`, cr)
	if err != nil {
		return err
	}

	chatroomID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	cp1 := ChatroomParticipant{ChatroomID: int(chatroomID), UserID: id1}
	_, err = tx.NamedExec(`INSERT INTO chatroom_participants (chatroom_id, user_id) VALUES (:chatroom_id, :user_id)`, cp1)
	if err != nil {
		return err
	}

	cp2 := ChatroomParticipant{ChatroomID: int(chatroomID), UserID: id2}
	_, err = tx.NamedExec(`INSERT INTO chatroom_participants (chatroom_id, user_id) VALUES (:chatroom_id, :user_id)`, cp2)
	if err != nil {
		return err
	}

	return nil
}

func GetChatRoomNameForUserChat(id1, id2 int) string {
	var name string
	if id1 < id2 {
		name = strconv.Itoa(id1) + "-" + strconv.Itoa(id2)
	} else {
		name = strconv.Itoa(id2) + "-" + strconv.Itoa(id1)
	}

	return name
}

func GetChatRoomID(db *sqlx.DB, id1, id2 int) (int, error) {
	var chatroomID int

	name := GetChatRoomNameForUserChat(id1, id2)
	query := `
        SELECT id FROM chatrooms
        WHERE name = ?
        LIMIT 1
    `

	err := db.Get(&chatroomID, query, name)
	if err != nil {

		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return chatroomID, nil
}
