package bot

import (
	"context"
	"runtime"
	"strconv"
	"tghwbot/bot/logger"
	"tghwbot/bot/tg"
)

func (b *Bot) makeContext(cmd *Command, msg *tg.Message) *Context {
	return &Context{
		bot:  b,
		cmd:  cmd,
		msg:  msg,
		chat: msg.Chat.ID,
	}
}

// Context type.
type Context struct {
	bot  *Bot
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

func api[T any](c *Context, method string, p params) T {
	result, err := performRequest[T](c.bot, method, p)
	switch err.(type) {
	case nil:
		return result
	case *tg.APIError:
		c.bot.log.Warn("from '%s'\n%s", c.cmd.Cmd, err.Error())
		c.Close()
		return result
	}

	switch err {
	case context.Canceled, context.DeadlineExceeded:
	default:
		c.bot.log.Error(err.Error())
	}
	c.Close()
	return result
}

// Close stops command execution.
func (c *Context) Close() {
	runtime.Goexit()
}

// Logger returns command logger.
func (c *Context) Logger() *logger.Logger {
	return c.bot.log.Child(c.cmd.Cmd)
}

// OpenChat makes chat interface.
func (c *Context) OpenChat(chatID int64) *Chat {
	return c.OpenChatUsername(strconv.FormatInt(chatID, 10))
}

// OpenChatUsername makes chat interface by username.
func (c *Context) OpenChatUsername(username string) *Chat {
	return &Chat{
		ctx:    c,
		chatID: username,
	}
}

// Chat makes current chat interface.
func (c *Context) Chat() *Chat {
	return c.OpenChat(c.chat)
}

// GetMe returns basic information about the bot.
func (c *Context) GetMe() tg.User {
	return *c.bot.Me
}

// GetUserPhotos returns a list of profile pictures for a user.
func (c *Context) GetUserPhotos(userID int64) *tg.UserProfilePhotos {
	p := params{}.set("user_id", userID)
	return api[*tg.UserProfilePhotos](c, "getUserProfilePhotos", p)
}

// GetFile returns basic information about a file
// and prepares it for downloading.
func (c *Context) GetFile(fileID string) *tg.File {
	p := params{}.set("file_id", fileID)
	return api[*tg.File](c, "getFile", p)
}
