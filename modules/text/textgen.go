package text

import (
	"strings"
	"tghwbot/bot"
	"tghwbot/bot/tg"
)

var Gen = bot.Command{
	Cmd:         "textgen",
	Description: "text generation",
	Run: func(ctx *bot.Context, msg *tg.Message, args []string) {
		query := strings.Join(args, " ")
		if query == "" {
			ctx.Reply("Придумайте начало истории")
		}
		replies, err := porfirevich(query, 30)
		if err != nil {
			ctx.Logger().Error(err.Error())
			ctx.Reply(err.Error())
		}

		var text string
		for _, r := range replies {
			text += query + r + "\n\n"
		}
		ctx.Reply(text)
	},
}
