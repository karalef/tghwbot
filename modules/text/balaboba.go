package text

import (
	"strconv"
	"strings"
	"sync"
	"tghwbot/modules"
	"time"

	"github.com/karalef/balaboba"
	"github.com/karalef/tgot"
	"github.com/karalef/tgot/tg"
)

var bb *Balaboba

func InitBalaboba() {
	bb = &Balaboba{
		client: balaboba.ClientRus,
		queue:  make(chan balabobaRequest, 20),
		api:    modules.API.Child("balaboba"),
	}
	go bb.listen()
}

var BalabobaCmd = tgot.Command{
	Cmd:         "yalm",
	Description: "Yandex Balaboba",
	Run: func(ctx tgot.MessageContext, msg *tg.Message, args []string) error {
		query := strings.Join(args, " ")
		if query == "" {
			return ctx.ReplyText("Think of the beginning of the story")
		}

		sent, err := ctx.Chat.Send(tgot.NewMessage("Choose the generation style"), tgot.SendOptions[tg.ReplyMarkup]{
			BaseSendOptions: tgot.BaseSendOptions{
				ReplyTo: msg.ID,
			},
			ReplyMarkup: StylesKeyboard,
		})
		if err != nil {
			return err
		}

		bb.Reg(ctx.MessageSignature(sent), msg.From.ID, query)
		return nil
	},
}

const requestTimeout = time.Minute

type Balaboba struct {
	client *balaboba.Client
	mut    sync.Mutex
	queue  chan balabobaRequest

	api tgot.Context
}

func (b *Balaboba) Reg(sig tgot.MessageSignature, from int64, query string) {
	modules.CallbackRouter.Reg(sig, &balabobaRequest{
		sig:     sig,
		query:   query,
		user:    from,
		timeout: time.Now().Add(requestTimeout),
	})
}

func (b *Balaboba) listen() {
	for {
		select {
		case r, ok := <-b.queue:
			if !ok {
				return
			}
			b.complete(r)
		}
	}
}

func (b *Balaboba) complete(req balabobaRequest) {
	edit := func(text string) {
		_, err := bb.api.EditText(req.sig, tgot.EditText{Text: text}, tg.InlineKeyboardMarkup{
			Keyboard: make([][]tg.InlineKeyboardButton, 0),
		})
		if err != nil {
			bb.api.Logger().Error(err.Error())
		}
	}
	b.mut.Lock()
	defer b.mut.Unlock()
	edit("wait up to 20 seconds...")
	var text string
	resp, err := b.client.Generate(req.query, req.style)
	if err != nil {
		text = "Generation error"
		bb.api.Logger().Error(err.Error())
	} else {
		text = resp.Text
	}
	edit(text)
}

type balabobaRequest struct {
	sig     tgot.MessageSignature
	query   string
	style   balaboba.Style
	user    int64
	timeout time.Time
}

func (balabobaRequest) Name() string {
	return "balaboba"
}

func (r *balabobaRequest) Handle(ctx tgot.Context, q *tg.CallbackQuery) (tgot.CallbackAnswer, bool, error) {
	if q.From.ID != r.user {
		return tgot.CallbackAnswer{}, false, nil
	}
	var err error
	if len(bb.queue) == 0 {
		_, err = ctx.EditText(ctx.CallbackSignature(q), tgot.EditText{
			Text: "request is added to the queue",
		}, tg.InlineKeyboardMarkup{
			Keyboard: make([][]tg.InlineKeyboardButton, 0),
		})
	}

	u, _ := strconv.ParseUint(q.Data, 10, 8)
	r.style = balaboba.Style(u)
	bb.queue <- *r

	return tgot.CallbackAnswer{}, true, err
}

func (r balabobaRequest) Timeout() time.Time {
	return r.timeout
}

func (r balabobaRequest) Close(ctx tgot.Context, sig tgot.MessageSignature) error {
	_, err := ctx.EditText(sig, tgot.EditText{Text: "request timeout"}, tg.InlineKeyboardMarkup{
		Keyboard: make([][]tg.InlineKeyboardButton, 0),
	})
	return err
}

// StylesKeyboard is a keybaord with all balaboba styles.
var StylesKeyboard = func() *tg.InlineKeyboardMarkup {
	str := func(s balaboba.Style) string {
		return strconv.FormatUint(uint64(s), 10)
	}
	return &tg.InlineKeyboardMarkup{
		Keyboard: [][]tg.InlineKeyboardButton{
			{
				{Text: "Standart", CallbackData: str(balaboba.Standart)},
			},
			{
				{Text: "User manual", CallbackData: str(balaboba.UserManual)},
				{Text: "Recipes", CallbackData: str(balaboba.Recipes)},
				{Text: "Short stories", CallbackData: str(balaboba.ShortStories)},
			},
			{
				{Text: "Wikipedia sipmlified", CallbackData: str(balaboba.WikipediaSipmlified)},
				{Text: "Movie synopses", CallbackData: str(balaboba.MovieSynopses)},
				{Text: "Folk wisdom", CallbackData: str(balaboba.FolkWisdom)},
			},
		},
	}
}()
