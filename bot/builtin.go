package bot

import (
	"strings"
	"tghwbot/bot/tg"
)

func MakeHelp(b *Bot) *Command {
	h := Command{
		Cmd:         "help",
		Description: "help",
		Args: []Arg{
			{
				Name: "command",
			},
		},
	}
	h.Run = func(ctx *Context, msg *tg.Message, args []string) {
		if len(args) > 0 {
			for _, c := range b.cmds {
				if c.Cmd != args[0] {
					continue
				}
				h, e := c.generateHelp()
				ctx.Reply(h, e...)
			}
			ctx.Reply("command not found")
		}
		var sb strings.Builder
		sb.WriteString("Commands list\n")

		for _, c := range b.cmds {
			sb.WriteByte('\n')
			sb.WriteByte(Prefix)
			sb.WriteString(c.Cmd + " - " + c.Description)
		}
		ctx.Reply(sb.String())
	}
	return &h
}
