package service

import (
	"errors"
)

var roomMap map[int]*RoomInfo
var roomIDIncVal = 0

func init() {
	roomMap = make(map[int]*RoomInfo)
	// FIXME: mocked data
	roomMap[1] = &RoomInfo{
		ID:         1,
		Name:       "123",
		Master:     "chow",
		PlayerList: []string{"chow", "sonong"},
		Capacity:   4,
	}
	roomMap[23123] = &RoomInfo{
		ID:         23123,
		Name:       "234234",
		Master:     "chow",
		PlayerList: []string{"chow", "sonong", "gg"},
		Capacity:   4,
	}
}

// ListRooms for player
func ListRooms() []*RoomInfo {
	roomList := make([]*RoomInfo, 0)
	for _, roomInfo := range roomMap {
		roomList = append(roomList, roomInfo)
	}
	return roomList
}

// CreateRoom by master player
func CreateRoom(roomInfo *RoomInfo) *RoomInfo {
	roomIDIncVal++
	roomInfo.ID = roomIDIncVal
	roomMap[roomInfo.ID] = roomInfo
	roomInfo.recvMsgChan = make(chan Message, 2)
	go roomInfo.HandleThread()
	return roomInfo
}

// DeleteRoom by master player or admin
func DeleteRoom(roomID, masterPlayerID int) {
	delete(roomMap, roomID)
}

// JoinRoom by player
func JoinRoom(roomID int, joinedUser string) (chan Message, chan Message, chan error, error) {
	room, ok := roomMap[roomID]
	// check if room exists
	if !ok {
		return nil, nil, nil, errors.New("room not exists")
	}
	// check if room not full
	if len(room.PlayerList) >= room.Capacity {
		return nil, nil, nil, errors.New("room is full")
	}
	// check if player already in another room
	if user, ok := userMap[joinedUser]; ok {
		if user.status == userInGame || user.status == userInRoom {
			// return nil, nil, nil, errors.New("player is already in another room")
		}
	} else {
		return nil, nil, nil, errors.New("player not exists")
	}
	// let player join
	room.PlayerList = append(room.PlayerList, joinedUser)
	// notify master and other players

	sendMsgChan := make(chan Message, 2)
	user, ok := userMap[joinedUser]
	user.setStatus(userInRoom)
	user.setMsgChan(sendMsgChan)
	roomQuit := make(chan error)
	return room.recvMsgChan, sendMsgChan, roomQuit, nil
}

// RoomInfo store room info
type RoomInfo struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Master      string   `json:"master"`
	PlayerList  []string `json:"players"`
	Capacity    int      `json:"capacity"`
	recvMsgChan chan Message
	game        *Game
}

// NewRoomInfo return *RoomInfo
func NewRoomInfo() *RoomInfo {
	ri := &RoomInfo{}
	return ri
}

// GetMasterPlayer is getter of RoomInfo
func (ri *RoomInfo) GetMasterPlayer() string {
	return ri.Master
}

func (r *RoomInfo) HandleThread() {
	for {
		select {
		case msg := <-r.recvMsgChan:
			switch msg.Type {
			case msgTypeHeartbeat:
				// reset ttl value
				// return response
			case msgTypeInit:

			case msgTypeCmd:
				// redirect msg to game

			case msgTypeChat:
				// send to all users
				for _, playerName := range r.PlayerList {
					player := userMap[playerName]
					player.sendQueue <- msg
				}
			}
		}
	}
}
