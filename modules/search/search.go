package search

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"tghwbot/common"
	"tghwbot/modules/random"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/commands"
)

var searxngURL = os.Getenv("IMG_URL")

// CMD is an "img" command.
var CMD = commands.SimpleCommand{
	Command: "img",
	Desc:    "random image",
	Func: func(ctx *tgot.Message, msg *tg.Message, args []string) error {
		logger := common.Log(ctx)

		q := strings.Join(args, " ")
		if q == "" {
			return ctx.ReplyText("Provide keywords")
		}
		ctx.Chat().SendChatAction(tg.ActionUploadPhoto)
		result, err := searchImages(q)
		if err != nil {
			logger.Err(err).Msg("search images failed")
			if err1 := errors.Unwrap(err); err1 != nil {
				err = err1
			}
			ctx.ReplyText("image search error: " + err.Error())
		}
		if len(result) == 0 {
			return ctx.ReplyText("No results")
		}
		u := result[random.RandP(len(result), 1.5)]
		p := tgot.NewPhoto(tg.FileURL(u))
		p.Caption = q
		return ctx.Reply(p)
	},
}

func searchImages(q string) ([]string, error) {
	vals := url.Values{
		"q":          {q},
		"format":     {"json"},
		"categories": {"images"},
	}
	resp, err := http.Get("http://" + searxngURL + "/search?" + vals.Encode())
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

var cacheTime = 60

func OnInline(ctx tgot.Query[tgot.InlineAnswer], q *tg.InlineQuery) {
	logger := common.Log(ctx)

	imgs, err := searchImages(q.Query)
	if len(imgs) == 0 {
		if err != nil {
			logger.Err(err).Msg("search images failed")
		}
		ctx.Answer(tgot.InlineAnswer{})
		return
	}

	n := len(imgs)
	if n > 10 {
		n = 10
	}

	answer := tgot.InlineAnswer{
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

	err = ctx.Answer(answer)
	if err != nil {
		logger.Err(err).Msg("answer failed")
	}
}
