package wikihow

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"tghwbot/bot"
	"tghwbot/modules"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Wikihow = bot.Command{
	Cmd:         "wikihow",
	Description: "random wikihow images or steps",
	Args: []bot.Arg{
		{
			Name:   "type",
			Consts: []string{"images", "steps"},
		},
		{
			Name: "count",
		},
	},
	Run: func(ctx *bot.Context, msg *tgbotapi.Message, args []string) {
		if !modules.IsRapidAPIReady() {
			ctx.ReplyText("Service is unavailable")
		}
		steps := false
		c := 3
		if len(args) > 0 {
			s, count := parseArg(args[0])
			if s != nil {
				if *s == "steps" {
					steps = true
				}
				if len(args) > 1 {
					_, count = parseArg(args[1])
					if count != nil {
						c = *count
					}
				}
			} else if count != nil {
				c = *count
			}
		}
		f := randomImages
		if steps {
			f = randomSteps
		}
		result, err := f(c)
		if err != nil || len(result) == 0 {
			e := err.Error()
			if len(result) == 0 {
				e = "no results"
			}
			log.Println("wikihow error:", e)
			ctx.ReplyText(e)
		}
		if steps {
			ctx.ReplyText(strings.Join(result, "\n"))
		}
		pgc := bot.PhotoGroupConfig{
			Photos: make([]tgbotapi.RequestFileData, 0, len(result)),
		}
		for i := range result {
			pgc.Photos = append(pgc.Photos, tgbotapi.FileURL(result[i]))
		}
		ctx.Chat().SendPhotoGroup(pgc)
	},
}

func parseArg(arg string) (*string, *int) {
	switch arg {
	case "steps", "images":
		return &arg, nil
	}
	count, err := strconv.ParseInt(arg, 10, 8)
	if err != nil {
		return nil, nil
	}
	c := int(count)
	if c > 10 {
		c = 10
	} else if c < 1 {
		c = 1
	}
	return nil, &c
}

func randomImages(count int) ([]string, error) {
	return request("images", count)
}

func randomSteps(count int) ([]string, error) {
	return request("steps", count)
}

func request(typ string, count int) ([]string, error) {
	res, err := modules.RapidAPI[map[string]string](nil, modules.RapidAPIRequest{
		Method: http.MethodGet,
		Host:   "hargrimm-wikihow-v1.p.rapidapi.com",
		Path:   "/" + typ,
		Query: url.Values{
			"count": {strconv.Itoa(count)},
		},
	})
	if err != nil {
		return nil, err
	}

	imgs := make([]string, 0, len(*res))
	for _, v := range *res {
		imgs = append(imgs, v)
	}
	return imgs, nil
}
