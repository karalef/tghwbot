package tg

// Contact represents a phone contact.
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserID      int64  `json:"user_id"`
	Vcard       string `json:"vcard"`
}

// DiceEmoji represenst dice emoji.
type DiceEmoji string

// all available animated emojis.
const (
	DiceCube DiceEmoji = "🎲"
	DiceDart DiceEmoji = "🎯"
	DiceBall DiceEmoji = "🏀"
	DiceGoal DiceEmoji = "⚽"
	DiceSlot DiceEmoji = "🎰"
	DiceBowl DiceEmoji = "🎳"
)

// Dice represents an animated emoji that displays a random value.
type Dice struct {
	Emoji DiceEmoji `json:"emoji"`
	Value int       `json:"value"`
}

// PollOption contains information about one answer option in a poll.
type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}

// PollAnswer represents an answer of a user in a non-anonymous poll.
type PollAnswer struct {
	PollID  string `json:"poll_id"`
	User    *User  `json:"user"`
	Options []int  `json:"option_ids"`
}

// PollType represents poll type.
type PollType string

// all available poll types.
const (
	PollQuiz    PollType = "quiz"
	PollRegular PollType = "regular"
)

// Poll contains information about a poll.
type Poll struct {
	ID                  string          `json:"id"`
	Question            string          `json:"question"`
	Options             []PollOption    `json:"options"`
	VoterCount          int             `json:"total_voter_count"`
	IsClosed            bool            `json:"is_closed"`
	IsAnonymous         bool            `json:"is_anonymous"`
	Type                PollType        `json:"type"`
	MultipleAnswers     bool            `json:"allows_multiple_answers"`
	CorrectOption       int             `json:"correct_option_id"`
	Explanation         string          `json:"explanation"`
	ExplanationEntities []MessageEntity `json:"explanation_entities"`
	OpenPeriod          int             `json:"open_period"`
	CloseDate           int64           `json:"close_date"`
}

// Location represents a point on the map.
type Location struct {
	Long               float32  `json:"longitude"`
	Lat                float32  `json:"latitude"`
	HorizontalAccuracy *float32 `json:"horizontal_accuracy,omitempty"`
	LivePeriod         int      `json:"live_period,omitempty"`
	Heading            int      `json:"heading,omitempty"`
	AlertRadius        int      `json:"proximity_alert_radius,omitempty"`
}

// Venue represents a venue.
type Venue struct {
	Location        Location `json:"location"`
	Title           string   `json:"title"`
	Address         string   `json:"address"`
	FoursquareID    string   `json:"foursquare_id"`
	FoursquareType  string   `json:"foursquare_type"`
	GooglePlaceID   string   `json:"google_place_id"`
	GooglePlaceType string   `json:"google_place_type"`
}

// StickerSet represents a sticker set.
type StickerSet struct {
	Name          string     `json:"name"`
	Title         string     `json:"title"`
	IsAnimated    bool       `json:"is_animated"`
	IsVideo       bool       `json:"is_video"`
	ContainsMasks bool       `json:"contains_masks"`
	Stickers      []Sticker  `json:"stickers"`
	Thumbnail     *PhotoSize `json:"thumb"`
}

// MaskPosition describes the position on faces where a mask should be placed by default.
type MaskPosition struct {
	Point  string  `json:"point"`
	XShift float64 `json:"x_shift"`
	YShift float64 `json:"y_shift"`
	Scale  float64 `json:"scale"`
}
