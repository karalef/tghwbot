package bot

import (
	"io"
	"net/http"
	"runtime"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) makeContext(cmd *Command, msg *tgbotapi.Message) *Context {
	return &Context{
		api:  b.api,
		cmd:  cmd,
		msg:  msg,
		chat: msg.Chat.ID,
	}
}

type Context struct {
	api  *tgbotapi.BotAPI
	cmd  *Command
	msg  *tgbotapi.Message
	chat int64
}

func (c *Context) err(e error) {
	if e == nil {
		return
	}
	//TODO
	c.Close()
}

func (c *Context) Close() {
	runtime.Goexit()
}

func (c *Context) OpenChat(chatID int64) *Chat {
	return &Chat{
		ctx:    c,
		chatID: chatID,
	}
}

func (c *Context) OpenChatUsername(username string) *Chat {
	return &Chat{
		ctx:      c,
		username: username,
	}
}

func (c *Context) Chat() *Chat {
	return c.OpenChat(c.chat)
}

func (c *Context) Download(fileID string, out io.Writer) []byte {
	url, err := c.api.GetFileDirectURL(fileID)
	c.err(err)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := c.api.Client.Do(req)
	c.err(err)
	defer resp.Body.Close()
	if out == nil {
		b, err := io.ReadAll(resp.Body)
		c.err(err)
		return b
	}
	_, err = io.Copy(out, resp.Body)
	c.err(err)
	return nil
}

func (c *Context) GetMe() tgbotapi.User {
	return c.api.Self
}

func (c *Context) GetUserPhotos(userID int64) tgbotapi.UserProfilePhotos {
	p, err := c.api.GetUserProfilePhotos(tgbotapi.NewUserProfilePhotos(userID))
	c.err(err)
	return p
}

func (c *Context) ReplyMessage(msg tgbotapi.MessageConfig) {
	c.api.Send(msg)
	c.Close()
}

func (c *Context) ReplyText(text string, ents ...tgbotapi.MessageEntity) {
	m := tgbotapi.NewMessage(c.msg.Chat.ID, text)
	m.Entities = ents
	c.ReplyMessage(m)
	c.Close()
}

func (c *Context) ReplyHelp() {
	m, e := generateHelp(c.cmd)
	c.ReplyText(m, e)
	c.Close()
}
