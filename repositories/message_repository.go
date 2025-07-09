package repositories

import (
	"context"

	"github.com/haileamlak/chat-system/infrastructure"
	"github.com/haileamlak/chat-system/models"
	"encoding/json"
)

type MessageRepository interface {
	SaveDirectMessage(key string, msg *models.DirectMessage) error
	GetDirectMessage(key string, index int64) (*models.DirectMessage, error)
	GetDMHistory(key string) ([]*models.DirectMessage, error)
	CreateGroup(key string, members []string) error
	GetGroupHistory(group string) ([]*models.GroupMessage, error)
	GroupExists(key string) (bool, error)
	AddMemberToGroup(groupKey string, member string) error
	IsMemberOfGroup(groupKey string, member string) (bool, error)
	SendGroupMessage(key string, msg *models.GroupMessage) error
	SendBroadcastMessage(key string, msg *models.BroadcastMessage) error
	GetBroadcastHistory(key string) ([]*models.BroadcastMessage, error)
}

type messageRepository struct {
	redisService infrastructure.RedisService
}

func NewMessageRepository(redisService infrastructure.RedisService) MessageRepository {
	return &messageRepository{
		redisService: redisService,
	}
}

func (r *messageRepository) SaveDirectMessage(key string, msg *models.DirectMessage) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.redisService.GetClient().RPush(context.Background(), key, msgJSON).Err()
}

func (r *messageRepository) GetDirectMessage(key string, index int64) (*models.DirectMessage, error) {
	msgJSON, err := r.redisService.GetClient().LIndex(context.Background(), key, index).Result()
	if err != nil {
		return nil, err
	}

	var msg models.DirectMessage
	if err := json.Unmarshal([]byte(msgJSON), &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *messageRepository) GetDMHistory(key string) ([]*models.DirectMessage, error) {
	msgs, err := r.redisService.GetClient().LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var messages []*models.DirectMessage
	for _, msgJSON := range msgs {
		var msg models.DirectMessage
		if err := json.Unmarshal([]byte(msgJSON), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}

func (r *messageRepository) CreateGroup(groupKey string, members []string) error {
	return r.redisService.GetClient().SAdd(context.Background(), groupKey, members).Err()
}

func (r *messageRepository) GetGroupHistory(groupKey string) ([]*models.GroupMessage, error) {
	msgs, err := r.redisService.GetClient().LRange(context.Background(), groupKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var messages []*models.GroupMessage
	for _, msgJSON := range msgs {
		var msg models.GroupMessage
		if err := json.Unmarshal([]byte(msgJSON), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}



func (r *messageRepository) GroupExists(groupKey string) (bool, error) {
	res, err := r.redisService.GetClient().Exists(context.Background(), groupKey).Result()
	return res == 1, err
}

func (r *messageRepository) AddMemberToGroup(groupKey string, member string) error {
	return r.redisService.GetClient().SAdd(context.Background(), groupKey, member).Err()
}

func (r *messageRepository) IsMemberOfGroup(groupKey string, member string) (bool, error) {
	return r.redisService.GetClient().SIsMember(context.Background(), groupKey, member).Result()
}

func (r *messageRepository) SendGroupMessage(groupKey string, msg *models.GroupMessage) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.redisService.GetClient().RPush(context.Background(), groupKey, msgJSON).Err()
}

func (r *messageRepository) SendBroadcastMessage(key string, msg *models.BroadcastMessage) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.redisService.GetClient().RPush(context.Background(), key, msgJSON).Err()
}

func (r *messageRepository) GetBroadcastHistory(key string) ([]*models.BroadcastMessage, error) {
	msgs, err := r.redisService.GetClient().LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var messages []*models.BroadcastMessage
	for _, msgJSON := range msgs {
		var msg models.BroadcastMessage
		if err := json.Unmarshal([]byte(msgJSON), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}