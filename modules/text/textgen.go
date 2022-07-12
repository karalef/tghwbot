package text

import (
	"log"
	"strings"
	"sync"
	"tghwbot/bot"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var textgenMut sync.Mutex

var Gen = bot.Command{
	Cmd:         "textgen",
	Description: "text generation",
	Run: func(ctx *bot.Context, msg *tgbotapi.Message, args []string) {
		query := strings.Join(args, " ")
		if query == "" {
			ctx.ReplyText("Придумайте начало истории")
			return
		}
		chat := ctx.Chat()
		replyID := chat.Send(bot.MessageConfig{
			ReplyToMessageID: msg.MessageID,
			MsgText: bot.MsgText{
				Text: "generating",
			},
		}).MessageID

		textgenMut.Lock()
		defer textgenMut.Unlock()

		ch := porfirevichAsync(query, 40)
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for i := 1; ; i++ {
			select {
			case result := <-ch:
				if result.Error != nil {
					log.Println("porfirevich error:", result.Error.Error())
					chat.Edit(replyID, result.Error.Error())
					return
				}
				var text string
				for _, r := range result.Replies {
					text += query + r + "\n\n"
				}
				chat.Edit(replyID, text)
				return
			case <-t.C:
				chat.Edit(replyID, "generating"+strings.Repeat(".", i))
				if i == 3 {
					i = 0
				}
			}
		}
	},
}
