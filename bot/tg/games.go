package tg

// CallbackGame is a placeholder, currently holds no information.
// Use BotFather to set up your game.
type CallbackGame struct{}

// Game represents a game.
type Game struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Photo       []PhotoSize     `json:"photo"`
	Text        string          `json:"text,omitempty"`
	Entities    []MessageEntity `json:"text_entities,omitempty"`
	Animation   *Animation      `json:"animation,omitempty"`
}

// GameHighScore represents one row of the high scores table for a game.
type GameHighScore struct {
	Position int   `json:"position"`
	User     *User `json:"user"`
	Score    int   `json:"score"`
}
