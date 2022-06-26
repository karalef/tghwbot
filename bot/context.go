package bot

import (
	"net/http"
	"runtime"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) makeContext(cmd *Command, chat *tgbotapi.Chat) *Context {
	return &Context{
		api:  b.api,
		cmd:  cmd,
		chat: chat,
	}
}

type Context struct {
	api  *tgbotapi.BotAPI
	cmd  *Command
	chat *tgbotapi.Chat
}

func (c *Context) Close() {
	runtime.Goexit()
}

func (c *Context) API() *tgbotapi.BotAPI {
	return c.api
}

func (c *Context) Download(file tgbotapi.FileConfig) (*http.Response, error) {
	f, err := c.api.GetFile(file)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(http.MethodGet, f.Link(c.api.Token), nil)
	return c.api.Client.Do(req)
}

func (c *Context) Send(ch tgbotapi.Chattable) {
	c.api.Send(ch)
}

func (c *Context) SendText(text string, ents ...tgbotapi.MessageEntity) {
	m := tgbotapi.NewMessage(c.chat.ID, text)
	m.Entities = ents
	c.Send(m)
}

func (c *Context) ReplyMessage(msg tgbotapi.MessageConfig) {
	c.Send(msg)
	c.Close()
}

func (c *Context) ReplyText(text string, ents ...tgbotapi.MessageEntity) {
	c.SendText(text, ents...)
	c.Close()
}

func (c *Context) ReplyHelp() {
	m, e := generateHelp(c.cmd)
	c.SendText(m, e)
	c.Close()
}
