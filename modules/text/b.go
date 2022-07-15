package text

import (
	"encoding/json"
	"log"
	"net/http"
	"tghwbot/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Buzzword = bot.Command{
	Cmd:         "b",
	Description: "bullshit generator",
	Run: func(ctx *bot.Context, _ *tgbotapi.Message, _ []string) {
		resp, err := http.Get("https://corporatebs-generator.sameerkumar.website/")
		if err != nil {
			log.Println("bullshit generator:", err.Error())
			ctx.ReplyText(err.Error())
		}
		defer resp.Body.Close()
		var res struct {
			Phrase string `json:"phrase"`
		}
		json.NewDecoder(resp.Body).Decode(&res)
		ctx.ReplyText(res.Phrase)
	},
}
