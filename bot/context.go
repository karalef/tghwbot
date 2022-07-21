package bot

import (
	"context"
	"strconv"
	"tghwbot/bot/logger"
	"tghwbot/bot/tg"
)

func (b *Bot) makeContext(cmd *Command, msg *tg.Message) *Context {
	c := Context{
		bot: b,
		cmd: cmd,
		msg: msg,
	}
	c.Chat = c.OpenChat(msg.Chat.ID)
	return &c
}

// Context type.
type Context struct {
	bot  *Bot
	cmd  *Command
	msg  *tg.Message
	Chat Chat
}

func (c *Context) getBot() *Bot {
	return c.bot
}

func (c *Context) caller() string {
	return c.cmd.Cmd
}

type commonContext interface {
	getBot() *Bot
	caller() string
}

func api[T any](c commonContext, method string, p params, files ...File) T {
	bot := c.getBot()
	var result T
	var err error
	if len(files) > 0 {
		result, err = uploadFiles[T](bot, method, p, files)
	} else {
		result, err = performRequest[T](bot, method, p)
	}
	switch err.(type) {
	case nil:
		return result
	case *tg.APIError:
		bot.log.Warn("from '%s'\n%s", c.caller(), err.Error())
		bot.closeExecution()
		return result
	}

	switch err {
	case context.Canceled, context.DeadlineExceeded:
	default:
		bot.log.Error(err.Error())
	}
	bot.closeExecution()
	return result
}

// Close stops command execution.
func (c *Context) Close() {
	c.bot.closeExecution()
}

// Logger returns command logger.
func (c *Context) Logger() *logger.Logger {
	return c.bot.log.Child(c.cmd.Cmd)
}

// OpenChat makes chat interface.
func (c *Context) OpenChat(chatID int64) Chat {
	return c.OpenChatUsername(strconv.FormatInt(chatID, 10))
}

// OpenChatUsername makes chat interface by username.
func (c *Context) OpenChatUsername(username string) Chat {
	return Chat{
		ctx:    c,
		chatID: username,
	}
}

// Reply sends message to the current chat and closes context.
func (c *Context) Reply(text string, entities ...tg.MessageEntity) {
	c.Chat.Send(Message{
		Text:     text,
		Entities: entities,
		SendOptions: SendOptions{
			ReplyTo: c.msg.ID,
		},
	})
	c.Close()
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
