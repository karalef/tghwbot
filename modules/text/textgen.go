package text

import (
	"log"
	"strings"
	"tghwbot/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Gen = bot.Command{
	Cmd:         "textgen",
	Description: "text generation",
	Run: func(ctx *bot.Context, msg *tgbotapi.Message, args []string) {
		query := strings.Join(args, " ")
		if query == "" {
			ctx.ReplyText("Придумайте начало истории")
			return
		}
		replies, err := porfirevich(query, 30)
		if err != nil {
			log.Println("porfirevich error:", err.Error())
			ctx.ReplyText(err.Error())
		}

		var text string
		for _, r := range replies {
			text += query + r + "\n\n"
		}
		ctx.ReplyText(text)
	},
}
