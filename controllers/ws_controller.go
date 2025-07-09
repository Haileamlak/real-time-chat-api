package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/haileamlak/chat-system/infrastructure"
	"github.com/haileamlak/chat-system/models"
	"github.com/haileamlak/chat-system/usecases"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketController interface {
	WebSocketHandler(c *gin.Context)
	handleConnection(conn *websocket.Conn, username string)
	handleMessage(msg models.WSMessage)
	StartSubscriber()
	SubscribeUser(username string)
	SubscribeGroup(group string)
	subscribeChannel(channel string)
	deliverToClient(msg models.WSMessage)
}

type webSocketController struct {
	msgUseCase   usecases.MessageUseCase
	redisService infrastructure.RedisService
	upgrader     websocket.Upgrader
	Clients      map[string]*websocket.Conn
}

func NewWebSocketController(messageUseCase usecases.MessageUseCase, redisService infrastructure.RedisService) WebSocketController {
	controller := &webSocketController{msgUseCase: messageUseCase, redisService: redisService, upgrader: websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// origin := r.Header.Get("Origin")
			// if origin == "ws://localhost:8080" {
			// 	return true
			// }
			return true
		},
	},
		Clients: make(map[string]*websocket.Conn),
	}

	controller.StartSubscriber()
	return controller
}

func (wsc *webSocketController) WebSocketHandler(c *gin.Context) {
	username := c.Query("user")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user param"})
		return
	}

	conn, err := wsc.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	// Subscribe user and groups via usecase
	wsc.SubscribeUser(username)
	groups, _ := wsc.redisService.GetClient().SMembers(context.Background(), "user:"+username+":groups").Result()
	for _, g := range groups {
		wsc.SubscribeGroup(g)
	}

	wsc.Clients[username] = conn
	log.Println(username, "connected via WebSocket")

	go wsc.handleConnection(conn, username)
}

func (wsc *webSocketController) handleConnection(conn *websocket.Conn, username string) {
	defer func() {
		conn.Close()
		delete(wsc.Clients, username)
		log.Println(username, "disconnected")
	}()

	for {
		var msg models.WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading:", err)
			return
		}

		wsc.handleMessage(msg)
	}

}

func (wsc *webSocketController) handleMessage(msg models.WSMessage) {
	msgJSON, _ := json.Marshal(msg)

	switch msg.Type {
	case "dm":
		directMessage := models.DirectMessage{
			From:    msg.From,
			To:      msg.To,
			Content: msg.Content,
		}

		wsc.msgUseCase.SaveDirectMessage(context.Background(), msg.From, msg.To, &directMessage)
		wsc.redisService.GetClient().Publish(context.Background(), "channel:user:"+msg.To, msgJSON)

	case "group":
		groupMessage := models.GroupMessage{
			From:    msg.From,
			Group:   msg.To,
			Content: msg.Content,
		}

		wsc.msgUseCase.SendGroupMessage(context.Background(), msg.To, &groupMessage)
		wsc.redisService.GetClient().Publish(context.Background(), "channel:group:"+msg.To, msgJSON)

	case "broadcast":
		broadcastMessage := models.BroadcastMessage{
			From:    msg.From,
			Content: msg.Content,
		}
		wsc.msgUseCase.SendBroadcastMessage(context.Background(), &broadcastMessage)
		wsc.redisService.GetClient().Publish(context.Background(), "channel:broadcast", msgJSON)
	}
}

func (wsc *webSocketController) StartSubscriber() {
	go wsc.subscribeChannel("channel:broadcast")
}

func (wsc *webSocketController) SubscribeUser(username string) {
	go wsc.subscribeChannel("channel:user:" + username)
}

func (wsc *webSocketController) SubscribeGroup(group string) {
	go wsc.subscribeChannel("channel:group:" + group)
}

func (wsc *webSocketController) subscribeChannel(channel string) {
	pubsub := wsc.redisService.GetClient().Subscribe(context.Background(), channel)
	ch := pubsub.Channel()

	log.Println("Subscribed to", channel)

	for msg := range ch {
		var wsMsg models.WSMessage
		err := json.Unmarshal([]byte(msg.Payload), &wsMsg)
		if err != nil {
			log.Println("Failed to unmarshal pubsub msg:", err)
			continue
		}

		log.Println("Received message on", channel, ":", wsMsg)

		wsc.deliverToClient(wsMsg)
	}
}

func (wsc *webSocketController) deliverToClient(msg models.WSMessage) {
	// if broadcast message, send to all clients except the sender
	if msg.Type == "broadcast" {
		for username, conn := range wsc.Clients {
			if username != msg.From {
				err := conn.WriteJSON(msg)
				if err != nil {
					log.Println("Error sending broadcast to", username, ":", err)
				}
			}
		}
		return
	} else {
		if conn, ok := wsc.Clients[msg.To]; ok {
			conn.WriteJSON(msg)
		}
	}
}
