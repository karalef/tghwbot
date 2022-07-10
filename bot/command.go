package bot

import "gopkg.in/telebot.v3"

// Prefix is the character with which commands must begin.
const Prefix = '/'

// Command respresents conversation command.
type Command struct {
	Cmd string
	Run func(*Context, *telebot.Message, []string)

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
