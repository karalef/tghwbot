package tg

// WebAppInfo describes a Web App.
type WebAppInfo struct {
	URL string `json:"url"`
}

// WebAppData describes data sent from a Web App to the bot.
type WebAppData struct {
	Data       string `json:"data"`
	ButtonText string `json:"button_text"`
}
