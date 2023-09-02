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
		fmt.Println("?")
	case "accept":
		models.MarkNotificationAsRead(s.db.DB, notif.ID)
		switch notif.Type {
		case "follow_accept":
			break
		case "follow_request":
			models.AcceptFollowRequest(s.db.DB, id, notif.SenderID)
		case "group_join_request":
			models.AcceptGroupJoinRequest(s.db.DB, id, notif.SenderID, notif.GroupID)
		}
		break
	case "reject":
		models.MarkNotificationAsRead(s.db.DB, notif.ID)
		switch notif.Type {
		case "follow_request":
			models.RejectFollowRequest(s.db.DB, id, notif.SenderID)
		case "group_join_request":
			models.RejectGroupJoinRequest(s.db.DB, id, notif.SenderID, notif.GroupID)
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

	case "get_groups":
		groups, _ := models.GetAllGroups(s.db.DB)
		fmt.Println(groups)
		groupBytes, _ := json.Marshal(groups)

		groupRaw := json.RawMessage(groupBytes)

		response := WebSocketMessage{
			Action: "get_groups",
			Data:   groupRaw,
		}
		s.writeMu.Lock()
		conn.WriteJSON(response)
		s.writeMu.Unlock()

	case "new_message":
		var req models.Message
		decoder := json.NewDecoder(bytes.NewReader(msg.Data))
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("new_message err: %v", err)
		}
		resultid, err := models.InsertMessage(s.db.DB, id, req.RecipientID, req.Content)
		if err != nil {
			log.Printf("insert message err: %v", err)
		}

		req.SenderID = id
		req.ID = int(resultid)
		s.broadcast <- req
		break
	}
}
