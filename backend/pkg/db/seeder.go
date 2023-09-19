package db

import (
	"fmt"
	"log"
	"social-network/pkg/models"
	"social-network/pkg/utils"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/jmoiron/sqlx"
)

func tableHasRows(db *sqlx.DB, tableName string) (bool, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := db.Get(&count, query)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func Seeder(db *sqlx.DB) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	ok, _ := tableHasRows(db, "users")
	if !ok {
		for i := 0; i < 25; i++ {
			user := CreateUser(i)
			_, err := tx.NamedExec(`INSERT INTO users (
									email, password, first_name, last_name, date_of_birth, nickname,
									about_me, avatar, privacy, created_at, updated_at)
								VALUES (
									:email, :password, :first_name, :last_name,
									:date_of_birth, :nickname, :about_me, :avatar,
									:privacy, :created_at, :updated_at)`, &user)
			if err != nil {
				log.Println(err)
			}
		}
	}
	ok, _ = tableHasRows(db, "posts")
	if !ok {
		for i := 20; i >= 0; i-- {
			post := CreatePost(i)
			_, err := tx.NamedExec(`INSERT INTO posts (
                            		user_id, title, content, image_url, privacy, created_at, updated_at) 
                        		VALUES (
                            		:user_id, :title, :content, :image_url, 
                           			:privacy, :created_at, :updated_at)`, &post)
			if err != nil {
				log.Println(err)
			}
		}

		for i := 0; i < 5; i++ {
			comment := CreatePostComment(i)
			_, err := tx.NamedExec(`INSERT INTO comments (
									post_id, user_id, content, created_at, updated_at)
								VALUES (
									:post_id, :user_id, :content, :created_at, :updated_at)`, &comment)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	ok, _ = tableHasRows(db, "groups")
	if !ok {
		notif := models.Notification{}

		for i := 0; i < 10; i++ {
			group := CreateGroup(i)
			_, err := tx.NamedExec(`INSERT INTO groups (
                            		creator_id, title, description, created_at, updated_at) 
                        		VALUES (
                            		:creator_id, :title, :description, :created_at, :updated_at)`, &group)
			if err != nil {
				log.Println(err)
			}

		}
		for i := 0; i < 20; i++ {
			groupMember := CreateGroupMember(i)
			_, err := tx.NamedExec(`INSERT INTO group_members (
                            		group_id, user_id, status, invited_by, joined_at) 
                        		VALUES (
                            		:group_id, :user_id, :status, :invited_by, :joined_at)`, &groupMember)
			if err != nil {
				log.Println(err)
			}
			if i == 1 || i == 2 {
				notif.GroupID = groupMember.GroupID
				notif.UserID = groupMember.GroupID
				notif.SenderID = groupMember.UserID
				notif.Type = "group_join_request"
				notif.Message = "wants to join your group"
				_, err = tx.NamedExec(`INSERT INTO user_notifications (user_id, sender_id, message, type, group_id)
									VALUES(:user_id, :sender_id, :message, :type, :group_id)
										`, &notif)
			}
		}

	}
	ok, _ = tableHasRows(db, "followers")
	if !ok {

		for i := 0; i < 6; i++ {
			follower := CreateFollowers(i)
			_, err = tx.NamedExec(`INSERT INTO followers (
									follower_id, followee_id, status) 
								VALUES (
									:follower_id, :followee_id, :status)`, &follower)
			var msg string
			var ntype string
			if i < 2 {
				msg = "wants to follow you"
				ntype = "follow_request"
			} else {
				msg = "is now following you"
				ntype = "follow_accept"
			}

			notification := models.Notification{
				UserID:   follower.FolloweeID, // followee_id as user_id
				SenderID: follower.FollowerID, // follower_id as sender_id
				Message:  msg,
				Type:     ntype,
			}

			_, err = tx.NamedExec(`INSERT INTO user_notifications (
                                user_id, sender_id, message, type)
                            VALUES (
                                :user_id, :sender_id, :message, :type)`, &notification)
			if err != nil {
				log.Println(err)
			}
		}
	}

	ok, _ = tableHasRows(db, "group_posts")
	if !ok {
		for i := 8; i >= 0; i-- {
			groupPost := CreateGroupPost(i)
			_, err := tx.NamedExec(`INSERT INTO group_posts (
                            		group_id, user_id, title, content, image_url, created_at, updated_at) 
                        		VALUES (
                            		:group_id, :user_id, :title, :content, :image_url, 
                           			:created_at, :updated_at)`, &groupPost)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		for i := 0; i < 9; i++ {
			groupPostComment := CreateGroupPostComment(i)
			_, err := tx.NamedExec(`INSERT INTO group_comments (
									post_id, user_id, content, created_at, updated_at)
								VALUES (
									:post_id, :user_id, :content, :created_at, :updated_at)`, &groupPostComment)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	ok, _ = tableHasRows(db, "messages")
	if !ok {
		chats := CreateChat(20)
		err = models.CreateChatRoomForUsers(db, 1, 2)

		for i := range chats {
			_, err = tx.NamedExec(`
        INSERT INTO messages (sender_id, recipient_id, content, chatroom_id, created_at)
        VALUES (:sender_id, :recipient_id, :content, :chatroom_id, :created_at)
   			`, chats[i])
		}

		cr := models.Chatroom{ID: 1, Name: "1-2"}
		_, err = tx.NamedExec(`INSERT INTO chatrooms (name) VALUES (:name)`, cr)
		if err != nil {
			return err
		}

		cp1 := models.ChatroomParticipant{ChatroomID: 1, UserID: 1}
		_, err = tx.NamedExec(`INSERT INTO chatroom_participants (chatroom_id, user_id) VALUES (:chatroom_id, :user_id)`, cp1)
		if err != nil {
			return err
		}

		cp2 := models.ChatroomParticipant{ChatroomID: 1, UserID: 2}
		_, err = tx.NamedExec(`INSERT INTO chatroom_participants (chatroom_id, user_id) VALUES (:chatroom_id, :user_id)`, cp2)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}

	}
	tx.Commit()
	return nil
}

func CreateUser(i int) models.User {
	var (
		password  string
		email     string
		firstName string
		lastName  string
		nickname  string
		aboutme   string
		privacy   int
		avatar    string
	)

	dateOfBirth := faker.Date()
	createdAt := time.Date(2020, time.January, i+20, i+5, i+5, i+5, i, time.UTC)
	if (i+1)%2 == 0 {
		avatar = "profile/default2.png"
	} else {
		avatar = "profile/default.png"
	}
	if i == 0 || i == 1 || i == 2 {
		switch i {
		case 0:
			email = "admin"
			firstName = "Linus"
			lastName = "Torvalds"
			privacy = 0
		case 1:
			email = "admin2"
			firstName = "Bill"
			lastName = "Gates"
			privacy = 1
		case 2:
			email = "admin3"
			firstName = "Jeff"
			lastName = "Bezos"
			privacy = 1
		}

		nickname = "Admin" + fmt.Sprint(i+1)
		password, _ = utils.HashPassword("secret")
		aboutme = fmt.Sprintf("Something is very interesting about me - I am %s %s of course", firstName, lastName)

	} else {
		if (i+1)%5 == 0 {
			privacy = 1
		}
		firstName = faker.FirstName()
		lastName = faker.LastName()
		email = faker.Email()
		nickname = faker.Username()
		password, _ = utils.HashPassword(faker.Password())
		aboutme = faker.Sentence()
	}

	return models.User{
		Email:       email,
		Password:    password,
		FirstName:   firstName,
		LastName:    lastName,
		DateOfBirth: dateOfBirth,
		Nickname:    nickname,
		AboutMe:     aboutme,
		Avatar:      &avatar,
		Privacy:     privacy,
		CreatedAt:   utils.FormatTime(&createdAt),
		UpdatedAt:   utils.FormatTime(&createdAt),
	}
}

func CreateFollowers(i int) models.Follower {
	followers := []models.Follower{
		{FollowerID: 1, FolloweeID: 3, Status: "pending"},
		{FollowerID: 2, FolloweeID: 3, Status: "pending"},
		{FollowerID: 1, FolloweeID: 2, Status: "accepted"},
		{FollowerID: 3, FolloweeID: 5, Status: "accepted"},
		{FollowerID: 3, FolloweeID: 10, Status: "accepted"},
		{FollowerID: 3, FolloweeID: 15, Status: "accepted"},
		{FollowerID: 3, FolloweeID: 20, Status: "accepted"},
	}
	return followers[i]
}

func CreatePost(i int) *models.PostResponse {
	var (
		userid    int
		title     string
		content   string
		imageUrl  *string
		privacy   int
		createdAt time.Time
	)
	baseDate := time.Date(2022, time.April, 22, 23, 55, 55, 0, time.UTC)
	createdAt = baseDate.AddDate(0, 0, -i) //
	titles := []string{
		"Linux has finally partnered with Microsoft",
		"Microsoft is happy to announce remake of MS-DOS",
		"Amazon is building a new factory in Amazon Rainforest",
	}
	contents := []string{
		// Post one
		`In a surprising turn of events, Linux and Microsoft have announced a groundbreaking partnership that 
		promises to reshape the tech landscape. The collaboration, unveiled at a joint press conference earlier 
		today, marks a significant shift in the relationship between the two tech giants, who have historically been seen as rivals.
		Satya Nadella, CEO of Microsoft, expressed his enthusiasm about the partnership, stating, "This collaboration is all 
		about bringing the best of both worlds to our users. We believe that by working together, we can drive innovation 
		forward in ways we couldn't have imagined separately."
		Linus Torvalds, the creator of Linux, added, "It's a new era for Linux and open-source software. 
		This partnership with Microsoft will allow us to reach more users and continue to champion the principles of open-source."
		While specific details of the partnership remain under wraps, insiders suggest that we can expect integrated tools, cross-platform compatibility, 
		and joint ventures in cloud computing. Both communities eagerly await the innovations this partnership will bring.`,
		// Post two
		`In a nostalgic nod to the past, Microsoft has unveiled plans to remake MS-DOS. "We're blending old-school charm with modern tech," 
		said a Microsoft spokesperson. This revamped version promises enhanced features while retaining the simplicity users loved. 
		Expected to launch next year, the tech community is buzzing with anticipation. Retro computing is back in style!`,
		// Post three
		`In a move that has raised eyebrows, Amazon announced plans to construct a state-of-the-art factory within the Amazon Rainforest. 
		The company states this initiative aims to boost local economies and introduce sustainable technologies. 
		"Our goal is to harmonize modern infrastructure with nature's wonders," said an Amazon representative. 
		While the project promises eco-friendly practices and job opportunities for locals, environmentalists are 
		closely monitoring its potential impact on the delicate ecosystem. The world watches as Amazon returns to its namesake.`,
	}
	switch {
	case i < 3:
		title = titles[i]
		content = contents[i]
		privacy = i
		url := fmt.Sprintf("post/default%d.png", i+1)
		imageUrl = &url
	case (i+1)%5 == 0:
		url := "post/default4.png"
		imageUrl = &url
		title = faker.Sentence()
		content = faker.Paragraph()
	case (i+1)%10 == 0:
		privacy = 1
		title = faker.Sentence()
		content = faker.Paragraph()
	default:
		title = faker.Sentence()
		content = faker.Paragraph()
	}
	userid = i + 1
	return &models.PostResponse{
		UserID:    userid,
		Title:     title,
		Content:   content,
		ImageUrl:  imageUrl,
		Privacy:   privacy,
		CreatedAt: utils.FormatTime(&createdAt),
		UpdatedAt: utils.FormatTime(&createdAt),
	}
}

func CreatePostComment(i int) models.Comment {
	var (
		post_id   int
		user_id   int
		content   string
		createdAt time.Time
	)
	createdAt = time.Date(2023, time.January, i+3, i+5, i+5, i+5, i, time.UTC)
	post_id = 21
	user_id = i + 4
	contents := []string{
		"Loving these news!",
		"Very disappointed!",
		"I am moving to macOS instead!",
		"Who cares?",
		"Windows XP was the best thing ever",
	}
	content = contents[i]
	return models.Comment{
		PostID:    post_id,
		UserID:    user_id,
		Content:   content,
		CreatedAt: utils.FormatTime(&createdAt),
		UpdatedAt: utils.FormatTime(&createdAt),
	}
}

func CreateGroup(i int) *models.Group {
	createdAt := time.Date(2021, time.November, i+20, i+5, i+5, i+5, i, time.UTC)
	titles := []string{
		"Support group for Windows users",
		"BTW I Use Arch",
		"Amazing Amazon",
		"JavaScript Enthusiasts Unite",
		"Journey with Java",
		"Pythonistas' Paradise",
		"Rust: Beyond the Basics",
		"Go-Getters of Golang",
		"Frontend Fanatics",
		"Deep Dive into Databases",
	}
	descriptions := []string{
		"A community for those who navigate the quirks of Windows. Share tips, tricks, and support.",
		"For those proud to announce their love for Arch Linux. Discuss setups, configurations, and more.",
		"Delve into the vast world of Amazon services, from AWS to Prime. Share experiences and advice.",
		"Join fellow JS developers in discussing the latest libraries, frameworks, and best practices.",
		"From beginners to experts, all Java lovers are welcome. Discuss projects, JVM, and more.",
		"Connect with Python developers. Discuss libraries, web frameworks, and Pythonic best practices.",
		"Explore the world of Rust. Share your experiences with memory safety and concurrency.",
		"Discuss the simplicity and efficiency of Go. Share projects, tools, and Goroutine tips.",
		"Web designers and developers unite! Discuss CSS tricks, JS animations, and responsive designs.",
		"From SQL to NoSQL, dive deep into database discussions. Share queries, optimizations, and tools.",
	}
	return &models.Group{
		CreatorID:   i + 1,
		Title:       titles[i],
		Description: descriptions[i],
		CreatedAt:   utils.FormatTime(&createdAt),
		UpdatedAt:   utils.FormatTime(&createdAt),
	}
}

func CreateGroupMember(i int) *models.GroupMember {
	var (
		groupID   int
		userID    int
		status    string
		invitedBy int
		joinedAt  time.Time
	)
	joinedAt = time.Date(2022, time.December, i+1, i+5, i+5, i+5, i, time.UTC)
	switch i {
	case 0:
		userID = 1
		groupID = 2
		status = "joined"
		invitedBy = 1
	case 1:
		userID = 3
		groupID = 1
		status = "requested"
		invitedBy = 3
	case 2:
		userID = 1
		groupID = 3
		status = "requested"
		invitedBy = userID
	case 10, 11, 12, 13:
		userID = i + 1
		groupID = 3
		status = "joined"
		invitedBy = userID
	case 14, 15, 16, 17, 18, 19:
		userID = i + 1
		groupID = 1
		status = "joined"
		invitedBy = userID
	default:
		userID = 5 + i
		groupID = 1
		status = "joined"
		invitedBy = groupID
	}
	return &models.GroupMember{
		GroupID:   groupID,
		UserID:    userID,
		Status:    status,
		InvitedBy: invitedBy,
		JoinedAt:  utils.FormatTime(&joinedAt),
	}
}

func CreateGroupPost(i int) models.PostResponse {
	var (
		userid    int
		groupid   int
		title     string
		content   string
		createdAt time.Time
	)
	baseDate := time.Date(2022, time.April, 30, 23, 55, 55, 0, time.UTC)
	createdAt = baseDate.AddDate(0, 0, -i) //
	titles := []string{
		"What do we think of the recent merge between Linux and Microsoft?",
		"What will happen to ArchLinux after recent news?",
		"How will Amazon's new factory in the Rainforest impact its global operations?",
	}
	contents := []string{
		`The tech world was taken by storm with the recent announcement of a merger between Linux and Microsoft. 
		Historically seen as rivals, this unexpected union promises a new era of collaboration and innovation. 
		While some purists express concerns over the potential dilution of Linux's open-source ethos, others are optimistic
		 about the integrated tools and solutions that could emerge. What are your thoughts on this groundbreaking development?`,
		`The ArchLinux community has been buzzing with speculation since the news of the Linux-Microsoft merger broke. 
		 As a distro known for its simplicity and user-centric approach, many wonder how ArchLinux will adapt to the changing landscape. 
		 Will we see a shift in its rolling release strategy? Or perhaps a deeper integration with Microsoft's tools? 
		 Join the discussion and share your predictions.`,
		`Amazon's decision to build a factory in the heart of the Amazon Rainforest has been met with a mix of intrigue and concern. 
		While the company promises eco-friendly practices and a boost to local economies, environmentalists are wary of the potential ecological impact. 
		On the business front, this move could position Amazon to tap into new markets and streamline its supply chain. 
		But at what cost? Dive into the debate and explore the potential ramifications on Amazon's global operations`,
	}
	userIDS := []int{10, 8, 9, 1, 11, 12, 13, 14, 20}
	image := ""
	switch i {
	case 0, 1, 2:
		if i == 1 {
			image = "post/default4.png"
		} else if i == 2 {
			image = "post/default3.png"
		} else {
			image = "post/default2.png"
		}
		userid = userIDS[i]
		groupid = i + 1
		title = titles[i]
		content = contents[i]
	default:
		userid = userIDS[i]
		groupid = 1
		title = faker.Sentence()
		content = faker.Paragraph()
	}
	return models.PostResponse{
		UserID:    userid,
		GroupID:   groupid,
		Title:     title,
		ImageUrl:  &image,
		Content:   content,
		CreatedAt: utils.FormatTime(&createdAt),
		UpdatedAt: utils.FormatTime(&createdAt),
	}
}

func CreateGroupPostComment(i int) models.Comment {
	var (
		post_id   int
		user_id   int
		content   string
		createdAt time.Time
	)
	createdAt = time.Date(2023, time.January, i+3, i+5, i+5, i+5, i, time.UTC)
	comments := []string{
		"Huge move! Didn't see this coming at all.",
		"Is this the end of open-source as we know it?",
		"I'm optimistic. Both have strengths to offer.",
		"Microsoft's strategy has been evolving lately.",
		"Linux purists won't be happy about this.",
		"Could be a game-changer for the tech industry.",
		"I hope they maintain Linux's core values.",
		"This might lead to some innovative products!",
		"I'm curious about the implications for developers.",
	}
	userIDs := []int{1, 8, 9, 10, 11, 12, 13, 14, 20}
	post_id = 9
	user_id = userIDs[i]
	content = comments[i]
	return models.Comment{
		PostID:    post_id,
		UserID:    user_id,
		Content:   content,
		CreatedAt: utils.FormatTime(&createdAt),
		UpdatedAt: utils.FormatTime(&createdAt),
	}
}

func CreateChat(count int) []*models.Message {
	messages := make([]*models.Message, 0)
	for i := 0; i < count; i++ {
		senderID := 1
		recipientID := 2

		if i%2 == 1 { // Alternate sender and recipient
			senderID, recipientID = recipientID, senderID
		}

		content := fmt.Sprintf("Message %d", i+1)

		createdAt := time.Date(2020, time.January, i+20, i+5, i+5, i+5, i, time.UTC)

		message := &models.Message{
			SenderID:    senderID,
			RecipientID: recipientID,
			Content:     content,
			CreatedAt:   utils.FormatTime(&createdAt),
			ChatroomID:  1,
		}

		messages = append(messages, message)
	}
	for i := 0; i < count; i++ {
		senderID := 1
		recipientID := 3 // Update recipient ID
		chatroomID := 2  // Update chatroom ID

		if i%2 == 1 { // Alternate sender and recipient
			senderID, recipientID = recipientID, senderID
		}

		content := fmt.Sprintf("Message %d", i+1)

		createdAt := time.Date(2020, time.January, i+20, i+5, i+5, i+5, i, time.UTC)

		message := &models.Message{
			SenderID:    senderID,
			RecipientID: recipientID,
			Content:     content,
			CreatedAt:   utils.FormatTime(&createdAt),
			ChatroomID:  chatroomID,
		}

		messages = append(messages, message)
	}

	return messages
}
