package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Chat struct {
	ctx      *Context
	chatID   int64
	username string
}

type MsgText struct {
	Text      string
	Entities  []tgbotapi.MessageEntity
	ParseMode string
}

type MessageConfig struct {
	MsgText
	DisableWebPagePreview    bool
	ReplyToMessageID         int
	ReplyMarkup              interface{}
	DisableNotification      bool
	AllowSendingWithoutReply bool
}

func (c *Chat) Send(msg MessageConfig) tgbotapi.Message {
	m, err := c.ctx.api.Send(tgbotapi.MessageConfig{
		Text:                  msg.Text,
		Entities:              msg.Entities,
		ParseMode:             msg.ParseMode,
		DisableWebPagePreview: msg.DisableWebPagePreview,
		BaseChat: tgbotapi.BaseChat{
			ChatID:                   c.chatID,
			ChannelUsername:          c.username,
			ReplyToMessageID:         msg.ReplyToMessageID,
			ReplyMarkup:              msg.ReplyMarkup,
			DisableNotification:      msg.DisableNotification,
			AllowSendingWithoutReply: msg.AllowSendingWithoutReply,
		},
	})
	c.ctx.err(err)
	return m
}

func (c *Chat) Edit(msgID int, text string) tgbotapi.Message {
	m, err := c.ctx.api.Send(tgbotapi.NewEditMessageText(c.chatID, msgID, text))
	c.ctx.err(err)
	return m
}

func NewPhoto(caption string, file tgbotapi.RequestFileData) PhotoConfig {
	return PhotoConfig{
		MsgText: MsgText{
			Text: caption,
		},
		File: file,
	}
}

type PhotoConfig struct {
	MsgText
	ReplyToMessageID         int
	ReplyMarkup              interface{}
	DisableNotification      bool
	AllowSendingWithoutReply bool
	File                     tgbotapi.RequestFileData
}

func (c *Chat) SendPhoto(msg PhotoConfig) tgbotapi.Message {
	m, err := c.ctx.api.Send(tgbotapi.PhotoConfig{
		BaseFile: tgbotapi.BaseFile{
			BaseChat: tgbotapi.BaseChat{
				ChatID:                   c.chatID,
				ChannelUsername:          c.username,
				ReplyToMessageID:         msg.ReplyToMessageID,
				ReplyMarkup:              msg.ReplyMarkup,
				DisableNotification:      msg.DisableNotification,
				AllowSendingWithoutReply: msg.AllowSendingWithoutReply,
			},
			File: msg.File,
		},
		Caption:         msg.Text,
		CaptionEntities: msg.Entities,
		ParseMode:       msg.ParseMode,
	})
	c.ctx.err(err)
	return m
}

type PhotoGroupConfig struct {
	Photos              []tgbotapi.RequestFileData
	DisableNotification bool
	ReplyToMessageID    int
}

func (c *Chat) SendPhotoGroup(pg PhotoGroupConfig) []tgbotapi.Message {
	mgc := tgbotapi.NewMediaGroup(c.chatID, make([]interface{}, 0, len(pg.Photos)))
	for i := range pg.Photos {
		mgc.Media = append(mgc.Media, tgbotapi.NewInputMediaPhoto(pg.Photos[i]))
	}
	m, err := c.ctx.api.SendMediaGroup(mgc)
	c.ctx.err(err)
	return m
}
