package controllers

import (
	"fmt"
	"net/http"
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
	scores    = make(map[string]map[int]int) // 存储房间内玩家得分
	roomsLock = sync.Mutex{}
)

func HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	fmt.Println(token)

	claims, err := utils.ParseToken(token)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	id := claims.UserID

	roomID := c.Query("room_id")
	if roomID == "" {
		roomID = utils.GenerateRoomID()
		for len(rooms[roomID]) >= 2 {
			roomID = utils.GenerateRoomID()
		}
		broadcastMessage(roomID, response.BroadcastBehave{
			Code:   response.ROOM_CREATED,
			Event:  "create and enter new room",
			UserID: id,
			RoomID: roomID,
			Score:  scores[roomID][id], // 发送当前玩家的分数
		})
	}

	if len(rooms[roomID]) > 2 {
		response.FailWithMessage("Room is full", c)
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
		scores[roomID] = make(map[int]int) // 初始化房间分数记录
	}
	rooms[roomID][conn] = true
	if _, exists := scores[roomID][id]; !exists {
		scores[roomID][id] = 0 // 初始化玩家分数
	}
	roomsLock.Unlock()

	if len(rooms[roomID]) == 1 {
		broadcastMessage(roomID, response.BroadcastBehave{
			Code:   response.ROOM_CREATED,
			Event:  "create and enter new room",
			UserID: id,
			RoomID: roomID,
			Score:  scores[roomID][id], // 发送当前玩家的分数
		})
	} else {
		// 通知新玩家加入
		broadcastMessage(roomID, response.BroadcastBehave{
			Code:   response.NEW_PLAYER_ENTER,
			Event:  "enter room",
			UserID: id,
			RoomID: roomID,
			Score:  scores[roomID][id], // 发送当前玩家的分数
		})
	}
	fmt.Printf("Client %d joined room %s\n", id, roomID)

	// 监听消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		broadcastMessage(roomID, message)
	}

	// 处理玩家退出
	roomsLock.Lock()
	delete(rooms[roomID], conn)
	delete(scores[roomID], id) // 移除玩家分数
	isRoomEmpty := len(rooms[roomID]) == 0
	roomsLock.Unlock()

	// 通知其他玩家该用户已退出
	broadcastMessage(roomID, response.BroadcastBehave{
		Code:   response.PLAYER_EXIT,
		Event:  "player exited",
		UserID: id,
		RoomID: roomID,
	})

	// 如果房间没人了，删除房间
	if isRoomEmpty {
		roomsLock.Lock()
		delete(rooms, roomID)
		delete(scores, roomID) // 删除房间分数数据
		roomsLock.Unlock()
		fmt.Printf("Room %s is now empty and deleted\n", roomID)
	}
}

func broadcastMessage(roomID string, message interface{}) {
	roomsLock.Lock()
	defer roomsLock.Unlock()

	for conn := range rooms[roomID] {
		if err := conn.WriteJSON(message); err != nil {
			fmt.Println("Failed to send message:", err)
			conn.Close()
			delete(rooms[roomID], conn)
		}
	}
}
