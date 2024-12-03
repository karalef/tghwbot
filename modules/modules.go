package modules

import "github.com/karalef/tgot"

// ReplyText replies to the message with text.
func ReplyText(ctx *tgot.Message, text string) error {
	return ctx.Chat().ReplyE(tgot.ReplyTo(ctx.ID().MessageID()), tgot.NewText(text))
}
