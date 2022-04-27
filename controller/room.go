package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/SpicyChickenFLY/game-server/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ListRooms for user
func ListRooms(c *gin.Context) {
	type roomList struct {
		RoomInfoList []*service.RoomInfo `json:"room_list"`
	}
	c.JSON(http.StatusOK, roomList{service.ListRooms()})
}

// CreateRoom by user
func CreateRoom(c *gin.Context) {
	// check login
	master := c.GetString("nickname")
	roomInfo := service.NewRoomInfo()
	// get params
	err := c.ShouldBindJSON(roomInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad request params"})
		return
	}
	if roomInfo.Master != master {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "wrong user"})
		return
	}

	service.CreateRoom(roomInfo)
	c.JSON(http.StatusOK, roomInfo)
}

// JoinRoom by user
func JoinRoom(c *gin.Context) {
	roomIDStr := c.Param("room-id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "room id must be a integer value")
		log.Println(http.StatusBadRequest, "room id must be a integer value")
		return
	}

	player := c.Query("nickname")
	recvMsgChan, sendMsgChan, userQuit, err := service.JoinRoom(roomID, player)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		log.Println(http.StatusInternalServerError, err.Error())
		return
	}

	wsUpgrader := websocket.Upgrader{
		CheckOrigin:  func(r *http.Request) bool { return true },
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}

	wsConn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": fmt.Sprintf("errors occurs when upgrading to ws %v", err)})
		log.Println(http.StatusInternalServerError, gin.H{"msg": fmt.Sprintf("errors occurs when upgrading to ws %v", err)})
		return
	}
	defer wsConn.Close()

	exit := make(chan struct{}, 2)
	user := service.GetUser(player)
	go user.ReaderThread(wsConn, recvMsgChan, userQuit, exit)
	go user.WriterThread(wsConn, sendMsgChan, userQuit, exit)

	select {
	case err := <-userQuit:
		log.Println(err)
		exit <- struct{}{}
		exit <- struct{}{}
		return
	} // blocked until quit

}
