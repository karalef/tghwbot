package bot

import (
	"io"
	"log"
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

// Context type.
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
	log.Println("telegram: ", e.Error())
	//TODO
	c.Close()
}

// Close stops command execution.
func (c *Context) Close() {
	runtime.Goexit()
}

// OpenChat makes chat interface.
func (c *Context) OpenChat(chatID int64) *Chat {
	return &Chat{
		ctx:    c,
		chatID: chatID,
	}
}

// OpenChatUsername makes chat interface by username intead of id.
func (c *Context) OpenChatUsername(username string) *Chat {
	return &Chat{
		ctx:      c,
		username: username,
	}
}

// Chat makes current chat interface.
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

func (c *Context) MessageClose(msg MessageConfig) {
	c.Chat().Send(msg)
	c.Close()
}

func (c *Context) TextClose(text string) {
	c.Chat().Send(MessageConfig{MsgText: MsgText{Text: text}})
	c.Close()
}

func (c *Context) ReplyText(text string, ents ...tgbotapi.MessageEntity) {
	c.Chat().Send(MessageConfig{
		MsgText: MsgText{
			Text:     text,
			Entities: ents,
		},
		ReplyToMessageID: c.msg.MessageID,
	})
	c.Close()
}

func (c *Context) ReplyHelp() {
	m, e := generateHelp(c.cmd)
	c.ReplyText(m, e)
	c.Close()
}
