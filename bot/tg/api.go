package tg

import (
	"encoding/json"
	"fmt"
)

// DefaultAPIURL is a default url for telegram api.
const DefaultAPIURL = "https://api.telegram.org"

// APIResponse represents telegram api response.
type APIResponse struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
	*APIError
}

// APIError describes telegram api error.
type APIError struct {
	Code        int    `json:"error_code"`
	Description string `json:"description"`
	Parameters  *struct {
		MigrateTo  *int64 `json:"migrate_to_chat_id"`
		RetryAfter *int   `json:"retry_after"`
	} `json:"parameters"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("telegram: (%d) %s", e.Code, e.Description)
}

// Update object represents an incoming update.
type Update struct {
	ID                int      `json:"update_id"`
	Message           *Message `json:"message"`
	EditedMessage     *Message `json:"edited_message"`
	ChannelPost       *Message `json:"channel_post"`
	EditedChannelPost *Message `json:"edited_channel_post"`
	//CallbackQuery     *CallbackQuery `json:"callback_query"`
	InlineQuery  *InlineQuery  `json:"inline_query"`
	InlineResult *InlineResult `json:"chosen_inline_result"`
	//ShippingQuery     *ShippingQuery `json:"shipping_query"`
	//PreCheckoutQuery  *PreCheckoutQuery `json:"pre_checkout_query"`
	Poll       *Poll       `json:"poll"`
	PollAnswer *PollAnswer `json:"poll_answer"`
	//MyChatMember    *ChatMemberUpdate `json:"my_chat_member"`
	//ChatMember      *ChatMemberUpdate `json:"chat_member"`
	//ChatJoinRequest *ChatJoinRequest  `json:"chat_join_request"`
}

// InlineQuery is an incoming inline query. When the user sends
// an empty query, your bot could return some default or
// trending results.
type InlineQuery struct {
	ID       string    `json:"id"`
	From     *User     `json:"from"`
	Query    string    `json:"query"` // up to 256 characters
	Offset   string    `json:"offset"`
	ChatType string    `json:"chat_type"`
	Location *Location `json:"location"`
}

// InlineResult represents a result of an inline query that was chosen
// by the user and sent to their chat partner.
type InlineResult struct {
	ResultID        string    `json:"result_id"`
	From            *User     `json:"from"`
	Location        *Location `json:"location"`
	InlineMessageID string    `json:"inline_message_id"`
	Query           string    `json:"query"`
}

// Command represents a bot command.
type Command struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}
