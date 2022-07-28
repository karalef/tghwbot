package images

import (
	"encoding/json"
	"errors"
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
		q := strings.Join(args, " ")
		if q == "" {
			ctx.Reply("Provide keywords")
		}
		ctx.Chat.SendChatAction(tg.ActionUploadPhoto)
		result, err := searchImages(q)
		if err != nil {
			ctx.Logger().Error(err.Error())
			ctx.Reply(errors.Unwrap(err).Error())
		}
		if len(result) == 0 {
			ctx.Reply("No results")
		}
		u := result[random.RandP(len(result), 1.5)]
		p := bot.NewPhoto(tg.FileURL(u))
		p.Caption = q
		ctx.Chat.Send(p)
	},
}

func searchImages(q string) ([]string, error) {
	q = url.QueryEscape(q)
	resp, err := http.Get("https://imsea.herokuapp.com/api/1?q=" + q)
	if err != nil {
		return nil, errors.Unwrap(err)
	}
	var res struct {
		ImageName string   `json:"image_name"`
		Results   []string `json:"results"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	resp.Body.Close()
	return res.Results, nil
}
