package bot

import (
	"runtime"
	"tghwbot/bot/logger"
	"tghwbot/bot/tg"
)

func (b *Bot) makeContext(cmd *Command, msg *tg.Message) *Context {
	return &Context{
		bot:  b,
		log:  b.log.Child(cmd.Cmd),
		cmd:  cmd,
		msg:  msg,
		chat: msg.Chat.ID,
	}
}

// Context type.
type Context struct {
	bot  *Bot
	log  *logger.Logger
	cmd  *Command
	msg  *tg.Message
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

func (c *Context) GetMe() *tg.User {
	return c.bot.Me
}
