package balaboba

import (
	"strconv"
	"strings"
	"sync"
	"tghwbot/queue"
	"time"

	"github.com/karalef/balaboba"
	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/commands"
	"github.com/karalef/tgot/router"
)

func Command(modulesCtx tgot.Context, r *router.Callbacks) *commands.Command {
	bb := &Balaboba{
		client: balaboba.ClientRus,
		api:    modulesCtx.Child("balaboba"),
		router: r,
	}
	bb.queue = queue.New(bb.complete)

	return &commands.Command{
		Cmd:         "yalm",
		Description: "Yandex Balaboba",
		Func: func(ctx tgot.ChatContext, msg *tg.Message, args []string) error {
			query := strings.Join(args, " ")
			if query == "" {
				return ctx.ReplyE(msg.ID, tgot.NewMessage("Think of the beginning of the story"))
			}

			sent, err := ctx.Reply(msg.ID, tgot.Message{
				Text:        "Choose the generation style",
				ReplyMarkup: StylesKeyboard,
			})
			if err != nil {
				return err
			}

			bb.Reg(tgot.MessageSignature(sent), msg.From.ID, query)
			return nil
		},
	}
}

type Balaboba struct {
	client *balaboba.Client
	mut    sync.Mutex
	queue  *queue.Queue[request]
	router *router.Callbacks

	api tgot.Context
}

func (b *Balaboba) Reg(sig tgot.MsgSignature, from int64, query string) {
	b.router.Reg(sig, &request{
		sig:     sig,
		query:   query,
		user:    from,
		timeout: time.Now().Add(requestTimeout),
		queue:   b.queue,
	})
}

func (b *Balaboba) complete(req request) {
	msg := b.api.OpenMessage(req.sig)
	edit := func(text string) {
		_, err := msg.EditText(tgot.EditText{Text: text}, emptyKeybaord)
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

var emptyKeybaord = tg.InlineKeyboardMarkup{Keyboard: make([][]tg.InlineKeyboardButton, 0)}

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
