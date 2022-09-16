package images

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"tghwbot/bot"
	"tghwbot/bot/tg"
	"tghwbot/modules/random"
	"unsafe"
)

var Search = bot.Command{
	Cmd:         "img",
	Description: "random image",
	Run: func(ctx bot.MessageContext, _ *tg.Message, args []string) error {
		q := strings.Join(args, " ")
		if q == "" {
			return ctx.ReplyText("Provide keywords")
		}
		ctx.Chat.SendChatAction(tg.ActionUploadPhoto)
		result, err := searchImages(q, true)
		if err != nil {
			ctx.Logger().Error(err.Error())
			return ctx.ReplyText(err.Error())
		}
		if len(result) == 0 {
			return ctx.ReplyText("No results")
		}
		u := result[random.RandP(len(result), 1.5)]
		p := bot.NewPhoto(tg.FileURL(u))
		p.Caption = q
		return ctx.Reply(p)
	},
}

/*
"https://searx.prvcy.eu/search?",
"https://search.unlocked.link/search?",
"https://search.sapti.me/search?",
"https://searx.zapashcanon.fr/search?",
*/

func searchImages(q string, safe bool) ([]string, error) {
	vals := url.Values{
		"q":          {q},
		"format":     {"json"},
		"categories": {"images"},
		"safesearch": {fmtBool(safe)},
	}
	resp, err := http.Get("https://searx.zapashcanon.fr/search?" + vals.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	var res struct {
		Results []struct {
			ImgSrc string `json:"img_src"`
		} `json:"results"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	for i := range res.Results {
		s := &res.Results[i].ImgSrc
		if strings.HasPrefix(*s, "//") {
			*s = "https:" + *s
		}
	}
	return *(*[]string)(unsafe.Pointer(&res)), nil
}

func fmtBool(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

var cacheTime = 60

func OnInline(ctx bot.QueryContext[bot.InlineAnswer], q *tg.InlineQuery) {
	imgs, err := searchImages(q.Query, false)
	if len(imgs) == 0 {
		if err != nil {
			ctx.Logger().Error(err.Error())
		}
		ctx.Answer(bot.InlineAnswer{})
	}

	n := len(imgs)
	if n > 10 {
		n = 10
	}

	answer := bot.InlineAnswer{
		Results:   make([]tg.InlineQueryResulter, n),
		CacheTime: &cacheTime,
	}

	for i := 0; i < n; i++ {
		answer.Results[i] = tg.InlineQueryResult[tg.InlineQueryResultPhoto]{
			ID: strconv.Itoa(i),
			Result: tg.InlineQueryResultPhoto{
				URL:          imgs[i],
				ThumbnailURL: imgs[i],
				Caption:      q.Query,
			},
		}
	}

	ctx.Answer(answer)
}
