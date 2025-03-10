package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"type/api/response"
	"type/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	rooms     = make(map[string]map[*websocket.Conn]bool)
	roomsLock = sync.Mutex{}
)

// HandleWebSocket 处理 WebSocket 连接
func HandleWebSocket(c *gin.Context) {
	var claims *utils.Claims

	// 获取 JWT Token
	token := c.GetHeader("Authorization")
	if token == "" {
		response.FailWithMessage("Authorization token is required", c)
		return
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		response.FailWithMessage("Invalid token", c)
		return
	}

	roomID := c.Query("room_id")
	if roomID == "" || len(rooms[roomID]) >= 2 {
		roomID = utils.GenerateRoomID()
		for len(rooms[roomID]) >= 2 {
			roomID = utils.GenerateRoomID()
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/ws?room_id=%s", roomID))
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	roomsLock.Lock()
	if rooms[roomID] == nil {
		rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	if len(rooms[roomID]) >= 2 {
		conn.WriteJSON(gin.H{"error": "Room is full"})
		conn.Close()
	}
	rooms[roomID][conn] = true
	roomsLock.Unlock()

	conn.WriteJSON(gin.H{
		"Client": claims.Username,
		"status": "joined",
	})
	fmt.Printf("Client %s joined room %s\n", claims.Username, roomID)

	initScore := 0

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			conn.WriteJSON(gin.H{
				"Client": claims.Username,
				"status": "exited",
			})
			break
		}

		intScore, _ := strconv.Atoi(string(message))
		initScore += intScore

		if len(rooms[roomID]) == 1 {
			for conn := range rooms[roomID] {
				conn.WriteJSON(gin.H{"1p": claims.UserID})
			}
		} else if len(rooms[roomID]) == 2 {
			for conn := range rooms[roomID] {
				conn.WriteJSON(gin.H{
					"1p": claims.UserID,
					"2p": "exists",
				})
			}
		}
		BroadcastToRoom(roomID, initScore, claims.UserID)
	}
}

func BroadcastToRoom(roomID string, score int, userID int) {
	roomsLock.Lock()
	defer roomsLock.Unlock()

	for conn := range rooms[roomID] {
		if err := conn.WriteJSON(gin.H{
			"userID": userID,
			"score":  strconv.Itoa(score),
		}); err != nil {
			fmt.Println("Failed to send score:", err)
			conn.Close()
			delete(rooms[roomID], conn)
		}
	}
}
