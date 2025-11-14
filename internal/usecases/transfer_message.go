package usecases

import (
	"log"
	"sync"

	"reaction/internal/entities"
)

// MessageReference - メッセージ参照情報
type MessageReference struct {
	GuildID   string
	ChannelID string
	MessageID string
}

// Message - メッセージ情報
type Message struct {
	GuildID   string
	ChannelID string
	ID        string
	Content   string
}

// DiscordClient - Discord APIとの通信インターフェース
type DiscordClient interface {
	GetMessage(channelID, messageID string) (*Message, error)
	SendMessageWithReference(channelID string, ref *MessageReference) (string, error)
	DeleteMessage(channelID, messageID string) error
}

// TransferMessageUseCase - メッセージ転送のビジネスロジック
type TransferMessageUseCase struct {
	config             *entities.Config
	transferMsgMapping map[string]string
	mappingMutex       sync.RWMutex
}

// NewTransferMessageUseCase - 新しいTransferMessageUseCaseを作成
func NewTransferMessageUseCase(config *entities.Config) *TransferMessageUseCase {
	return &TransferMessageUseCase{
		config:             config,
		transferMsgMapping: make(map[string]string),
	}
}

// TransferMessage - メッセージを転送し、マッピングを保存
func (uc *TransferMessageUseCase) TransferMessage(
	client DiscordClient,
	originalMsg *Message,
) error {
	transferRef := &MessageReference{
		GuildID:   originalMsg.GuildID,
		ChannelID: originalMsg.ChannelID,
		MessageID: originalMsg.ID,
	}

	transferredMsgID, err := client.SendMessageWithReference(uc.config.TransferChannelID, transferRef)
	if err != nil {
		log.Printf("メッセージ転送に失敗: %v", err)
		return err
	}

	uc.mappingMutex.Lock()
	uc.transferMsgMapping[originalMsg.ID] = transferredMsgID
	uc.mappingMutex.Unlock()

	log.Printf("メッセージ %s をチャンネル %s に転送しました (転送メッセージID: %s)", originalMsg.ID, uc.config.TransferChannelID, transferredMsgID)
	return nil
}

// DeleteTransferredMessage - 転送メッセージを削除
func (uc *TransferMessageUseCase) DeleteTransferredMessage(
	client DiscordClient,
	originalMsgID string,
) error {
	uc.mappingMutex.RLock()
	transferredMsgID, exists := uc.transferMsgMapping[originalMsgID]
	uc.mappingMutex.RUnlock()

	if !exists {
		log.Printf("メッセージ %s の転送記録が見つかりません", originalMsgID)
		return nil
	}

	err := client.DeleteMessage(uc.config.TransferChannelID, transferredMsgID)
	if err != nil {
		log.Printf("転送メッセージの削除に失敗: %v", err)
		return err
	}

	uc.mappingMutex.Lock()
	delete(uc.transferMsgMapping, originalMsgID)
	uc.mappingMutex.Unlock()

	log.Printf("転送メッセージ %s を削除しました", transferredMsgID)
	return nil
}

// IsTransferredMessage - メッセージが転送先かどうかを判定
func (uc *TransferMessageUseCase) IsTransferredMessage(msgID string) bool {
	uc.mappingMutex.RLock()
	defer uc.mappingMutex.RUnlock()

	for _, transferredMsgID := range uc.transferMsgMapping {
		if transferredMsgID == msgID {
			return true
		}
	}
	return false
}
