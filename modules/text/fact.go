package text

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"tghwbot/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Fact = bot.Command{
	Cmd:         "fact",
	Description: "random useless (but true) fact",
	Run: func(ctx *bot.Context, _ *tgbotapi.Message, args []string) {
		t := "random.json"
		if len(args) > 0 && args[0] == "today" {
			t = "today.json"
		}
		resp, err := http.Get("https://uselessfacts.jsph.pl/" + t + "?" + url.Values{
			"language": {"en"},
		}.Encode())
		if err != nil {
			log.Println("useless fact:", err.Error())
			err = errors.Unwrap(err)
			ctx.ReplyText(err.Error())
		}
		defer resp.Body.Close()
		var res struct {
			Text string `json:"text"`
		}
		json.NewDecoder(resp.Body).Decode(&res)
		ctx.ReplyText(res.Text)
	},
}
