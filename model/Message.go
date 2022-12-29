package model

type Message struct {
	ID     int    `json:"id"`
	Sender string `json:"sender"`
	Target string `json:"target"`
	Body   string `json:"body"`
}

func NewMessage(sender string, target string, body string) *Message {
	return &Message{
		Sender: sender,
		Target: target,
		Body:   body,
	}
}
