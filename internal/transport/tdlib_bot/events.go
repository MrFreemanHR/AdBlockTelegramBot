package tdlibbot

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/core/entity"
	"errors"

	"github.com/zelenin/go-tdlib/client"
)

func (b *tdlibbot) ProcessEvents(event client.Type) {
	if event.GetType() == client.TypeUpdateNewMessage {
		message, ok := (event).(*client.UpdateNewMessage)
		if !ok {
			logger.Logger().Warn("[TDLIB Processing] Can't cast new message to struct")
			return
		}
		b.OnMessage(message)
	}
}

func (b *tdlibbot) OnMessage(event *client.UpdateNewMessage) {
	msg, err := b.processMessageEntity(event.Message)
	if err == nil {
		if *msg.From == b.GetTelegramAPI().Self() {
			return
		}
		logger.Logger().UselessInfo("Message: %s", msg.Text)
		for _, handler := range b.messageHandlers {
			processed := handler.ProcessMessage(&msg)
			if processed {
				break
			}
		}
	}
}

func (b *tdlibbot) processMessageEntity(message *client.Message) (entity.TelegramMessage, error) {
	var messageEntity entity.TelegramMessage
	switch message.Content.MessageContentType() {
	case client.TypeMessageText:
		messageEntity = b.processTextMessage((message.Content).(*client.MessageText))
	}
	messageEntity.MessageID = message.Id
	// Sender processing
	sender, ok := (message.SenderId).(*client.MessageSenderUser)
	if ok {
		user, err := b.client.GetUser(&client.GetUserRequest{
			UserId: sender.UserId,
		})
		if err != nil {
			logger.Logger().Warn("[TDLIB Proccessing] Can't get user by ID: %s", err.Error())
			return entity.TelegramMessage{}, err
		}
		messageEntity.From = &entity.TelegramUser{
			ID:           user.Id,
			IsBot:        user.Type.UserTypeType() == client.TypeUserTypeBot,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			UserName:     user.Username,
			LanguageCode: user.LanguageCode,
		}
	} else {
		logger.Logger().Warn("[TDLIB Processing] Can't cast sender to MessageSenderUser")
		return entity.TelegramMessage{}, errors.New("cast error: sender to messagesender")
	}
	// Chat processing
	chat, err := b.client.GetChat(&client.GetChatRequest{
		ChatId: message.ChatId,
	})
	if err != nil {
		logger.Logger().Warn("[TDLIB Processing] Can't get chat by ID: %s", err.Error())
		return entity.TelegramMessage{}, err
	}
	messageEntity.Chat = &entity.TelegramChat{
		ID:    chat.Id,
		Type:  chat.Type.ChatTypeType(),
		Title: chat.Title,
		Permissions: &entity.TelegramChatPermissions{
			CanSendMessages:       chat.Permissions.CanSendMessages,
			CanSendMediaMessages:  chat.Permissions.CanSendMediaMessages,
			CanSendPolls:          chat.Permissions.CanSendPolls,
			CanSendOtherMessages:  chat.Permissions.CanSendOtherMessages,
			CanAddWebPagePreviews: chat.Permissions.CanAddWebPagePreviews,
			CanChangeInfo:         chat.Permissions.CanChangeInfo,
			CanInviteUsers:        chat.Permissions.CanInviteUsers,
			CanPinMessages:        chat.Permissions.CanPinMessages,
		},
	}
	return messageEntity, nil
}

func (b *tdlibbot) processTextMessage(message *client.MessageText) entity.TelegramMessage {
	var messageEntity = entity.TelegramMessage{
		Text: message.Text.Text,
	}
	return messageEntity
}
