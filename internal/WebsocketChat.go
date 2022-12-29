package internal

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go-chat/model"
	"go-chat/utils"
	"net/http"
	"strings"
)

var upgradeWebsocket = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	fmt.Printf("Request's method: %s, host: %s. URI: %s, protocol: %s\n", r.Method, r.Host, r.RequestURI, r.Proto)
	return r.Method == http.MethodGet
}

type MessageChannel chan *model.Message
type UserChannel chan *UserChat

type Channel struct {
	messageChannel MessageChannel
	leaveChannel   UserChannel
}

type WebsocketChat struct {
	users       map[string]*UserChat
	joinChannel UserChannel
	channel     *Channel
}

func NewWebsocketChat() *WebsocketChat {
	return &WebsocketChat{
		users:       make(map[string]*UserChat),
		joinChannel: make(UserChannel),
		channel: &Channel{
			messageChannel: make(MessageChannel),
			leaveChannel:   make(UserChannel),
		},
	}
}

func (w *WebsocketChat) UserConnectionHandler(rw http.ResponseWriter, r *http.Request) {
	connection, err := upgradeWebsocket.Upgrade(rw, r, nil)
	if err != nil {
		fmt.Printf("Error connecting to %s, error: %v\n", r.Host, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	keys := r.URL.Query()
	username := strings.TrimSpace(keys.Get("username"))
	if username == "" {
		username = fmt.Sprintf("user-%d", utils.GenerateRandomID())
	}

	userChat := NewUserChat(w.channel, username, connection)

	w.joinChannel <- userChat
	userChat.OnlineListen()
}

func (w *WebsocketChat) UsersChatManager() {
	for {
		select {
		case userChat := <-w.joinChannel:
			w.AddUserChat(userChat)
		case message := <-w.channel.messageChannel:
			w.SendMessage(message)
		case userChat := <-w.channel.leaveChannel:
			w.LeaveUserChat(userChat.Username)
		}
	}
}

func (w *WebsocketChat) AddUserChat(userChat *UserChat) {
	if user, ok := w.users[userChat.Username]; ok {
		user.Connection = userChat.Connection
	} else {
		w.users[userChat.Username] = userChat
		fmt.Printf("New user: %s joined\n", userChat.Username)
	}
}

func (w *WebsocketChat) SendMessage(message *model.Message) {
	if targetChat, ok := w.users[message.Target]; ok {
		err := targetChat.SendMessageToClient(message)
		if err != nil {
			fmt.Printf("Error connecting to client: %s, error: %v\n", message.Target, err)
		}
	}
}

func (w *WebsocketChat) LeaveUserChat(username string) {
	if user, ok := w.users[username]; ok {
		defer user.Connection.Close()
		delete(w.users, username)
		fmt.Printf("User: %s left the chat\n", username)
	}
}
