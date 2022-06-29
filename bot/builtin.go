package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ping = Command{
	Cmd:         "ping",
	Description: "check the bot for availability",
	Run: func(ctx *Context, _ *tgbotapi.Message, _ []string) {
		ctx.ReplyText("pong")
	},
}

func makeHelp(b *Bot) *Command {
	return &Command{
		Cmd:         "help",
		Description: "help",
		Args: []Arg{
			{
				Names: []string{"command"},
			},
		},
		Run: func(ctx *Context, _ *tgbotapi.Message, args []string) {
			if len(args) > 0 {
				for _, c := range b.cmds {
					if c.Cmd != args[0] {
						continue
					}
					h, e := generateHelp(c)
					ctx.ReplyText(h, e)
				}
				ctx.ReplyText("command not found")
			}
			var sb strings.Builder
			sb.WriteString("Commands list\n")

			for _, c := range b.cmds {
				sb.WriteByte('\n')
				sb.WriteByte(Prefix)
				sb.WriteString(c.Cmd + " - " + c.Description)
			}
			ctx.ReplyText(sb.String())
		},
	}
}

func generateHelp(c *Command) (string, tgbotapi.MessageEntity) {
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
		if len(a.Names) > 0 {
			sb.WriteString(strings.Join(a.Names, "|"))
			sb.WriteByte('|')
		}
		if len(a.Consts) > 0 {
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
	return sb.String(), tgbotapi.MessageEntity{
		Type:   "pre",
		Offset: 0,
		Length: sb.Len() - len(c.Description) - 1,
	}
}
