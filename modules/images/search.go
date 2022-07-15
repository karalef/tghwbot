package images

import (
	"encoding/json"
	"errors"
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
		if q == "" {
			ctx.ReplyText("Provide keywords")
		}
		resp, err := http.Get("https://imsea.herokuapp.com/api/1?q=" + q)
		if err != nil {
			log.Println("image search:", err.Error())
			err = errors.Unwrap(err)
			ctx.ReplyText(err.Error())
		}
		defer resp.Body.Close()
		var res struct {
			ImageName string   `json:"image_name"`
			Results   []string `json:"results"`
		}
		json.NewDecoder(resp.Body).Decode(&res)
		if len(res.Results) == 0 {
			ctx.ReplyText("No results")
		}
		u := res.Results[random.RandP(len(res.Results), 1.4)]
		ctx.Chat().SendPhoto(bot.NewPhoto(res.ImageName, tgbotapi.FileURL(u)))
	},
}
