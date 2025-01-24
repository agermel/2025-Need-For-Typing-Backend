package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"type/models"
	"type/utils"

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
	Client    models.Client
)

// @Summary      WebSocket 连接处理
// @Description  建立 WebSocket 连接并处理房间逻辑
// @Tags         WebSocket
// @Param        Authorization  header   string  true  "JWT Token"
// @Param        room_id        query    string  false "房间 ID"
// @Success      200            {object}  map[string]interface{}  "返回房间加入信息"
// @Failure      401            {string}  string                  "未授权"
// @Router       /ws [get]
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	var claims *utils.Claims

	// 从请求头中提取JWT
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization token is required", http.StatusUnauthorized)
		return
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		fmt.Println("Failed to fetch claims")
		return
	}

	Client.Username = claims.Username

	roomID := r.URL.Query().Get("room_id")

	// 房间号处理逻辑
	if roomID == "" || len(rooms[roomID]) >= 2 {
		fmt.Println("没输入房间号或房间已满，正在切换到其他房间")
		roomID = utils.GenerateRoomID()
		for len(rooms[roomID]) >= 2 {
			roomID = utils.GenerateRoomID()
		}

		// 获取当前URL
		newURL := r.URL

		// 获取URL查询参数
		query := newURL.Query()

		// 更新房间号
		query.Set("room_id", roomID)

		// 设置更新后的查询参数到 URL 中
		newURL.RawQuery = query.Encode()

		http.Redirect(w, r, newURL.String(), http.StatusFound)
	}

	Client.Conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		return
	}
	defer Client.Conn.Close()

	roomsLock.Lock()
	if rooms[roomID] == nil {
		rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	if len(rooms[roomID]) >= 2 {
		Client.Conn.WriteJSON(map[string]string{
			"error": "the room" + string(roomID) + " is full",
		})
		// 如果检测到房间满了就直接关闭client.conn
		Client.Conn.Close()
	}
	rooms[roomID][Client.Conn] = true
	roomsLock.Unlock()

	// 进入房间有json放回提示
	Client.Conn.WriteJSON(map[string]string{
		"Client": Client.Username,
		"status": "joined",
	})
	fmt.Printf("Client %s joined room %s\n", roomID, Client.Username)

	init_score := 0

	for {
		_, score, err := Client.Conn.ReadMessage()
		strScore := string(score)
		intScore, _ := strconv.Atoi(strScore)
		init_score += intScore
		if err != nil {
			Client.Conn.WriteJSON(map[string]string{
				"Client": Client.Username,
				"status": "exited",
			})
		}

		if len(rooms[roomID]) == 1 {
			for conn := range rooms[roomID] {
				conn.WriteJSON(map[string]int{
					"1p": claims.UserID,
				})
			}
		}

		if len(rooms[roomID]) == 2 {
			for conn := range rooms[roomID] {
				conn.WriteJSON(map[string]interface{}{
					"1p": claims.UserID,
					"2p": "exists",
				})
			}
		}

		BroadcastToRoom(roomID, init_score, claims.UserID)
	}
}

func BroadcastToRoom(roomID string, score int, userID int) {
	roomsLock.Lock()
	defer roomsLock.Unlock()

	for conn := range rooms[roomID] {
		strScore := strconv.Itoa(score)
		err := conn.WriteJSON(map[string]interface{}{
			"userID": userID,
			"score":  strScore,
		})
		if err != nil {
			fmt.Println("Failed to send score:", err)
			conn.Close()
			delete(rooms[roomID], conn)
		}
	}
}
