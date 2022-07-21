package text

import (
	"strings"
	"sync"
	"tghwbot/bot"
	"tghwbot/bot/tg"
)

var textgenMut sync.Mutex

var Gen = bot.Command{
	Cmd:         "textgen",
	Description: "text generation",
	Run: func(ctx *bot.Context, msg *tg.Message, args []string) {
		query := strings.Join(args, " ")
		if query == "" {
			ctx.Reply("Think of the beginning of the story")
		}

		textgenMut.Lock()
		defer textgenMut.Unlock()
		ctx.Chat.Send(bot.ChatAction(tg.ActionTyping))
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
