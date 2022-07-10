package text

import (
	"strings"
	"tghwbot/bot"

	"gopkg.in/telebot.v3"
)

var Gen = bot.Command{
	Cmd:         "textgen",
	Description: "text generation",
	Run: func(ctx *bot.Context, msg *telebot.Message, args []string) {
		query := strings.Join(args, " ")
		if query == "" {
			ctx.ReplyClose("Придумайте начало истории")
		}
		replies, err := porfirevich(query, 30)
		if err != nil {
			ctx.ReplyError(err)
		}

		var text string
		for _, r := range replies {
			text += query + r + "\n\n"
		}
		ctx.Reply(text)
	},
}
