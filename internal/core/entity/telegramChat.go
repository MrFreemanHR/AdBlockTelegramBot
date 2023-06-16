package entity

// Entity that store needed information from Chat to process internal logic
type TelegramChat struct {
	// ID is a unique identifier for this chat
	ID int64 `json:"id"`
	// Type of chat, can be either “private”, “group”, “supergroup” or “channel”
	Type string `json:"type"`
	// Title for supergroups, channels and group chats
	//
	// optional
	Title string `json:"title,omitempty"`
	// Permissions are default chat member permissions, for groups and
	// supergroups. Returned only in getChat.
	//
	// optional
	Permissions *TelegramChatPermissions `json:"permissions,omitempty"`
}

func (c TelegramChat) IsSuperGroup() bool {
	return c.Type == "supergroup" || c.Type == "chatTypeSupergroup"
}

func (c TelegramChat) IsPrivate() bool {
	return c.Type == "private" || c.Type == "chatTypePrivate"
}
