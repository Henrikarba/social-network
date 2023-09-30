package models

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type Event struct {
	ID         int        `json:"id"`
	GroupID    int        `json:"group_id" db:"group_id"`
	CreatedBy  int        `json:"created_by" db:"created_by"`
	User       *User      `json:"user,omitempty"`
	Title      string     `json:"title" db:"title"`
	Content    string     `json:"content" db:"content"`
	EventStart *time.Time `json:"event_start" db:"event_start"`
	EventEnd   *time.Time `json:"event_end" db:"event_end"`
	CreatedAt  string     `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  string     `json:"updated_at,omitempty" db:"updated_at"`

	EventResponses []EventResponse `json:"responses,omitempty"`
}

type EventResponse struct {
	ID       int    `json:"id" db:"id"`
	EventID  int    `json:"event_id,omitempty" db:"event_id"`
	Response string `json:"response,omitempty" db:"response"`
	UserID   int    `json:"user_id" db:"user_id"`
	User     *User  `json:"user"`
}

func InsertEvent(db *sqlx.DB, e Event) (*Event, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	fmt.Println(e.EventStart)
	res, err := tx.NamedExec(`
        INSERT INTO events (group_id, created_by, title, content, event_start, event_end)
        VALUES (:group_id, :created_by, :title, :content, :event_start, :event_end)
    `, e)
	id, _ := res.LastInsertId()
	e.ID = int(id)

	return &e, nil
}

func GetEventsForUserID(db *sqlx.DB, id int) (*[]Event, error) {
	var events []Event
	query := `
        SELECT e.* FROM events e
        JOIN group_members gm ON e.group_id = gm.group_id
        WHERE gm.user_id = ? 
            AND gm.status = 'joined'
            AND e.event_end >= ?  
        ORDER BY e.event_end ASC
    `
	currentTime := time.Now().UTC()
	fmt.Println(currentTime)
	err := db.Select(&events, query, id, currentTime)
	if err != nil {
		return nil, err
	}

	return &events, nil
}

func GetEventsForGroupID(db *sqlx.DB, id int) ([]Event, error) {
	var events []Event
	query := `SELECT * from events WHERE group_id = ? ORDER BY event_end DESC`
	err := db.Select(&events, query, id)
	if err != nil {
		return nil, err
	}

	for i := range events {
		profile, _ := GetPrivateProfile(db, events[i].CreatedBy)
		events[i].User = profile
		responses, _ := GetEventResponsesForEventID(db, events[i].ID)
		fmt.Println(responses)
		events[i].EventResponses = responses
	}
	return events, nil
}

func InsertEventResponse(db *sqlx.DB, eventID int, response string, userID int) error {
	query := `
        INSERT OR REPLACE INTO event_responses (event_id, user_id, response)
        VALUES (?, ?, ?)
    `

	_, err := db.Exec(query, eventID, userID, response)
	if err != nil {
		log.Printf("error inserting event response: %v", err)
		return err
	}

	return nil
}

func GetEventResponsesForEventID(db *sqlx.DB, id int) ([]EventResponse, error) {
	var eventResponses []EventResponse
	fmt.Println(id)
	query := `SELECT * from event_responses WHERE event_id = ?`
	err := db.Select(&eventResponses, query, id)
	if err != nil {
		return nil, err
	}
	for i := range eventResponses {
		user, _ := GetPrivateProfile(db, eventResponses[i].UserID)
		eventResponses[i].User = user
	}
	return eventResponses, nil
}
