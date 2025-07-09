package models

type DirectMessage struct {
	From      string `json:"from" binding:"required"`
	To        string `json:"to" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Timestamp string `json:"timestamp" binding:"required"`
}
