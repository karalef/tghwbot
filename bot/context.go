package bot

import (
	"io"
	"runtime"
	"tghwbot/logger"

	"gopkg.in/telebot.v3"
)

func (b *Bot) makeContext(cmd *Command, msg *telebot.Message) *Context {
	return &Context{
		api:  b.api,
		log:  b.log.Child(cmd.Cmd),
		cmd:  cmd,
		msg:  msg,
		chat: msg.Chat.ID,
	}
}

// Context type.
type Context struct {
	api  *telebot.Bot
	log  *logger.Logger
	cmd  *Command
	msg  *telebot.Message
	chat int64
}

func (c *Context) err(e error) {
	if e == nil {
		return
	}
	println(e.Error())
	//TODO
	c.Close()
}

// Close stops command execution.
func (c *Context) Close() {
	runtime.Goexit()
}

// Logger returns command logger.
func (c *Context) Logger() *logger.Logger {
	return c.log
}

// OpenChat makes chat interface.
func (c *Context) OpenChat(chatID int64) *Chat {
	return &Chat{
		ctx:    c,
		chatID: chatID,
	}
}

// Chat makes current chat interface.
func (c *Context) Chat() *Chat {
	return c.OpenChat(c.chat)
}

func (c *Context) Download(fileID string) io.ReadCloser {
	body, err := c.api.File(&telebot.File{FileID: fileID})
	c.err(err)
	return body
}

func (c *Context) GetMe() *telebot.User {
	return c.api.Me
}

func (c *Context) GetUserPhotos(userID int64) []telebot.Photo {
	p, err := c.api.ProfilePhotosOf(&telebot.User{ID: userID})
	c.err(err)
	return p
}

func (c *Context) Send(text string, opts ...*telebot.SendOptions) *telebot.Message {
	var o *telebot.SendOptions
	if opts != nil {
		o = opts[0]
	}
	return c.Chat().SendMessage(text, o)
}

func (c *Context) SendClose(text string, opts ...*telebot.SendOptions) {
	c.Send(text, opts...)
	c.Close()
}

func (c *Context) Reply(text string, opts ...*telebot.SendOptions) {
	var o telebot.SendOptions
	if opts != nil {
		o = *opts[0]
	}
	o.ReplyTo = c.msg
	c.Send(text, &o)
}

func (c *Context) ReplyClose(text string, opts ...*telebot.SendOptions) {
	c.Reply(text, opts...)
	c.Close()
}

func (c *Context) ReplyHelp() {
	m, e := generateHelp(c.cmd)
	c.ReplyClose(m, &telebot.SendOptions{Entities: e})
}

func (c *Context) ReplyError(err error) {
	c.log.Error(err.Error())
	c.ReplyClose(err.Error())
}
