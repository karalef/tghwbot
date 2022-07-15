package images

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"tghwbot/bot"
	"tghwbot/modules/random"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Search = bot.Command{
	Cmd:         "img",
	Description: "random image",
	Run: func(ctx *bot.Context, _ *tgbotapi.Message, args []string) {
		q := url.QueryEscape(strings.Join(args, " "))
		resp, err := http.Get("https://imsea.herokuapp.com/api/1?q=" + q)
		if err != nil {
			log.Println("bullshit generator:", err.Error())
			ctx.ReplyText(err.Error())
		}
		defer resp.Body.Close()
		var res struct {
			ImageName string   `json:"image_name"`
			Results   []string `json:"results"`
		}
		json.NewDecoder(resp.Body).Decode(&res)
		u := res.Results[random.Rand(len(res.Results))]
		ctx.Chat().SendPhoto(bot.NewPhoto(res.ImageName, tgbotapi.FileURL(u)))
	},
}
