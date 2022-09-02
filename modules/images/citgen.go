package images

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"strings"
	"tghwbot/bot"
	"tghwbot/bot/tg"
	"tghwbot/modules/images/fonts"
	"tghwbot/modules/images/utils"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var CitgenCmd = bot.Command{
	Cmd:         "citgen",
	Description: "Генерация цитаты",
	Run: func(ctx bot.MessageContext, msg *tg.Message, args []string) error {
		if msg.ReplyTo == nil {
			return ctx.ReplyText("Ответьте на сообщение")
		}
		from := msg.ReplyTo.From
		text := msg.ReplyTo.Text
		date := msg.ReplyTo.Time()
		caption := ""
		if text == "" {
			return ctx.ReplyText("Сообщение не содержит текста")
		}
		ctx.Chat.SendChatAction(tg.ActionUploadPhoto)

		log := ctx.Logger()
		photo, err := getPhoto(&ctx.Context, from.ID, 200)
		if err != nil {
			log.Warn("citgen: %s", err.Error())
			caption = err.Error()
		}

		user := from.Username
		if user == "" || strings.ReplaceAll(user, " ", "") == "" {
			user = from.FirstName + from.LastName
		}
		data, err := DefaultCitgen.GeneratePNGReader(photo, user, text, date)
		if err != nil {
			log.Error("citgen generate: %s", err.Error())
			return ctx.ReplyText(err.Error())
		}
		p := bot.NewPhoto(tg.FileReader("citgen.png", data))
		p.Caption = caption
		return ctx.Send(p)
	},
}

func getPhoto(ctx *bot.Context, from int64, minSize int) (image.Image, error) {
	ph, err := ctx.GetUserPhotos(from)
	if err != nil || ph.TotalCount == 0 {
		return nil, err
	}

	var fid string
	for _, p := range ph.Photos[0] {
		if p.Height >= minSize {
			fid = p.FileID
			break
		}
	}
	rc, err := ctx.DownloadReaderFile(fid)
	if err != nil {
		return nil, err
	}
	i, _, err := image.Decode(rc)
	rc.Close()
	return i, err
}

var DefaultCitgen = Citgen{
	FontFace:  fonts.GoFontFace(20),
	PhotoSize: 200,
	Width:     700,
	MinHeight: 400,
	Padding:   40,
	BG:        color.Black,
	FG:        color.White,
}

type Citgen struct {
	FontFace  font.Face
	PhotoSize int
	Width     int
	MinHeight int
	Padding   int
	BG, FG    color.Color
}

func (c *Citgen) GeneratePNGReader(photo image.Image, name, quote string, t time.Time) (io.Reader, error) {
	p, err := c.Generate(photo, name, quote, t)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	err = png.Encode(buf, p)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *Citgen) Generate(photo image.Image, name, quote string, t time.Time) (image.Image, error) {
	if photo == nil {
		photo = image.Rect(0, 0, c.PhotoSize, c.PhotoSize)
	} else {
		photo = utils.Resize(photo, c.PhotoSize, c.PhotoSize)
	}
	lineHeight := fonts.Height(c.FontFace)
	textOffsetX := c.Padding*2 + c.PhotoSize
	bottomContentPadding := c.Padding*2 + lineHeight
	lines, textHeight := c.splitLines(quote, c.Width-textOffsetX-c.Padding, c.MinHeight-c.Padding-bottomContentPadding)

	img := image.NewRGBA(image.Rect(0, 0, c.Width, textHeight+c.Padding+bottomContentPadding))

	// bg
	draw.Draw(img, img.Bounds(), image.NewUniform(c.BG), image.ZP, draw.Src)

	// draw photo
	offset := image.Pt(c.Padding, c.Padding+(img.Bounds().Dy()-bottomContentPadding-c.Padding)/2-c.PhotoSize/2)
	utils.DrawClipCircle(img, offset, photo, image.Pt(c.PhotoSize/2, c.PhotoSize/2), c.PhotoSize/2)

	d := font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c.FG),
		Face: c.FontFace,
	}

	// draw quote
	for i := -(len(lines) / 2); i < len(lines)/2+len(lines)%2; i++ {
		d.Dot.X = fixed.I(textOffsetX + 10)
		d.Dot.Y = fixed.I(c.Padding + textHeight/2 + i*lineHeight)
		d.DrawString(lines[i+len(lines)/2])
	}

	// draw name
	d.Dot.X = fixed.I(c.Padding)
	d.Dot.Y = fixed.I(img.Bounds().Dy() - c.Padding)
	d.DrawString(name + " " + copyrightSymbol)

	// draw time
	ft := t.Format("02.01.2006 15:04")
	d.Dot.X = fixed.I(img.Bounds().Dx() - c.Padding - fonts.StringWidth(d.Face, ft))
	d.Dot.Y = fixed.I(img.Bounds().Dy() - c.Padding)
	d.DrawString(ft)

	return img, nil
}

func (c *Citgen) splitLines(s string, width, minHeight int) ([]string, int) {
	var lines []string

	for _, line := range strings.Split(s, "\n") {
		var newLine string
		for _, word := range strings.Split(line, " ") {
			if newLine == "" {
				newLine = word
				continue
			}
			if fonts.StringWidth(c.FontFace, newLine+" "+word) <= width {
				newLine += " " + word
				continue
			}
			lines = append(lines, newLine)
			newLine = word
		}

		if newLine != "" {
			lines = append(lines, newLine)
		}
	}

	lines[0] = "«" + lines[0]
	lines[len(lines)-1] += "»"

	textHeight := fonts.Height(c.FontFace) * len(lines)
	if textHeight < minHeight {
		textHeight = minHeight
	}

	return lines, textHeight
}
