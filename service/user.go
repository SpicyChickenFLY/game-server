package service

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

var userMap map[string]*User

func init() {
	userMap = make(map[string]*User)
}

func Login(userName string) error {
	if _, ok := userMap[userName]; ok {
		return errors.New("user name used")
	}
	userMap[userName] = newUser(userName)
	return nil
}

func Logout(userName string) {
	delete(userMap, userName)
}

func GetUser(userName string) *User {
	return userMap[userName]
}

// User Status
const (
	userInLobby = iota
	userInRoom
	userInGame
)

// User store infomation of user
type User struct {
	nickname  string
	status    int
	sendQueue chan Message
}

// newUser return *UserInfo
func newUser(nickName string) *User {
	return &User{
		nickName,
		userInLobby,
		nil,
	}
}

// setStatus is setter of UserInfo
func (u *User) setStatus(status int) {
	u.status = status
}

func (u *User) setMsgChan(sendMsgChan chan Message) {
	u.sendQueue = sendMsgChan
}

func (u *User) ReaderThread(wsConn *websocket.Conn, recvMsgChan chan Message, errChan chan error, exit chan struct{}) {
	for {
		select {
		case <-exit:
			return
		default:
			_, msg, err := wsConn.ReadMessage()
			if err != nil {
				log.Println(err)
				errChan <- err
				return
			}
			clientMsg := Message{}
			err = json.Unmarshal(msg, &clientMsg)
			if err != nil {
				log.Println(err)
				errChan <- err
				return
			}
			recvMsgChan <- clientMsg
		}
	}
}

func (u *User) WriterThread(wsConn *websocket.Conn, sendMsgChan chan Message, errChan chan error, exit chan struct{}) {
	for {
		select {
		case <-exit:
			return
		case msg := <-sendMsgChan: // send message
			msgStr, err := json.Marshal(msg)
			if err != nil {
				log.Println(err)
				errChan <- err
				return
			}
			err = wsConn.WriteMessage(websocket.TextMessage, msgStr)
			if err != nil {
				log.Println(err)
				errChan <- err
				return
			}
		}
	}
}
