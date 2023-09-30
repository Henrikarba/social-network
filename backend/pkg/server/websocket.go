package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/models"
	mw "social-network/pkg/server/middleware"

	"github.com/gorilla/websocket"
)

type WebSocketMessage struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

func (s *Server) ListenForMessages(messages <-chan models.Message) {
	for m := range messages {
		conn, ok := s.conns[m.RecipientID]
		if ok {
			// profile, _ := models.GetPrivateProfile(s.db.DB, m.SenderID)
			// m.Sender = profile
			// mBytes, err := json.Marshal(m)
			// if err != nil {
			// 	log.Printf("Error marshaling message: %v", err)
			// 	continue
			// }
			// mRaw := json.RawMessage(mBytes)
			// msg := WebSocketMessage{
			// 	Action: "new_message",
			// 	Data:   mRaw,
			// }
			data, _ := models.GetAuthenticatedUserDate(s.db.DB, m.RecipientID)
			s.writeMu.Lock()
			err := conn.WriteJSON(data)
			s.writeMu.Unlock()
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(mw.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("can't upgrade connection: %v", err)
		return
	}
	defer func() {
		s.RemoveConn(id)
		conn.Close()
	}()

	profile, err := models.GetAuthenticatedUserDate(s.db.DB, id)
	if err != nil {
		log.Printf("can't get user data: %v", err)
	}

	s.AddConn(id, conn)

	s.writeMu.Lock()
	conn.WriteJSON(profile)
	s.writeMu.Unlock()

	for {
		var msg WebSocketMessage
		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("id %d disconnected %v", id, err)
			break
		}
		fmt.Println("receiving")
		s.messageHandler(msg, id)
	}
}

func (s *Server) messageHandler(msg WebSocketMessage, id int) {
	var notif models.Notification
	if msg.Action == "accept" || msg.Action == "reject" {
		decoder := json.NewDecoder(bytes.NewReader(msg.Data))
		err := decoder.Decode(&notif)
		if err != nil {
			log.Printf("can't decode msg: %v", err)
		}
	}
	conn := s.GetConn(id)

	switch msg.Action {
	case "follow_request":
		var req models.User
		decoder := json.NewDecoder(bytes.NewReader(msg.Data))
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("getting user %v: ", err)
		}
		err = models.CreateFollowRequest(s.db.DB, id, req.ID)
		if err != nil {
			log.Printf("making request: %v", err)
		}
		msg := &models.Message{
			RecipientID: req.ID,
		}
		s.broadcast <- *msg
		break
	case "accept":
		models.MarkNotificationAsRead(s.db.DB, notif.ID)
		switch notif.Type {
		case "follow_accept":
			break
		case "follow_request":
			models.AcceptFollowRequest(s.db.DB, id, notif.SenderID)
			msg := &models.Message{
				RecipientID: notif.SenderID,
			}
			s.broadcast <- *msg
			break
		case "group_join_request":
			models.AcceptGroupJoinRequest(s.db.DB, id, notif.SenderID, notif.GroupID)
			msg := &models.Message{
				RecipientID: notif.SenderID,
			}
			s.broadcast <- *msg
			break
		case "group_join_invite":
			models.AcceptGroupJoinInvite(s.db.DB, id, notif.SenderID, notif.GroupID)
			msg := &models.Message{
				RecipientID: notif.SenderID,
			}
			s.broadcast <- *msg
			break
		}

		break
	case "reject":
		models.MarkNotificationAsRead(s.db.DB, notif.ID)
		switch notif.Type {
		case "follow_request":
			models.RejectFollowRequest(s.db.DB, id, notif.SenderID)
			msg := &models.Message{
				RecipientID: notif.SenderID,
			}
			s.broadcast <- *msg
			break
		case "group_join_request":
			models.RejectGroupJoinRequest(s.db.DB, id, notif.SenderID, notif.GroupID)
			msg := &models.Message{
				RecipientID: notif.SenderID,
			}
			s.broadcast <- *msg
			break
		case "group_join_invite":
			models.DeclineGroupJoinInvite(s.db.DB, id, notif.SenderID, notif.GroupID)
			msg := &models.Message{
				RecipientID: notif.SenderID,
			}
			s.broadcast <- *msg
			break

		}

		break
	case "toggle_privacy":
		models.TogglePrivacy(s.db.DB, id)
		break
	case "get_user":
		var req models.UserResponse
		decoder := json.NewDecoder(bytes.NewReader(msg.Data))
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("getting user %v: ", err)
		}
		profile, err := models.GetUserProfile(s.db.DB, id, req.User.ID)
		profBytes, _ := json.Marshal(profile)

		profRaw := json.RawMessage(profBytes)
		response := WebSocketMessage{
			Action: "get_user",
			Data:   profRaw,
		}
		s.writeMu.Lock()
		conn.WriteJSON(response)
		s.writeMu.Unlock()
		break
	case "get_chat":
		var req models.Message
		decoder := json.NewDecoder(bytes.NewReader(msg.Data))
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("getting chat %v: ", err)
		}

		chat := models.GetFullChat(s.db.DB, id, req.SenderID)
		ID, _ := models.GetChatRoomID(s.db.DB, id, req.SenderID)
		chat.ChatroomID = ID
		postBytes, _ := json.Marshal(chat)

		chatRaw := json.RawMessage(postBytes)
		response := WebSocketMessage{
			Action: "get_chat",
			Data:   chatRaw,
		}

		s.writeMu.Lock()
		conn.WriteJSON(response)
		s.writeMu.Unlock()
		break

	case "new_event":
		var req models.Event
		decoder := json.NewDecoder(bytes.NewReader(msg.Data))
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("creating new event %v: ", err)
		}
		req.CreatedBy = id
		event, err := models.InsertEvent(s.db.DB, req)
		fmt.Println(event)
		break
	case "get_group_chat":
		var req models.Message
		decoder := json.NewDecoder(bytes.NewReader(msg.Data))
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("getting group_chat %v: ", err)
		}
		chat, _ := models.GetGroupMessageHistory(s.db.DB, req.SenderID)
		ID, _ := models.GetGroupChatRoomID(s.db.DB, req.SenderID)
		fullchat := models.FullChat{
			ChatroomID: ID,
			Messages:   chat,
		}
		postBytes, _ := json.Marshal(fullchat)

		chatRaw := json.RawMessage(postBytes)
		response := WebSocketMessage{
			Action: "get_group_chat",
			Data:   chatRaw,
		}

		s.writeMu.Lock()
		conn.WriteJSON(response)
		s.writeMu.Unlock()
		break

	case "event_response":
		var req models.EventResponse
		decoder := json.NewDecoder(bytes.NewReader(msg.Data))
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("getting event_response %v: ", err)
		}
		err = models.InsertEventResponse(s.db.DB, req.EventID, req.Response, id)
		if err != nil {
			log.Printf("inserting eventresponse: %v", err)
		}
		break

	case "get_groups":
		groups, _ := models.GetAllGroups(s.db.DB)
		groupBytes, _ := json.Marshal(groups)

		groupRaw := json.RawMessage(groupBytes)

		response := WebSocketMessage{
			Action: "get_groups",
			Data:   groupRaw,
		}
		s.writeMu.Lock()
		conn.WriteJSON(response)
		s.writeMu.Unlock()
		break

	case "get_group":
		var req models.Group
		decoder := json.NewDecoder(bytes.NewReader(msg.Data))
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("getting group %v: ", err)
		}
		group := models.GetGroup(s.db.DB, req.ID, id)
		groupBytes, _ := json.Marshal(group)

		groupRaw := json.RawMessage(groupBytes)

		response := WebSocketMessage{
			Action: "get_group",
			Data:   groupRaw,
		}
		s.writeMu.Lock()
		conn.WriteJSON(response)
		s.writeMu.Unlock()
		break

	case "join_group":
		var req models.Group
		decoder := json.NewDecoder(bytes.NewReader(msg.Data))
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("joining group %v: ", err)
		}
		models.JoinGroup(s.db.DB, id, req.CreatorID, req.ID)
		msg := &models.Message{
			RecipientID: req.CreatorID,
		}
		s.broadcast <- *msg
		break

	case "group_join_invite":
		var req models.GroupMember
		decoder := json.NewDecoder(bytes.NewReader(msg.Data))
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("inviting to group %v: ", err)
		}
		models.InviteToGroup(s.db.DB, req.GroupID, req.UserID, id)
		msg := &models.Message{
			RecipientID: req.UserID,
		}
		s.broadcast <- *msg
		break

	case "new_message":
		var req models.Message
		decoder := json.NewDecoder(bytes.NewReader(msg.Data))
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("new_message err: %v", err)
		}
		if req.Type == "regular" {
			resultid, err := models.InsertMessage(s.db.DB, id, req.RecipientID, req.Content)
			if err != nil {
				log.Printf("insert message err: %v", err)
			}
			req.SenderID = id
			req.ID = int(resultid)
			s.broadcast <- req
		} else {
			members := models.GetGroupMembers(s.db.DB, req.RecipientID)
			err := models.InsertGroupMessage(s.db.DB, id, req.RecipientID, req.Content)
			if err != nil {
				log.Printf("inserting group msg: %v", err)
			}
			profile, _ := models.GetPrivateProfile(s.db.DB, id)
			req.CreatedBy = profile
			req.SenderID = id

			for _, mem := range members {
				if mem == id {
					continue
				}
				conn, ok := s.conns[mem]
				if ok {
					s.writeMu.Lock()
					conn.WriteJSON(req)
					s.writeMu.Unlock()
				}
			}
		}

		break

	}
}
