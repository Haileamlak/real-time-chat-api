package usecases

import (
	
	"fmt"
	"sort"
	"context"

	"github.com/haileamlak/chat-system/models"
	"github.com/haileamlak/chat-system/repositories"
)
type MessageUseCase interface {
	SaveDirectMessage(ctx context.Context, user1, user2 string, msg *models.DirectMessage) error
	GetDirectMessage(ctx context.Context, user1, user2 string, index int64) (*models.DirectMessage, error)
	GetDMHistory(ctx context.Context, user1, user2 string) ([]*models.DirectMessage, error)
	CreateGroup(ctx context.Context, groupName string, members []string) error
	GetGroupHistory(ctx context.Context, groupName string) ([]*models.GroupMessage, error)
	GroupExists(ctx context.Context, groupName string) (bool, error)
	AddMemberToGroup(ctx context.Context, groupName, member string) error
	IsMemberOfGroup(ctx context.Context, groupName, member string) (bool, error)
	SendGroupMessage(ctx context.Context, groupName string, msg *models.GroupMessage) error
	SendBroadcastMessage(ctx context.Context, msg *models.BroadcastMessage) error
	GetBroadcastHistory(ctx context.Context) ([]*models.BroadcastMessage, error)
}

type messageUseCase struct {
	messageRepo repositories.MessageRepository
}

func NewMessageUseCase(messageRepo repositories.MessageRepository) MessageUseCase {
	return &messageUseCase{
		messageRepo: messageRepo,
	}
}

func (m *messageUseCase) SaveDirectMessage(ctx context.Context, user1, user2 string, msg *models.DirectMessage) error {
	key := getDMKey(user1, user2)
	return m.messageRepo.SaveDirectMessage(key, msg)
}

func (m *messageUseCase) GetDirectMessage(ctx context.Context, user1, user2 string, index int64) (*models.DirectMessage, error) {
	key := getDMKey(user1, user2)
	return m.messageRepo.GetDirectMessage(key, index)
}

func (m *messageUseCase) GetDMHistory(ctx context.Context, user1, user2 string) ([]*models.DirectMessage, error) {
	key := getDMKey(user1, user2)
	return m.messageRepo.GetDMHistory(key)
}

func (m *messageUseCase) CreateGroup(ctx context.Context, groupName string, members []string) error {
	groupKey := fmt.Sprintf("group:%s:members", groupName)
	return m.messageRepo.CreateGroup(groupKey, members)
}
func (m *messageUseCase) GetGroupHistory(ctx context.Context, groupName string) ([]*models.GroupMessage, error) {
	groupKey := fmt.Sprintf("group:%s:messages", groupName)
	return m.messageRepo.GetGroupHistory(groupKey)
}
func (m *messageUseCase) GroupExists(ctx context.Context, groupName string) (bool, error) {
	groupKey := fmt.Sprintf("group:%s:members", groupName)
	return m.messageRepo.GroupExists(groupKey)
}
func (m *messageUseCase) AddMemberToGroup(ctx context.Context, groupName, member string) error {
	groupKey := fmt.Sprintf("group:%s:members", groupName)
	return m.messageRepo.AddMemberToGroup(groupKey, member)
}
func (m *messageUseCase) IsMemberOfGroup(ctx context.Context, groupName, member string) (bool, error) {
	groupKey := fmt.Sprintf("group:%s:members", groupName)
	return m.messageRepo.IsMemberOfGroup(groupKey, member)
}
func (m *messageUseCase) SendGroupMessage(ctx context.Context, groupName string, msg *models.GroupMessage) error {
	groupKey := fmt.Sprintf("group:%s:messages", groupName)
	return m.messageRepo.SendGroupMessage(groupKey, msg)
}
func (m *messageUseCase) SendBroadcastMessage(ctx context.Context, msg *models.BroadcastMessage) error {
	broadcastKey := "broadcast:messages"
	return m.messageRepo.SendBroadcastMessage(broadcastKey, msg)
}
func (m *messageUseCase) GetBroadcastHistory(ctx context.Context) ([]*models.BroadcastMessage, error) {	
	broadcastKey := "broadcast:messages"
	return m.messageRepo.GetBroadcastHistory(broadcastKey)
}


func getDMKey(user1, user2 string) string {
	users := []string{user1, user2}
	sort.Strings(users)
	return fmt.Sprintf("dm:%s:%s", users[0], users[1])
}