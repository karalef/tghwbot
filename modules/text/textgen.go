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
	Run: func(ctx bot.MessageContext, msg *tg.Message, args []string) error {
		query := strings.Join(args, " ")
		if query == "" {
			return ctx.ReplyText("Think of the beginning of the story")
		}

		textgenMut.Lock()
		defer textgenMut.Unlock()
		ctx.Chat.SendChatAction(tg.ActionTyping)
		replies, err := porfirevich(query, 30)
		if err != nil {
			ctx.Logger().Error(err.Error())
			return ctx.ReplyText(err.Error())
		}

		var text string
		for _, r := range replies {
			text += query + r + "\n\n"
		}
		return ctx.ReplyText(text)
	},
}
