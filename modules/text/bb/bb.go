package bb

import (
	"errors"
	"strconv"
	"sync"
	"tghwbot/bot"
	"tghwbot/bot/tg"
	"tghwbot/modules"
	"time"

	"github.com/karalef/balaboba"
)

const requestTimeout = time.Minute

var client = balaboba.ClientRus
var mut sync.Mutex

var cacheReq = make(map[bot.MessageSignature]balabobaRequest)

type balabobaRequest struct {
	query   string
	user    int64
	timeout time.Time
}

func cacheGC() {
	for k, r := range cacheReq {
		if time.Now().After(r.timeout) {
			delete(cacheReq, k)
			modules.CallbackRouter.Unreg(k)
		}
	}
}

func Reg(sig bot.MessageSignature, userID int64, query string) {
	cacheGC()
	cacheReq[sig] = balabobaRequest{
		query:   query,
		user:    userID,
		timeout: time.Now().Add(requestTimeout),
	}
	modules.CallbackRouter.Reg(sig, "balaboba styles", styleHandler)
}

func complete(sig bot.MessageSignature, from int64, style balaboba.Style) (*balaboba.Response, bool, error) {
	cacheGC()
	req, ok := cacheReq[sig]
	if !ok || req.user != from {
		return nil, ok, nil
	}
	delete(cacheReq, sig)
	mut.Lock()
	resp, err := client.Generate(nil, req.query, style)
	mut.Unlock()
	return resp, true, err
}

// StylesKeyboard is a keybaord with all balaboba styles.
var StylesKeyboard = &tg.InlineKeyboardMarkup{
	Keyboard: [][]tg.InlineKeyboardButton{
		{
			{Text: "Standart", CallbackData: styleToStr(balaboba.Standart)},
		},
		{
			{Text: "User manual", CallbackData: styleToStr(balaboba.UserManual)},
			{Text: "Recipes", CallbackData: styleToStr(balaboba.Recipes)},
			{Text: "Short stories", CallbackData: styleToStr(balaboba.ShortStories)},
		},
		{
			{Text: "Wikipedia sipmlified", CallbackData: styleToStr(balaboba.WikipediaSipmlified)},
			{Text: "Movie synopses", CallbackData: styleToStr(balaboba.MovieSynopses)},
			{Text: "Folk wisdom", CallbackData: styleToStr(balaboba.FolkWisdom)},
		},
	},
}

func styleToStr(s balaboba.Style) string {
	return strconv.FormatUint(uint64(s), 10)
}

func styleFromStr(str string) balaboba.Style {
	u, _ := strconv.ParseUint(str, 10, 8)
	return balaboba.Style(u)
}

func styleHandler(ctx bot.Context, q *tg.CallbackQuery) (bot.CallbackAnswer, bool, error) {
	sig := ctx.CallbackSignature(q)
	edit := func(text string) error {
		_, err := ctx.EditText(sig, bot.EditText{Text: text}, tg.InlineKeyboardMarkup{
			Keyboard: make([][]tg.InlineKeyboardButton, 0),
		})
		return err
	}
	err := edit("wait up to 20 seconds...")
	if err != nil {
		return bot.CallbackAnswer{}, true, err
	}
	resp, found, err := complete(sig, q.From.ID, styleFromStr(q.Data))
	if !found {
		return bot.CallbackAnswer{}, true, edit("request timeout")
	}
	if err != nil {
		edit("Generation error")
		return bot.CallbackAnswer{
			Text:      err.Error(),
			ShowAlert: true,
		}, true, err
	}

	if resp.Error != 0 {
		edit("Generation error")
		return bot.CallbackAnswer{
			Text:      "Unknown balaboba error",
			ShowAlert: true,
		}, true, errors.New("unknown balaboba error " + strconv.Itoa(resp.Error))
	}
	if resp.BadQuery != 0 {
		resp.Query = ""
		resp.Text = balaboba.BadQueryEng
	}
	return bot.CallbackAnswer{
		Text: "Generation complete",
	}, true, edit(resp.FullText())
}
