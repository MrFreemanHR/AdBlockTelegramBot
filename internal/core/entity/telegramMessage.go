package entity

// Entity that store needed information from Message to process internal logic
type TelegramMessage struct {
	// MessageID is a unique message identifier inside this chat
	MessageID int64 `json:"message_id"`
	// From is a sender, empty for messages sent to channels;
	//
	// optional
	From *TelegramUser `json:"from,omitempty"`
	// Chat is the conversation the message belongs to
	Chat *TelegramChat `json:"chat"`
	// Text is for text messages, the actual UTF-8 text of the message, 0-4096 characters;
	//
	// optional
	Text string `json:"text,omitempty"`
	// Audio message
	Audio interface{} `json:"audio,omitempty"`
	// Photo message
	Photo interface{} `json:"photo,omitempty"`
	// VideoNote message
	VideoNote interface{} `json:"video_note,omitempty"`
	// ReplyToMessage for replies, the original message.
	// Note that the Message object in this field will not contain further ReplyToMessage fields
	// even if it itself is a reply;
	//
	// optional
	ReplyToMessage *TelegramMessage `json:"reply_to_message,omitempty"`
}

func NewMessage(chatID int64, text string) TelegramMessage {
	return TelegramMessage{
		Chat: &TelegramChat{ID: chatID},
		Text: text,
	}
}
