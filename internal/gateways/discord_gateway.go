package gateways

import (
	"reaction/internal/usecases"

	"github.com/bwmarrin/discordgo"
)

// DiscordGateway - Discord APIとの通信ゲートウェイ
type DiscordGateway struct {
	session *discordgo.Session
}

// NewDiscordGateway - 新しいDiscordGatewayを作成
func NewDiscordGateway(session *discordgo.Session) *DiscordGateway {
	return &DiscordGateway{
		session: session,
	}
}

// GetMessage - メッセージを取得
func (g *DiscordGateway) GetMessage(channelID, messageID string) (*usecases.Message, error) {
	msg, err := g.session.ChannelMessage(channelID, messageID)
	if err != nil {
		return nil, err
	}

	return &usecases.Message{
		GuildID:   msg.GuildID,
		ChannelID: msg.ChannelID,
		ID:        msg.ID,
		Content:   msg.Content,
	}, nil
}

// SendMessageWithReference - メッセージ参照を使って転送
func (g *DiscordGateway) SendMessageWithReference(channelID string, ref *usecases.MessageReference) (string, error) {
	discordRef := &discordgo.MessageReference{
		GuildID:   ref.GuildID,
		ChannelID: ref.ChannelID,
		MessageID: ref.MessageID,
	}

	transferMsgSend := &discordgo.MessageSend{
		Reference: discordRef,
	}

	transferredMsg, err := g.session.ChannelMessageSendComplex(channelID, transferMsgSend)
	if err != nil {
		return "", err
	}

	return transferredMsg.ID, nil
}

// DeleteMessage - メッセージを削除
func (g *DiscordGateway) DeleteMessage(channelID, messageID string) error {
	return g.session.ChannelMessageDelete(channelID, messageID)
}
