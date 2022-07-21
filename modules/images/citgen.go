package images

import (
	"bytes"
	"flag"
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
	Run: func(ctx *bot.Context, msg *tg.Message, args []string) {
		if msg.ReplyTo == nil {
			ctx.Reply("Ответьте на сообщение")
		}
		from := msg.ReplyTo.From
		text := msg.ReplyTo.Text
		date := msg.ReplyTo.Time()
		caption := ""
		if text == "" {
			ctx.Reply("Сообщение не содержит текста")
		}
		ctx.Chat.Send(bot.ChatAction(tg.ActionUploadPhoto))

		log := ctx.Logger()
		photo, err := getPhoto(ctx, from.ID, 200)
		if err != nil {
			log.Warn("citgen: %s", err.Error())
			caption = err.Error()
		}

		config := DefaultCitgen
		if len(args) > 0 {
			citgenFlagSet := flag.NewFlagSet("citgen", flag.ContinueOnError)
			citgenFlagSet.BoolVar(&config.PhotoQuad, "q", false, "")
			citgenFlagSet.Parse(args)
		}
		user := from.Username
		if user == "" || strings.ReplaceAll(user, " ", "") == "" {
			user = from.FirstName + from.LastName
		}
		data, err := config.GeneratePNGReader(photo, user, text, date, from.ID == msg.From.ID)
		if err != nil {
			log.Error("citgen generate: %s", err.Error())
			ctx.Reply(err.Error())
		}
		p := bot.NewPhoto(tg.FileReader("citgen.png", data))
		p.Caption = caption
		ctx.Chat.Send(p)
	},
}

func getPhoto(ctx *bot.Context, from int64, minSize int) (image.Image, error) {
	ph := ctx.GetUserPhotos(from)
	if ph.TotalCount == 0 {
		return nil, nil
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
	i, _, err := utils.Decode(rc)
	rc.Close()
	return i, err
}

var DefaultCitgen = Citgen{
	FontFace:  fonts.RobotoFace(20),
	PhotoSize: 200,
	PhotoQuad: false,
	Width:     700,
	MinHeight: 400,
	Padding:   40,
	BG:        color.Black,
	FG:        color.White,
}

type Citgen struct {
	FontFace  font.Face
	PhotoSize int
	PhotoQuad bool
	Width     int
	MinHeight int
	Padding   int
	BG, FG    color.Color
}

func (c *Citgen) GeneratePNGReader(photo image.Image, name, quote string, t time.Time, self bool) (io.Reader, error) {
	p, err := c.Generate(photo, name, quote, t, self)
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

func (c *Citgen) Generate(photo image.Image, name, quote string, t time.Time, self bool) (image.Image, error) {
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
	if c.PhotoQuad {
		draw.Draw(img, img.Bounds().Add(offset), photo, image.ZP, draw.Over)
	} else {
		utils.DrawClipCircle(img, offset, photo, image.Pt(c.PhotoSize/2, c.PhotoSize/2), c.PhotoSize/2)
	}

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
			if newLine != "" {
				word = " " + word
			}
			if fonts.StringWidth(c.FontFace, newLine+word) <= width {
				newLine += word
				continue
			}
			lines = append(lines, newLine)
			newLine = word
		}

		if newLine != "" {
			lines = append(lines, strings.TrimSpace(newLine))
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
