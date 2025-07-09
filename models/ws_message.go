package models


type WSMessage struct {
	Type    string `json:"type" binding:"required"`    // dm | group | broadcast
	From    string `json:"from" binding:"required"`
	To      string `json:"to" binding:"required"`      // user or group
	Content string `json:"content" binding:"required"`
}