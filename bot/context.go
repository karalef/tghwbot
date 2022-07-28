package bot

import (
	"context"
	"io"
	"strconv"
	"tghwbot/bot/internal"
	"tghwbot/bot/logger"
	"tghwbot/bot/tg"
)

type contextBase struct {
	bot       *Bot
	Chat      Chat
	getCaller func() string
}

func (c *contextBase) caller() string {
	if c.getCaller != nil {
		return c.getCaller()
	}
	return "unknown"
}

func (c *contextBase) getBot() *Bot {
	return c.bot
}

// Close stops handler execution.
func (c *contextBase) Close() {
	c.bot.closeExecution()
}

// Logger returns command logger.
func (c *contextBase) Logger() *logger.Logger {
	return c.bot.log.Child(c.caller())
}

// OpenChat makes chat interface.
func (c *contextBase) OpenChat(chatID int64) Chat {
	return c.OpenChatUsername(strconv.FormatInt(chatID, 10))
}

// OpenChatUsername makes chat interface by username.
func (c *contextBase) OpenChatUsername(username string) Chat {
	return Chat{
		ctx:    c,
		chatID: username,
	}
}

// GetMe returns basic information about the bot.
func (c *contextBase) GetMe() tg.User {
	return *c.bot.Me
}

// GetUserPhotos returns a list of profile pictures for a user.
func (c *contextBase) GetUserPhotos(userID int64) *tg.UserProfilePhotos {
	p := params{}.set("user_id", userID)
	return api[*tg.UserProfilePhotos](c, "getUserProfilePhotos", p)
}

// GetFile returns basic information about a file
// and prepares it for downloading.
func (c *contextBase) GetFile(fileID string) *tg.File {
	p := params{}.set("file_id", fileID)
	return api[*tg.File](c, "getFile", p)
}

// DownloadReader downloads file as io.ReadCloser from Telegram servers.
func (c *contextBase) DownloadReader(f *tg.File) (io.ReadCloser, error) {
	return c.bot.downloadFile(f.FilePath)
}

// Download downloads file from Telegram servers.
func (c *contextBase) Download(f *tg.File) ([]byte, error) {
	rc, err := c.bot.downloadFile(f.FilePath)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

// DownloadReaderFile downloads file as io.ReadCloser from Telegram servers
// by file id.
func (c *contextBase) DownloadReaderFile(fid string) (io.ReadCloser, error) {
	return c.DownloadReader(c.GetFile(fid))
}

// DownloadFile downloads file from Telegram servers by file id.
func (c *contextBase) DownloadFile(fid string) ([]byte, error) {
	return c.Download(c.GetFile(fid))
}

type commonContext interface {
	getBot() *Bot
	caller() string
	Close()
}

func api[T any](c commonContext, method string, p params, files ...File) T {
	bot := c.getBot()
	result, err := performRequest[T](bot, method, p, files...)
	switch err.(type) {
	case nil:
		return result
	case *tg.APIError:
		bot.log.Error("%s\n%s\n%s", err.Error(), c.caller(), backtrace())
		c.Close()
		return result
	}

	switch err {
	case context.Canceled, context.DeadlineExceeded:
	default:
		bot.log.Error("%s\n%s\n%s", err.Error(), c.caller(), backtrace())
	}
	c.Close()
	return result
}

func backtrace() string {
	return internal.FramesString(internal.BackTrace(2, 2), true)
}

func (b *Bot) makeContext(cmd *Command, msg *tg.Message) *Context {
	c := Context{
		contextBase: contextBase{
			bot: b,
		},
		cmd: cmd,
		msg: msg,
	}
	c.getCaller = c.caller
	c.Chat = c.OpenChat(msg.Chat.ID)
	return &c
}

// Context type.
type Context struct {
	contextBase
	cmd *Command
	msg *tg.Message
}

func (c *Context) caller() string {
	return "command '" + c.cmd.Cmd + "'"
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
