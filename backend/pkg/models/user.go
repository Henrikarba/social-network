package models

import (
	"database/sql"
	"fmt"
	"log"
	"social-network/pkg/utils"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int     `json:"id,omitempty" db:"id"`
	Email       string  `json:"email,omitempty" db:"email"`
	FirstName   string  `json:"first_name,omitempty" db:"first_name"`
	Password    string  `json:"-" db:"password"`
	LastName    string  `json:"last_name,omitempty" db:"last_name" `
	DateOfBirth string  `json:"date_of_birth,omitempty" db:"date_of_birth"`
	Avatar      *string `json:"avatar,omitempty" db:"avatar"`
	Nickname    string  `json:"nickname,omitempty" db:"nickname"`
	AboutMe     string  `json:"about_me,omitempty" db:"about_me"`
	Privacy     int     `json:"privacy" db:"privacy"`
	CreatedAt   string  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   string  `json:"updated_at,omitempty" db:"updated_at"`

	Status string `json:"status,omitempty"`
}

type UserRequest struct {
	Email       string `json:"email" schema:"email"`
	Password    string `json:"password" schema:"password"`
	Password2   string `schema:"password2"`
	FirstName   string `schema:"first_name"`
	LastName    string `schema:"last_name"`
	DateOfBirth string `schema:"date_of_birth"`
	ImageData   []byte `schema:"image"`
	Nickname    string `schema:"nickname"`
	AboutMe     string `schema:"about_me"`
}

