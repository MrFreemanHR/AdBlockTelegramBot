package bot

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/core/entity"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *bot) ProcessEvents(event *tgbotapi.Update) {
	if event.Message != nil {
		b.OnMessage(event)
	}
}

func (b *bot) OnMessage(event *tgbotapi.Update) {
	msg, err := b.processMessageEntity(event)
	if err == nil {
		logger.Logger().UselessInfo("Message: %s", msg.Text)
		for _, handler := range b.messageHandlers {
			processed := handler.ProcessMessage(&msg)
			if processed {
				break
			}
		}
	}
}

func (b *bot) processMessageEntity(message *tgbotapi.Update) (entity.TelegramMessage, error) {
	var messageEntity entity.TelegramMessage
	messageEntity.Text = message.Message.Text
	messageEntity.Audio = message.Message.Audio
	messageEntity.Photo = message.Message.Photo
	messageEntity.VideoNote = message.Message.VideoNote
	messageEntity.MessageID = int64(message.Message.MessageID)
	// Sender processing
	messageEntity.From = &entity.TelegramUser{
		ID:           message.Message.From.ID,
		IsBot:        message.Message.From.IsBot,
		FirstName:    message.Message.From.FirstName,
		LastName:     message.Message.From.LastName,
		UserName:     message.Message.From.UserName,
		LanguageCode: message.Message.From.LanguageCode,
	}
	// Chat processing
	messageEntity.Chat = &entity.TelegramChat{
		ID:    message.Message.Chat.ID,
		Type:  message.Message.Chat.Type,
		Title: message.Message.Chat.Title,
	}
	if message.Message.Chat.Permissions != nil {
		messageEntity.Chat.Permissions = &entity.TelegramChatPermissions{
			CanSendMessages:       message.Message.Chat.Permissions.CanSendMessages,
			CanSendMediaMessages:  message.Message.Chat.Permissions.CanSendMediaMessages,
			CanSendPolls:          message.Message.Chat.Permissions.CanSendPolls,
			CanSendOtherMessages:  message.Message.Chat.Permissions.CanSendOtherMessages,
			CanAddWebPagePreviews: message.Message.Chat.Permissions.CanAddWebPagePreviews,
			CanChangeInfo:         message.Message.Chat.Permissions.CanChangeInfo,
			CanInviteUsers:        message.Message.Chat.Permissions.CanInviteUsers,
			CanPinMessages:        message.Message.Chat.Permissions.CanPinMessages,
		}
	}
	return messageEntity, nil
}
