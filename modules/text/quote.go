package text

import (
	"encoding/json"
	"errors"
	"image"
	"log"
	"net/http"
	"net/url"
	"tghwbot/bot"
	"tghwbot/modules/images"
	"tghwbot/modules/images/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Quote = bot.Command{
	Cmd:         "quote",
	Description: "programming quote",
	Run: func(ctx *bot.Context, _ *tgbotapi.Message, args []string) {
		useCitgen := false
		if len(args) > 0 {
			useCitgen = args[0] == "citgen"
		}
		resp, err := http.Get("https://programming-quotes-api.herokuapp.com/quotes/random")
		if err != nil {
			log.Println("quote:", err.Error())
			err = errors.Unwrap(err)
			ctx.ReplyText(err.Error())
		}
		defer resp.Body.Close()
		var res struct {
			ID     string `json:"id"`
			Author string `json:"author"`
			Text   string `json:"en"`
		}
		json.NewDecoder(resp.Body).Decode(&res)
		if !useCitgen {
			ctx.Chat().Send(bot.MessageConfig{
				MsgText: bot.MsgText{
					Text:      res.Text + "\n\n" + "<b>" + res.Author + "</b>",
					ParseMode: tgbotapi.ModeHTML,
				},
			})
			return
		}

		d, err := generateQuoteCitgen(res.Author, res.Text)
		if err != nil {
			log.Println("quote:", err.Error())
			ctx.ReplyText(err.Error())
		}
		ctx.Chat().SendPhoto(bot.NewPhoto("", tgbotapi.FileBytes{
			Name:  "citgen.png",
			Bytes: d,
		}))
	},
}

func generateQuoteCitgen(author, text string) ([]byte, error) {
	v := url.Values{
		"action":      {"query"},
		"prop":        {"pageimages"},
		"format":      {"json"},
		"pithumbsize": {"250"},
		"titles":      {author},
	}
	resp, err := http.Get("https://en.wikipedia.org/w/api.php?" + v.Encode())
	if err != nil {
		return nil, errors.Unwrap(err)
	}
	defer resp.Body.Close()

	var result struct {
		Query struct {
			Pages map[string]struct {
				Thumbnail struct {
					Source string `json:"source"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnail"`
			} `json:"pages"`
		} `json:"query"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if len(result.Query.Pages) == 0 {
		return nil, errors.New("author not found")
	}
	var src string
	var size int
	for _, v := range result.Query.Pages {
		src = v.Thumbnail.Source
		size = v.Thumbnail.Height
		if size > v.Thumbnail.Width {
			size = v.Thumbnail.Width
		}
		break
	}
	resp, err = http.Get(src)
	if err != nil {
		return nil, errors.Unwrap(err)
	}
	defer resp.Body.Close()
	i, _, err := utils.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	i = utils.Crop(i, image.Rect(0, 0, size, size))
	return images.DefaultCitgen.GeneratePNGBytes(i, author, text, nil)
}
