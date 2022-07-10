package bot

import (
	"io"

	"gopkg.in/telebot.v3"
)

type Chat struct {
	ctx    *Context
	chatID int64
}

func (c *Chat) send(what interface{}, opts *telebot.SendOptions) *telebot.Message {
	m, err := c.ctx.api.Send(telebot.ChatID(c.chatID), what, opts)
	c.ctx.err(err)
	return m
}

func (c *Chat) SendMessage(text string, opts *telebot.SendOptions) *telebot.Message {
	return c.send(text, opts)
}

func (c *Chat) SendPhoto(ph *telebot.Photo, opts *telebot.SendOptions) *telebot.Message {
	return c.send(ph, opts)
}

func NewPhoto(data io.Reader, caption string) *telebot.Photo {
	return &telebot.Photo{
		File:    telebot.FromReader(data),
		Caption: caption,
	}
}
