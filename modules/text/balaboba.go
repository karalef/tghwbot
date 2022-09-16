package text

import (
	"strings"
	"tghwbot/bot"
	"tghwbot/bot/tg"
	"tghwbot/modules/text/bb"
)

var Balaboba = bot.Command{
	Cmd:         "yalm",
	Description: "Yandex Balaboba",
	Run: func(ctx bot.MessageContext, msg *tg.Message, args []string) error {
		query := strings.Join(args, " ")
		if query == "" {
			return ctx.ReplyText("Think of the beginning of the story")
		}

		sent, err := ctx.Chat.Send(bot.NewMessage("Choose the generation style"), bot.SendOptions[tg.ReplyMarkup]{
			BaseSendOptions: bot.BaseSendOptions{
				ReplyTo: msg.ID,
			},
			ReplyMarkup: bb.StylesKeyboard,
		})
		if err != nil {
			return err
		}

		bb.Reg(ctx.MessageSignature(sent), msg.From.ID, query)
		return nil
	},
}
