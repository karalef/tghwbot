package tg

import (
	"time"
)

// Message represents a message.
type Message struct {
	ID                  int             `json:"message_id"`
	From                *User           `json:"from"`
	SenderChat          *Chat           `json:"sender_chat"`
	Date                int64           `json:"date"`
	Chat                *Chat           `json:"chat"`
	FrowardFrom         *User           `json:"forward_from"`
	ForwardChat         *Chat           `json:"forward_from_chat"`
	ForwardMessageID    int             `json:"forward_from_message_id"`
	ForwardSignature    string          `json:"forward_signature"`
	ForwardSenderName   string          `json:"forward_sender_name"`
	ForwardDate         int             `json:"forward_date"`
	IsAutomaticForward  bool            `json:"is_automatic_forward"`
	ReplyTo             *Message        `json:"reply_to_message"`
	ViaBot              *User           `json:"via_bot"`
	EditDate            int64           `json:"edit_date"`
	HasProtectedContent bool            `json:"has_protected_content"`
	MediaGroupID        string          `json:"media_group_id"`
	AuthorSignature     string          `json:"author_signature"`
	Text                string          `json:"text"`
	Entities            []MessageEntity `json:"entities"`
	Animation           *Animation      `json:"animation"`
	Audio               *Audio          `json:"audio"`
	Document            *Document       `json:"document"`
	Photo               []PhotoSize     `json:"photo"`
	Sticker             *Sticker        `json:"sticker"`
	Video               *Video          `json:"video"`
	VideoNote           *VideoNote      `json:"video_note"`
	Voice               *Voice          `json:"voice"`
	Caption             string          `json:"caption"`
	CaptionEntities     []MessageEntity `json:"caption_entities"`
	Contact             *Contact        `json:"contact"`
	Dice                *Dice           `json:"dice"`
	// Game *Game `json:"game"`
	Poll              *Poll       `json:"poll"`
	Venue             *Venue      `json:"venue"`
	Location          *Location   `json:"location"`
	NewChatMembers    []*User     `json:"new_chat_members"`
	LeftChatMember    *User       `json:"left_chat_member"`
	NewChatTitle      string      `json:"new_chat_title"`
	NewChatPhoto      []PhotoSize `json:"new_chat_photo"`
	DeleteChatPhoto   bool        `json:"delete_chat_photo"`
	GroupCreated      bool        `json:"group_chat_created"`
	SuperGroupCreated bool        `json:"supergroup_chat_created"`
	ChannelCreated    bool        `json:"channel_chat_created"`
	MigrateTo         int64       `json:"migrate_to_chat_id"`
	MigrateFrom       int64       `json:"migrate_from_chat_id"`
	PinnedMessage     *Message    `json:"pinned_message"`
	// Invoice *Invoice `json:"invoice"`
	// SuccessfulPayment *SuccessfulPayment `json:"successful_payment"`
	// ConnectedWebsite string `json:"connected_website"`
	// PassportData *PassportData `json:"passport_data"`
	ProximityAlert *ProximityAlert `json:"proximity_alert_triggered"`
	// VideoChatScheduled *VideoChatScheduled `json:"video_chat_scheduled"`
	// VideoChatStarted *VideoChatStarted `json:"video_chat_started"``
	// VideoChatEnded *VideoChatEnded `json:"video_chat_ended"`
	// VideoChatInvited *VideoChatInvited `json:"video_chat_participants_invited"``
	WebAppData  *WebAppData           `json:"web_app_data"`
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
}

// Time converts unixtime to time.Time.
func (m *Message) Time() time.Time {
	return time.Unix(m.Date, 0)
}

// MessageEntity represents one special entitty in a text message.
// For example, hashtag, usernames, URLs, etc.
type MessageEntity struct {
	Type     EntityType `json:"type"`
	Offset   int        `json:"offset"` // in UTF-16
	Length   int        `json:"length"`
	URL      string     `json:"url"`
	User     *User      `json:"user"`
	Language string     `json:"language"`
}

// EntityType is a MessageEntity type.
type EntityType string

// all available entity types.
const (
	EntityMention       EntityType = "mention"
	EntityTMention      EntityType = "text_mention"
	EntityHashtag       EntityType = "hashtag"
	EntityCashtag       EntityType = "cashtag"
	EntityCommand       EntityType = "bot_command"
	EntityURL           EntityType = "url"
	EntityEmail         EntityType = "email"
	EntityPhone         EntityType = "phone_number"
	EntityBold          EntityType = "bold"
	EntityItalic        EntityType = "italic"
	EntityUnderline     EntityType = "underline"
	EntityStrikethrough EntityType = "strikethrough"
	EntityCode          EntityType = "code"
	EntityCodeBlock     EntityType = "pre"
	EntityTextLink      EntityType = "text_link"
	EntitySpoiler       EntityType = "spoiler"
)

// ProximityAlert represents the content of a service message,
// sent whenever a user in the chat triggers a proximity alert
// set by another user.
type ProximityAlert struct {
	Traveler *User `json:"traveler"`
	Watcher  *User `json:"watcher"`
	Distance int   `json:"distance"`
}

// ParseMode type.
type ParseMode string

// all available parse modes.
const (
	Markdown   ParseMode = "Markdown"
	MarkdownV2 ParseMode = "MarkdownV2"
	HTML       ParseMode = "HTML"
)
