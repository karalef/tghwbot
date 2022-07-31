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
	Run: func(ctx *bot.Context, _ *tg.Message, args []string) {
		q := strings.Join(args, " ")
		if q == "" {
			ctx.Reply("Provide keywords")
		}
		ctx.Chat.SendChatAction(tg.ActionUploadPhoto)
		result, err := searchImages(q, false)
		if err != nil {
			ctx.Logger().Error(err.Error())
			ctx.Reply(err.Error())
		}
		if len(result) == 0 {
			ctx.Reply("No results")
		}
		u := result[random.RandP(len(result), 1.5)]
		p := bot.NewPhoto(tg.FileURL(u))
		p.Caption = q
		ctx.Chat.Send(p)
	},
}

func searchImages(q string, safe bool) ([]string, error) {
	vals := url.Values{
		"q":          {q},
		"format":     {"json"},
		"categories": {"images"},
		"safesearch": {"0"},
	}
	if safe {
		vals.Set("safesearch", "1")
	}
	resp, err := http.Get("https://searx.zapashcanon.fr/search?" + vals.Encode())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	var res struct {
		Results []struct {
			ImgSrc string `json:"img_src"`
		} `json:"results"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	resp.Body.Close()
	for i := range res.Results {
		s := &res.Results[i].ImgSrc
		if strings.HasPrefix(*s, "//") {
			*s = "https:" + *s
		}
	}
	return *(*[]string)(unsafe.Pointer(&res)), nil
}

func OnInline(ctx *bot.InlineContext, q *tg.InlineQuery) {
	imgs, err := searchImages(q.Query, false)
	if len(imgs) == 0 {
		if err != nil {
			ctx.Logger().Error(err.Error())
		}
		ctx.Answer(&bot.InlineAnswer{})
	}

	n := len(imgs)
	if n > 10 {
		n = 10
	}

	results := make([]tg.InlineQueryResult, n)
	for i := 0; i < n; i++ {
		results[i].ID = strconv.Itoa(i)
		results[i].Result = tg.InlineQueryResultPhoto{
			URL:          imgs[i],
			ThumbnailURL: imgs[i],
			Caption:      q.Query,
		}
	}

	ctx.Answer(&bot.InlineAnswer{
		Results: results,
	}, 60)
}
