package models

type BroadcastMessage struct {
	From      string `json:"from" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Timestamp string `json:"timestamp" binding:"required"`
}
