package bot

import (
	"strings"
	"tghwbot/bot/tg"
)

var ping = Command{
	Cmd:         "ping",
	Description: "check the bot for availability",
	Run: func(ctx *Context, _ *tg.Message, _ []string) {
		ctx.Send("pong")
	},
}

var help = Command{
	Cmd:         "help",
	Description: "help",
	Args: []Arg{
		{
			Name: "command",
		},
	},
}

func makeHelp(b *Bot) *Command {
	h := help
	h.Run = func(ctx *Context, msg *tg.Message, args []string) {
		if len(args) > 0 {
			for _, c := range b.cmds {
				if c.Cmd != args[0] {
					continue
				}
				h, e := generateHelp(c)
				ctx.ReplyClose(h, &SendOptions{Entities: e})
			}
			ctx.ReplyClose("command not found")
		}
		var sb strings.Builder
		sb.WriteString("Commands list\n")

		for _, c := range b.cmds {
			sb.WriteByte('\n')
			sb.WriteByte(Prefix)
			sb.WriteString(c.Cmd + " - " + c.Description)
		}
		ctx.ReplyClose(sb.String())
	}
	return &h
}

func generateHelp(c *Command) (string, []tg.MessageEntity) {
	sb := strings.Builder{}
	sb.WriteByte(Prefix)
	sb.WriteString(c.Cmd)
	for _, a := range c.Args {
		sb.WriteByte(' ')
		if a.Required {
			sb.WriteByte('[')
		} else {
			sb.WriteByte('{')
		}
		sb.WriteString(a.Name)

		if len(a.Consts) > 0 {
			sb.WriteByte(':')
			sb.WriteString("\"" + strings.Join(a.Consts, "\"|\"") + "\"")
		}
		if a.Required {
			sb.WriteByte(']')
		} else {
			sb.WriteByte('}')
		}
	}
	sb.WriteByte('\n')
	sb.WriteString(c.Description)
	return sb.String(), []tg.MessageEntity{tg.MessageEntity{
		Type:   "pre",
		Offset: 0,
		Length: sb.Len() - len(c.Description) - 1,
	}}
}
