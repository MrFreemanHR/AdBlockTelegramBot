package transport

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/core/entity"
	"errors"
	"io"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zelenin/go-tdlib/client"
	"golang.org/x/net/html"
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
		if message.ParseMode != "" {
			msg.ParseMode = message.ParseMode
		}
		_, err := api.tgapi.Send(msg)
		return err
	}
	if api.tdlib != nil {
		var entities []*client.TextEntity
		var messageText = message.Text
		if message.ParseMode == "html" {
			messageText, entities = api.processHTMLEntities(message.Text)
		}
		content := client.InputMessageText{
			Text: &client.FormattedText{
				Text:     messageText,
				Entities: entities,
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
	if api.tdlib != nil {
		request := client.DeleteMessagesRequest{
			ChatId:     chatID,
			MessageIds: []int64{messageID},
			Revoke:     true,
		}
		_, err := api.tdlib.DeleteMessages(&request)
		return err
	}
	return nil
}

func (api *TelegramAPI) processHTMLEntities(text string) (string, []*client.TextEntity) {
	reader := strings.NewReader(text)
	tokenizer := html.NewTokenizer(reader)
	currentPos := 0
	newText := ""
	var bold, italic, strike, code, underline []*client.TextEntity
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				break
			}
			logger.Logger().Warn("[TDLIB Processing] Processing HTML tokens error: %s", tokenizer.Err().Error())
			break
		}
		token := tokenizer.Token()
		htmlToken := false
		if token.String() == "<b>" {
			api.setEntityParam(&currentPos, nil, &bold, &client.TextEntityTypeBold{})
			htmlToken = true
		}
		if token.String() == "</b>" {
			api.setEntityParam(nil, &currentPos, &bold, &client.TextEntityTypeBold{})
			htmlToken = true
		}
		if token.String() == "<i>" {
			api.setEntityParam(&currentPos, nil, &italic, &client.TextEntityTypeItalic{})
			htmlToken = true
		}
		if token.String() == "</i>" {
			api.setEntityParam(nil, &currentPos, &italic, &client.TextEntityTypeItalic{})
			htmlToken = true
		}
		if token.String() == "<s>" {
			api.setEntityParam(&currentPos, nil, &strike, &client.TextEntityTypeStrikethrough{})
			htmlToken = true
		}
		if token.String() == "</s>" {
			api.setEntityParam(nil, &currentPos, &strike, &client.TextEntityTypeStrikethrough{})
			htmlToken = true
		}
		if token.String() == "<code>" {
			api.setEntityParam(&currentPos, nil, &code, &client.TextEntityTypeCode{})
			htmlToken = true
		}
		if token.String() == "</code>" {
			api.setEntityParam(nil, &currentPos, &code, &client.TextEntityTypeCode{})
			htmlToken = true
		}
		if token.String() == "<u>" {
			api.setEntityParam(&currentPos, nil, &underline, &client.TextEntityTypeUnderline{})
			htmlToken = true
		}
		if token.String() == "</u>" {
			api.setEntityParam(nil, &currentPos, &underline, &client.TextEntityTypeUnderline{})
			htmlToken = true
		}
		if !htmlToken {
			currentPos += len(token.String())
			newText += token.Data
		}
	}
	var allEntities []*client.TextEntity
	allEntities = append(allEntities, bold...)
	allEntities = append(allEntities, italic...)
	allEntities = append(allEntities, strike...)
	allEntities = append(allEntities, code...)
	allEntities = append(allEntities, underline...)
	return newText, allEntities
}

func (api *TelegramAPI) setEntityParam(start *int, end *int, slice *[]*client.TextEntity, entityType client.TextEntityType) {
	var entity *client.TextEntity
	if len(*slice) == 0 {
		entity = &client.TextEntity{
			Offset: -1,
			Length: -1,
			Type:   entityType,
		}
	} else {
		entity = (*slice)[len(*slice)-1]
	}
	if start != nil {
		if entity.Offset == -1 {
			entity.Offset = int32(*start)
		}
	}
	if end != nil {
		if entity.Length == -1 {
			entity.Length = int32(*end) - entity.Offset
		}
	}
	if len(*slice) == 0 {
		*slice = append(*slice, entity)
	} else {
		(*slice)[len(*slice)-1] = entity
	}
}
