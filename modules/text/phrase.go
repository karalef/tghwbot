package text

import (
	"errors"
	"io"
	"log"
	"net/http"
	"tghwbot/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Phrase = bot.Command{
	Cmd:         "phrase",
	Description: "generate cool tech-savvy sounding phrases",
	Run: func(ctx *bot.Context, _ *tgbotapi.Message, _ []string) {
		resp, err := http.Get("https://techy-api.vercel.app/api/text")
		if err != nil {
			log.Println("phrase:", err.Error())
			err = errors.Unwrap(err)
			ctx.ReplyText(err.Error())
		}
		defer resp.Body.Close()
		d, _ := io.ReadAll(resp.Body)
		ctx.ReplyText(string(d))
	},
}
