package entity

// Entity that store needed information from Message to process internal logic
type TelegramMessage struct {
	// MessageID is a unique message identifier inside this chat
	MessageID int `json:"message_id"`
	// From is a sender, empty for messages sent to channels;
	//
	// optional
	From *TelegramUser `json:"from,omitempty"`
	// Chat is the conversation the message belongs to
	Chat *TelegramChat `json:"chat"`
}