func GetUserProfile(db *sqlx.DB, viewerID, profileID int) (*UserResponse, error) {
	var user User
	var response UserResponse
	var follower Follower

	query := `
		SELECT u.id, u.first_name, u.last_name, u.privacy, u.email, u.about_me, u.avatar, u.date_of_birth
		FROM users u
		WHERE u.id = ? AND (u.privacy = 0 OR u.privacy = 1)
	`

	err := db.Get(&user, query, profileID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if user.Privacy == 1 {
		followerQuery := `
			SELECT f.status
			FROM followers f
			WHERE f.follower_id = ? AND f.followee_id = ?
		`

		err = db.Get(&follower, followerQuery, viewerID, profileID)

		if err == sql.ErrNoRows || follower.Status != "accepted" {
			user.Email = ""
			user.DateOfBirth = ""
			user.Nickname = ""
			user.AboutMe = ""
			response.Following = nil
			response.Followers = nil
		} else {
			followers, _ := GetFollowers(db, user.ID)
			following, _ := GetFollowing(db, user.ID)
			response.Followers = &followers
			response.Following = &following
		}
	} else if user.Privacy == 0 {
		followers, _ := GetFollowers(db, user.ID)
		following, _ := GetFollowing(db, user.ID)
		response.Followers = &followers
		response.Following = &following
	}

	response.User = &user
	return &response, nil
}

func RegisterUser(db *sqlx.DB, Email, Password, FirstName, LastName, DateOfBirth, Nickname, AboutMe, MimeType string, ImageData []byte) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	newUser := &User{
		Email:       Email,
		Password:    Password,
		FirstName:   FirstName,
		LastName:    LastName,
		DateOfBirth: DateOfBirth,
		Nickname:    Nickname,
		AboutMe:     AboutMe,
	}

	if ImageData != nil {
		path, err := utils.SaveImage(ImageData, MimeType, "post")
		if err != nil {
			log.Println(err)
		}
		newUser.Avatar = &path
	} else {
		defaultPath := "profile/default.png"
		newUser.Avatar = &defaultPath
	}

	result, err := tx.NamedExec(`
		INSERT INTO users (email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, privacy, created_at, updated_at)
		VALUES (:email, :password, :first_name, :last_name, :date_of_birth, :avatar, :nickname, :about_me, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, newUser)

	if err != nil {
		return err
	}

	UserID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	newUser.ID = int(UserID)
	return nil
}

func ValidateLogin(db *sqlx.DB, email, password string) (bool, int, error) {
	var user User
	err := db.Get(&user, "SELECT id, password FROM users WHERE lower(email) = lower(?)", email)
	if err != nil {
		return false, 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, 0, err
	}
	return true, user.ID, nil
}

type UserResponse struct {
	User          *User            `json:"user,omitempty"`
	Followers     *[]User          `json:"followers,omitempty"`
	Following     *[]User          `json:"following,omitempty"`
	Notifications *[]Notification  `json:"notifications,omitempty"`
	Events        *[]Event         `json:"events"`
	Feed          *Feed            `json:"feed,omitempty"`
	Groups        *[]Group         `json:"groups,omitempty"`
	Messages      *[]Message       `json:"messages,omitempty"`
	GroupMessages *[]Message       `json:"group_messages,omitempty"`
	Chatlist      *[]ChatListEntry `json:"chatlist,omitempty"`
}

type Feed struct {
	Posts      *[]PostResponse `json:"posts"`
	GroupPosts *[]PostResponse `json:"group_posts"`
}

func GetAuthenticatedUserDate(db *sqlx.DB, id int) (*UserResponse, error) {
	var profile UserResponse

	// Profile
	user, err := GetPublicProfile(db, id)
	if err != nil {
		return nil, err
	}
	profile.User = user

	// Connections
	followers, err := GetFollowers(db, id)
	profile.Followers = &followers

	following, err := GetFollowing(db, id)
	profile.Following = &following

	// Notifications
	notifications, err := GetNotifications(db, id)
	for i := range notifications {
		sender, err := GetPrivateProfile(db, notifications[i].SenderID)
		if err != nil {
			log.Printf("cant get private profile: %v", err)
		}
		notifications[i].Sender = sender

		if notifications[i].GroupID > 0 {
			group := GetGroup(db, notifications[i].GroupID, id)
			notifications[i].Group = group
		}

	}

	profile.Notifications = &notifications

	// Post feed
	posts, err := GetRegularPosts(db, id)
	if err != nil {
		log.Printf("cant get regular posts: %v", err)
	}
	groupPosts, err := GetGroupPosts(db, id)
	if err != nil {
		log.Printf("cant get group posts: %v", err)
	}

	profile.Feed = &Feed{
		Posts:      &posts,
		GroupPosts: &groupPosts,
	}

	// Groups
	groups, err := GetUserGroups(db, id)
	profile.Groups = &groups

	unread, _ := GetLastUnreadMessages(db, id)
	profile.Messages = &unread

	unread2, _ := GetLastUnreadGroupMessages(db, id)
	profile.GroupMessages = &unread2

	chatlist := GetChatList(db, id)
	profile.Chatlist = chatlist

	events, err := GetEventsForUserID(db, id)
	if err != nil {
		log.Printf("get events: %v", err)
	}
	profile.Events = events

	return &profile, nil
}

func GetPublicProfile(db *sqlx.DB, id int) (*User, error) {
	var user User
	err := db.Get(&user, "SELECT id, email, first_name, last_name, about_me, date_of_birth, avatar, privacy FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetPrivateProfile(db *sqlx.DB, id int) (*User, error) {
	var user User
	err := db.Get(&user, "SELECT id, first_name, last_name, avatar FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetPrivacySetting(db *sqlx.DB, id int) (int, error) {
	var user User
	err := db.Get(&user, "SELECT privacy FROM users WHERE id = ?", id)
	if err != nil {
		return 0, err
	}
	return user.Privacy, nil
}

func GetUserGroups(db *sqlx.DB, userID int) ([]Group, error) {
	var groups []Group

	groupIDsQuery := `
        SELECT group_id, status
        FROM group_members
        WHERE user_id = ?`

	// Define a struct to hold the results of the query
	type GroupMember struct {
		GroupID int    `db:"group_id"`
		Status  string `db:"status"`
	}

	var groupMembers []GroupMember

	err := db.Select(&groupMembers, groupIDsQuery, userID)
	if err != nil {
		return nil, err
	}

	// Create a map to store status for each group
	groupStatusMap := make(map[int]string)

	// Populate the map with group statuses
	for _, member := range groupMembers {
		groupStatusMap[member.GroupID] = member.Status
	}

	// Fetch the groups based on the user's group IDs
	query := `
        SELECT id, title, creator_id
        FROM groups
        WHERE id IN (?)`

	// Prepare the list of group IDs from the map
	var groupIDs []int
	for groupID := range groupStatusMap {
		groupIDs = append(groupIDs, groupID)
	}

	query, args, err := sqlx.In(query, groupIDs)
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)
	err = db.Select(&groups, query, args...)
	if err != nil {
		return nil, err
	}

	for i, group := range groups {
		memberIDS := GetAnyGroupMembers(db, group.ID)
		groups[i].MemberIDS = memberIDS

		status, ok := groupStatusMap[group.ID]
		if ok {
			groups[i].Status = status
		}
	}

	return groups, nil
}

func TogglePrivacy(db *sqlx.DB, userid int) error {
	// Fetch the current privacy value
	var currentPrivacy int
	err := db.Get(&currentPrivacy, "SELECT privacy FROM users WHERE id = ?", userid)
	if err != nil {
		return err
	}

	newPrivacy := 1 - currentPrivacy
	_, err = db.Exec("UPDATE users SET privacy = ? WHERE id = ?", newPrivacy, userid)
	if err != nil {
		return err
	}

	return nil
}
