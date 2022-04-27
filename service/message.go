package service

const (
	msgTypeHeartbeat = iota
	msgTypeInit
	msgTypeCmd
	msgTypeChat
)

// Message format for frontend
type Message struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
	User    string `json:"user"`
}

var msgChan = make(chan Message)
