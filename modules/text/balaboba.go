package text

import (
	"strconv"
	"strings"
	"sync"
	"tghwbot/modules"
	"tghwbot/queue"
	"time"

	"github.com/karalef/balaboba"
	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/commands"
)

var bb *Balaboba

func InitBalaboba() {
	bb = &Balaboba{
		client: balaboba.ClientRus,
		api:    modules.API.Child("balaboba"),
	}
	bb.queue = queue.New(bb.complete)
}

var BalabobaCmd = commands.Command{
	Cmd:         "yalm",
	Description: "Yandex Balaboba",
	Func: func(ctx tgot.ChatContext, msg *tg.Message, args []string) error {
		query := strings.Join(args, " ")
		if query == "" {
			return modules.ReplyText(ctx, msg, "Think of the beginning of the story")
		}

		sent, err := ctx.Send(tgot.NewMessage("Choose the generation style"), tgot.SendOptions[tg.ReplyMarkup]{
			ReplyTo:     msg.ID,
			ReplyMarkup: StylesKeyboard,
		})
		if err != nil {
			return err
		}

		bb.Reg(tgot.MessageSignature(sent), msg.From.ID, query)
		return nil
	},
}

const requestTimeout = time.Minute

type Balaboba struct {
	client *balaboba.Client
	mut    sync.Mutex
	queue  *queue.Queue[balabobaRequest]

	api tgot.Context
}

func (b *Balaboba) Reg(sig tgot.MsgSignature, from int64, query string) {
	modules.CallbackRouter.Reg(sig, &balabobaRequest{
		sig:     sig,
		query:   query,
		user:    from,
		timeout: time.Now().Add(requestTimeout),
	})
}

func (b *Balaboba) complete(req balabobaRequest) {
	msg := b.api.OpenMessage(req.sig)
	edit := func(text string) {
		_, err := msg.EditText(tgot.EditText{Text: text}, tg.InlineKeyboardMarkup{
			Keyboard: make([][]tg.InlineKeyboardButton, 0),
		})
		if err != nil {
			b.api.Logger().Error(err.Error())
		}
	}
	b.mut.Lock()
	defer b.mut.Unlock()
	edit("wait up to 20 seconds...")
	var text string
	resp, err := b.client.Generate(req.query, req.style)
	if err != nil {
		text = "Generation error"
		b.api.Logger().Error(err.Error())
	} else {
		text = resp.Text
	}
	edit(text)
}

type balabobaRequest struct {
	sig     tgot.MsgSignature
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
	if bb.queue.Len() > 0 {
		msg := ctx.OpenMessage(tgot.CallbackSignature(q))
		_, err = msg.EditText(tgot.EditText{
			Text: "request is added to the queue",
		}, tg.InlineKeyboardMarkup{
			Keyboard: make([][]tg.InlineKeyboardButton, 0),
		})
		if err != nil {
			return tgot.CallbackAnswer{}, true, err
		}
	}

	u, _ := strconv.ParseUint(q.Data, 10, 8)
	r.style = balaboba.Style(u)
	bb.queue.Push(*r)

	return tgot.CallbackAnswer{}, true, err
}

func (r balabobaRequest) Timeout() time.Time {
	return r.timeout
}

func (r balabobaRequest) Close(ctx tgot.Context, sig tgot.MsgSignature) error {
	_, err := ctx.OpenMessage(sig).EditText(tgot.EditText{Text: "request timeout"}, tg.InlineKeyboardMarkup{
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
