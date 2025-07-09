package models

type GroupMessage struct {
	From      string `json:"from" binding:"required"`
	Group     string `json:"group" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Timestamp string `json:"timestamp" binding:"required"`
}
