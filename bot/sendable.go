package bot

import "tghwbot/bot/tg"

// Sendable interface for Chat.Send.
type Sendable interface {
	what() string
	params() params
}

// Fileable interface.
type Fileable interface {
	Sendable
	files() []File
}

// CaptionData represents caption with entities and parse mode.
type CaptionData struct {
	Caption   string
	ParseMode tg.ParseMode
	Entities  []tg.MessageEntity
}

func (c *CaptionData) embed(p params) {
	p.set("caption", c.Caption)
	p.set("parse_mode", string(c.ParseMode))
	p.set("caption_entities", c.Entities)
}

// BaseFile is common structure for single file with thumbnail.
type BaseFile struct {
	File      *tg.InputFile
	Thumbnail *tg.InputFile
}

func (b *BaseFile) files(field string) []File {
	f := make([]File, 0, 2)
	f = append(f, File{field, b.File})
	if b.Thumbnail != nil {
		f = append(f, File{"thumb", b.Thumbnail})
	}
	return f
}

// SendOptions cointains common send* parameters.
type SendOptions struct {
	DisableNotification      bool
	ProtectContent           bool
	ReplyTo                  int
	AllowSendingWithoutReply bool
	ReplyMarkup              tg.ReplyMarkup
}

func (o *SendOptions) embed(p params) {
	p.set("disable_notification", o.DisableNotification)
	p.set("protect_content", o.ProtectContent)
	p.set("reply_to_message_id", o.ReplyTo)
	p.set("allow_sending_without_reply", o.AllowSendingWithoutReply)
	p.set("reply_markup", o.ReplyMarkup)
}

// NewMessage makes a new message.
func NewMessage(text string) *Message {
	return &Message{
		Text: text,
	}
}

// Message contains information about the message to be sent.
type Message struct {
	Text                  string
	ParseMode             tg.ParseMode
	Entities              []tg.MessageEntity
	DisableWebPagePreview bool
	SendOptions
}

func (Message) what() string {
	return "Message"
}

func (m Message) params() params {
	p := params{}
	p.set("text", m.Text)
	p.set("parse_mode", string(m.ParseMode))
	p.set("entities", m.Entities)
	p.set("disable_web_page_preview", m.DisableWebPagePreview)
	m.SendOptions.embed(p)
	return p
}

// NewPhoto makes a new photo.
func NewPhoto(photo *tg.InputFile) *Photo {
	return &Photo{
		Photo: photo,
	}
}

// Photo contains information about the photo to be sent.
type Photo struct {
	Photo *tg.InputFile
	CaptionData
	SendOptions
}

func (Photo) what() string {
	return "Photo"
}

func (ph Photo) params() params {
	p := params{}
	ph.CaptionData.embed(p)
	ph.SendOptions.embed(p)
	return p
}

func (ph Photo) files() []File {
	return []File{{"photo", ph.Photo}}
}

// NewAudio makes a new audio.
func NewAudio(audio *tg.InputFile) *Audio {
	return &Audio{
		BaseFile: BaseFile{File: audio},
	}
}

// Audio contains information about the audio to be sent.
type Audio struct {
	BaseFile
	CaptionData
	Duration  int
	Performer string
	Title     string
	SendOptions
}

func (Audio) what() string {
	return "Audio"
}

func (a Audio) params() params {
	p := params{}
	a.CaptionData.embed(p)
	p.set("duration", a.Duration)
	p.set("performer", a.Performer)
	p.set("title", a.Title)
	a.SendOptions.embed(p)
	return p
}

func (a Audio) files() []File {
	return a.BaseFile.files("audio")
}

// NewDocument makes a new document.
func NewDocument(document *tg.InputFile) *Document {
	return &Document{
		BaseFile: BaseFile{File: document},
	}
}

// Document contains information about the document to be sent.
type Document struct {
	BaseFile
	CaptionData
	DisableTypeDetection bool
	SendOptions
}

func (Document) what() string {
	return "Document"
}

func (d Document) params() params {
	p := params{}
	d.CaptionData.embed(p)
	p.set("disable_content_type_detection", d.DisableTypeDetection)
	d.SendOptions.embed(p)
	return p
}

func (d Document) files() []File {
	return d.BaseFile.files("document")
}

// NewVideo makes a new video.
func NewVideo(video *tg.InputFile) *Video {
	return &Video{
		BaseFile: BaseFile{File: video},
	}
}

// Video contains information about the video to be sent.
type Video struct {
	BaseFile
	CaptionData
	Duration          int
	Width             int
	Height            int
	SupportsStreaming bool
	SendOptions
}

func (Video) what() string {
	return "Video"
}

func (v Video) params() params {
	p := params{}
	p.set("duration", v.Duration)
	p.set("width", v.Width)
	p.set("height", v.Height)
	v.CaptionData.embed(p)
	p.set("supports_streaming", v.SupportsStreaming)
	v.SendOptions.embed(p)
	return p
}

func (v Video) files() []File {
	return v.BaseFile.files("video")
}

// NewAnimation makes a new animation.
func NewAnimation(animation *tg.InputFile) *Video {
	return &Video{
		BaseFile: BaseFile{File: animation},
	}
}

// Animation contains information about the animation to be sent.
type Animation struct {
	BaseFile
	CaptionData
	Duration int
	Width    int
	Height   int
	SendOptions
}

func (Animation) what() string {
	return "Animation"
}

func (a Animation) params() params {
	p := params{}
	p.set("duration", a.Duration)
	p.set("width", a.Width)
	p.set("height", a.Height)
	a.CaptionData.embed(p)
	a.SendOptions.embed(p)
	return p
}

func (a Animation) files() []File {
	return a.BaseFile.files("animation")
}

// NewVoice makes a new voice.
func NewVoice(voice *tg.InputFile) *Voice {
	return &Voice{
		Voice: voice,
	}
}

// Voice contains information about the voice to be sent.
type Voice struct {
	Voice *tg.InputFile
	CaptionData
	Duration int
	SendOptions
}

func (Voice) what() string {
	return "Voice"
}

func (v Voice) params() params {
	p := params{}
	v.CaptionData.embed(p)
	p.set("duration", v.Duration)
	v.SendOptions.embed(p)
	return p
}

func (v Voice) files() []File {
	return []File{{"voice", v.Voice}}
}

// NewVideoNote makes a new video note.
func NewVideoNote(videoNote *tg.InputFile) *VideoNote {
	return &VideoNote{
		BaseFile: BaseFile{File: videoNote},
	}
}

// VideoNote contains information about the video note to be sent.
type VideoNote struct {
	BaseFile
	Duration int
	Length   int
	SendOptions
}

func (VideoNote) what() string {
	return "VideoNote"
}

func (v VideoNote) params() params {
	p := params{}
	p.set("duration", v.Duration)
	p.set("length", v.Length)
	v.SendOptions.embed(p)
	return p
}

func (v VideoNote) files() []File {
	return v.BaseFile.files("video_note")
}

// NewDice makes a new dice.
func NewDice(dice tg.DiceEmoji) *Dice {
	return &Dice{
		Emoji: dice,
	}
}

// Dice contains information about the dice to be sent.
type Dice struct {
	Emoji tg.DiceEmoji
	SendOptions
}

func (Dice) what() string {
	return "Dice"
}

func (d Dice) params() params {
	p := params{}
	p.set("emoji", string(d.Emoji))
	d.SendOptions.embed(p)
	return p
}
