package bot

import "tghwbot/bot/tg"

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
