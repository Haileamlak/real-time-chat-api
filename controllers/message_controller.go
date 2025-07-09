package controllers

import (
	"github.com/haileamlak/chat-system/models"
	"github.com/haileamlak/chat-system/usecases"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MessageController interface {
	SendDM(c *gin.Context)
	GetDMHistory(c *gin.Context)
	CreateGroup(c *gin.Context)
	JoinGroup(c *gin.Context)
	SendGroupMessage(c *gin.Context)
	GetGroupHistory(c *gin.Context)
	SendBroadcast(c *gin.Context)
	GetBroadcastHistory(c *gin.Context)
}

type messageController struct{
	messageUseCase usecases.MessageUseCase
}

func NewMessageController(messageUseCase usecases.MessageUseCase) MessageController {
	return &messageController{
		messageUseCase: messageUseCase,
	}
}


func (m *messageController) SendDM(c *gin.Context) {
	var msg models.DirectMessage
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message"})
		return
	}

	if msg.From != c.GetString("user") {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot send messages as another user"})
		return
	}

	msg.Timestamp = time.Now().UTC().Format(time.RFC3339)
	err := m.messageUseCase.SaveDirectMessage(c.Request.Context(), msg.From, msg.To, &msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sent"})
}

func (m *messageController) GetDMHistory(c *gin.Context) {
	otherUser := c.Param("user")
	fromUser := c.Query("from")

	if otherUser == "" || fromUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both 'user' and 'from' parameters are required"})
		return
	}
	if fromUser != c.GetString("user") {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot access messages as another user"})
		return
	}

	messages, err := m.messageUseCase.GetDMHistory(c.Request.Context(), fromUser, otherUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (m *messageController) CreateGroup(c *gin.Context) {
	type Req struct {
		GroupName string `json:"group" binding:"required"`
		Creator   string `json:"creator" binding:"required"`
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil || req.GroupName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group data"})
		return
	}

	err := m.messageUseCase.CreateGroup(c.Request.Context(), req.GroupName, []string{req.Creator})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group created"})
}

func (m *messageController) JoinGroup(c *gin.Context) {
	type Req struct {
		GroupName string `json:"group" binding:"required"`
		User      string `json:"user" binding:"required"`
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if req.GroupName == "" || req.User == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group name and user are required"})
		return	
	}

	err := m.messageUseCase.AddMemberToGroup(c.Request.Context(), req.GroupName, req.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined group"})
}

func (m *messageController) SendGroupMessage(c *gin.Context) {
	var msg models.GroupMessage
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message"})
		return
	}

	if msg.From != c.GetString("user") {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot send messages as another user"})
		return		
	}

	msg.Timestamp = time.Now().UTC().Format(time.RFC3339)
	err := m.messageUseCase.SendGroupMessage(c.Request.Context(), msg.Group, &msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send group message"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Message sent"})
}

func (m *messageController) GetGroupHistory(c *gin.Context) {
	group := c.Param("name")
	if group == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group name is required"})
		return
	}

	if exists, err := m.messageUseCase.GroupExists(c.Request.Context(), group); err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	msgs, err := m.messageUseCase.GetGroupHistory(c.Request.Context(), group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve group messages"})
		return
	}

	c.JSON(http.StatusOK, msgs)
}

func (m *messageController) SendBroadcast(c *gin.Context) {
	var msg models.BroadcastMessage
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message"})
		return
	}

	msg.Timestamp = time.Now().UTC().Format(time.RFC3339)
	
	err := m.messageUseCase.SendBroadcastMessage(c.Request.Context(), &msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send broadcast message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Broadcast sent"})
}

func (m *messageController) GetBroadcastHistory(c *gin.Context) {
	msgs, err := m.messageUseCase.GetBroadcastHistory(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve broadcast messages"})
		return
	}

	c.JSON(http.StatusOK, msgs)
}
