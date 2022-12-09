package modules

import (
	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

// Reply replies to message.
func Reply(ctx tgot.ChatContext, to *tg.Message, s tgot.Sendable) error {
	return ctx.SendE(s, tgot.SendOptions[tg.ReplyMarkup]{
		ReplyTo: to.ID,
	})
}

// ReplyText replies to message just with text.
func ReplyText(ctx tgot.ChatContext, to *tg.Message, text string) error {
	return Reply(ctx, to, tgot.NewMessage(text))
}
