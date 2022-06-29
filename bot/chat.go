package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Chat struct {
	ctx      *Context
	chatID   int64
	username string
}

func (c *Chat) SendMessage(msg tgbotapi.MessageConfig) tgbotapi.Message {
	msg.ChatID = c.chatID
	msg.ChannelUsername = c.username
	m, err := c.ctx.api.Send(msg)
	c.ctx.err(err)
	return m
}

func (c *Chat) SendPhoto(msg tgbotapi.PhotoConfig) tgbotapi.Message {
	msg.ChatID = c.chatID
	msg.ChannelUsername = c.username
	m, err := c.ctx.api.Send(msg)
	c.ctx.err(err)
	return m
}
