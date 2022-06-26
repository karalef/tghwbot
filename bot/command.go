package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Prefix is the character with which commands must begin.
const Prefix = '/'

// Command respresents conversation command.
type Command struct {
	Cmd string
	Run func(*Context, *tgbotapi.Message, []string)

	Description string
	Help        string
	Args        []Arg
}

// Arg type.
type Arg struct {
	Required bool
	Names    []string
	Consts   []string
}
