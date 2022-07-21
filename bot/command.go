package bot

import (
	"strings"
	"tghwbot/bot/tg"
)

// Prefix is the character with which commands must begin.
const Prefix = '/'

// Command respresents conversation command.
type Command struct {
	Cmd string
	Run func(*Context, *tg.Message, []string)

	Description string
	Help        string
	Args        []Arg
}

// Arg type.
type Arg struct {
	Required bool
	Name     string
	Consts   []string
}

func (c *Command) generateHelp() (string, []tg.MessageEntity) {
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
	return sb.String(), []tg.MessageEntity{{
		Type:   "pre",
		Offset: 0,
		Length: sb.Len() - len(c.Description) - 1,
	}}
}
