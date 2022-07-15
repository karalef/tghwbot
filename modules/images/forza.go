package images

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"tghwbot/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Forza = bot.Command{
	Cmd:         "forza",
	Description: "random forza car photo",
	Run: func(ctx *bot.Context, _ *tgbotapi.Message, args []string) {
		resp, err := http.Get("https://forza-api.tk/")
		if err != nil {
			log.Println("forza photo:", err.Error())
			err = errors.Unwrap(err)
			ctx.ReplyText(err.Error())
		}
		defer resp.Body.Close()
		var res struct {
			Image string `json:"image"`
		}
		json.NewDecoder(resp.Body).Decode(&res)
		ls := strings.LastIndexByte(res.Image, '/') + 1
		ld := strings.LastIndexByte(res.Image, '.')
		car := strings.ReplaceAll(res.Image[ls:ld], "_", " ")
		ctx.Chat().SendPhoto(bot.NewPhoto(car, tgbotapi.FileURL(res.Image)))
	},
}
