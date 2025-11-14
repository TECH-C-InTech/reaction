package handlers

import (
	"log"

	"reaction/internal/entities"
	"reaction/internal/usecases"

	"github.com/bwmarrin/discordgo"
)

// DiscordHandler - Discordイベントを処理するハンドラー
type DiscordHandler struct {
	transferUseCase *usecases.TransferMessageUseCase
	gateway         usecases.DiscordClient
	config          *entities.Config
}

// NewDiscordHandler - 新しいDiscordHandlerを作成
func NewDiscordHandler(
	transferUseCase *usecases.TransferMessageUseCase,
	gateway usecases.DiscordClient,
	config *entities.Config,
) *DiscordHandler {
	return &DiscordHandler{
		transferUseCase: transferUseCase,
		gateway:         gateway,
		config:          config,
	}
}

// HandleReactionAdd - リアクション追加イベントを処理
func (h *DiscordHandler) HandleReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	if !h.isTriggerReactionEmoji(r.Emoji) {
		return
	}

	if h.transferUseCase.IsTransferredMessage(r.MessageID) {
		log.Printf("メッセージ %s は転送メッセージのため、リアクション追加を無視します", r.MessageID)
		return
	}

	originalMsg, err := h.gateway.GetMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Printf("メッセージ取得に失敗: %v", err)
		return
	}

	discordMsg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Printf("リアクション確認のためのメッセージ取得に失敗: %v", err)
		return
	}

	reactionCount := h.getTriggerReactionCount(discordMsg)
	if reactionCount > 1 {
		log.Printf("メッセージ %s には既にトリガーリアクションが %d 個ついているため、転送をスキップします", r.MessageID, reactionCount)
		return
	}

	err = h.transferUseCase.TransferMessage(h.gateway, originalMsg)
	if err != nil {
		log.Printf("メッセージ転送に失敗: %v", err)
		return
	}
}

// HandleReactionRemove - リアクション削除イベントを処理
func (h *DiscordHandler) HandleReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.UserID == s.State.User.ID {
		return
	}

	if !h.isTriggerReactionEmoji(r.Emoji) {
		return
	}

	if h.transferUseCase.IsTransferredMessage(r.MessageID) {
		log.Printf("メッセージ %s は転送メッセージのため、リアクション削除を無視します", r.MessageID)
		return
	}

	discordMsg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Printf("メッセージ取得に失敗: %v", err)
		return
	}

	triggerReactionCount := h.getTriggerReactionCount(discordMsg)
	if triggerReactionCount > 0 {
		log.Printf("メッセージ %s にはまだトリガーリアクションが %d 個残っているため、削除をスキップします", r.MessageID, triggerReactionCount)
		return
	}

	err = h.transferUseCase.DeleteTransferredMessage(h.gateway, r.MessageID)
	if err != nil {
		log.Printf("転送メッセージの削除に失敗: %v", err)
		return
	}
}

func (h *DiscordHandler) isTriggerReactionEmoji(emoji discordgo.Emoji) bool {
	return emoji.ID == h.config.TriggerReactionEmoji
}

func (h *DiscordHandler) getTriggerReactionCount(msg *discordgo.Message) int {
	for _, reaction := range msg.Reactions {
		if reaction.Emoji.ID == h.config.TriggerReactionEmoji {
			return reaction.Count
		}
	}
	return 0
}
