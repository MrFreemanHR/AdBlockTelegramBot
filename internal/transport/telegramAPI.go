package transport

import (
	"adblock_bot/internal/core/entity"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zelenin/go-tdlib/client"
)

type TelegramAPI struct {
	tgapi *tgbotapi.BotAPI
	tdlib *client.Client
	self  entity.TelegramUser
}

func NewUnifiedTransportWithBotAPI(tgapi *tgbotapi.BotAPI, self entity.TelegramUser) *TelegramAPI {
	return &TelegramAPI{
		tgapi: tgapi,
		self:  self,
	}
}

func NewUnifiedTransportWithTDlib(tdlib *client.Client, self entity.TelegramUser) *TelegramAPI {
	return &TelegramAPI{
		tdlib: tdlib,
		self:  self,
	}
}

func (api *TelegramAPI) Self() entity.TelegramUser {
	return api.self
}

func (api *TelegramAPI) SendMessage(message entity.TelegramMessage) error {
	if api.tdlib == nil && api.tgapi == nil {
		return errors.New("no transport available")
	}
	if api.tgapi != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
		if message.ReplyToMessage != nil {
			msg.ReplyToMessageID = int(message.ReplyToMessage.MessageID)
		}
		_, err := api.tgapi.Send(msg)
		return err
	}
	if api.tdlib != nil {
		content := client.InputMessageText{
			Text: &client.FormattedText{
				Text: message.Text,
			},
		}
		request := client.SendMessageRequest{
			ChatId:              message.Chat.ID,
			InputMessageContent: &content,
		}
		if message.ReplyToMessage != nil {
			request.ReplyToMessageId = message.ReplyToMessage.MessageID
		}
		_, err := api.tdlib.SendMessage(&request)
		return err
	}
	return nil
}

func (api *TelegramAPI) RemoveMessage(chatID, messageID int64) error {
	if api.tdlib == nil && api.tgapi == nil {
		return errors.New("no transport available")
	}
	if api.tgapi != nil {
		msg := tgbotapi.NewDeleteMessage(chatID, int(messageID))
		_, err := api.tgapi.Send(msg)
		return err
	}
	return nil
}
