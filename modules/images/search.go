package images

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"tghwbot/bot"
	"tghwbot/bot/tg"
	"tghwbot/modules/random"
)

var Search = bot.Command{
	Cmd:         "img",
	Description: "random image",
	Run: func(ctx *bot.Context, _ *tg.Message, args []string) {
		q := url.QueryEscape(strings.Join(args, " "))
		if q == "" {
			ctx.Reply("Provide keywords")
		}
		ctx.Chat.SendChatAction(tg.ActionUploadPhoto)
		resp, err := http.Get("https://imsea.herokuapp.com/api/1?q=" + q)
		if err != nil {
			log.Println("image search:", err.Error())
			err = errors.Unwrap(err)
			ctx.Reply(err.Error())
		}
		defer resp.Body.Close()
		var res struct {
			ImageName string   `json:"image_name"`
			Results   []string `json:"results"`
		}
		json.NewDecoder(resp.Body).Decode(&res)
		if len(res.Results) == 0 {
			ctx.Reply("No results")
		}
		u := res.Results[random.RandP(len(res.Results), 1.5)]
		p := bot.NewPhoto(tg.FileURL(u))
		p.Caption = res.ImageName
		ctx.Chat.Send(p)
	},
}
