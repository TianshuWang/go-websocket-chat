package internal

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go-chat/model"
	"go-chat/utils"
	"strings"
)

type UserChat struct {
	Channel    *Channel
	Username   string
	Connection *websocket.Conn
}

func NewUserChat(channel *Channel, username string, conn *websocket.Conn) *UserChat {
	return &UserChat{
		Channel:    channel,
		Username:   username,
		Connection: conn,
	}
}

func (u *UserChat) OnlineListen() {
	for {
		_, message, err := u.Connection.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message: [%s]\n", err.Error())
			break
		}

		msg := new(model.Message)
		if err := json.Unmarshal(message, msg); err != nil {
			fmt.Printf("Error unmarshalling message: [%s]\n", err.Error())
			break
		}

		if strings.TrimSpace(msg.Sender) != strings.TrimSpace(u.Username) {
			msg.Sender = u.Username
		}
		fmt.Printf("Message: [%+v]\n", msg)
		u.Channel.messageChannel <- msg
	}

	u.Channel.leaveChannel <- u
}

func (u *UserChat) SendMessageToClient(message *model.Message) error {
	message.ID = utils.GenerateRandomID()
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message to json: [%v]\n", err)
	}
	err = u.Connection.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return fmt.Errorf("error writing message to json: [%v]\n", err)
	}
	fmt.Printf("Message sended from: [%s] to: [%s]\n", message.Sender, message.Target)
	return nil
}
