package balaboba

import (
	"strconv"
	"tghwbot/queue"
	"time"

	"github.com/karalef/balaboba"
	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

const requestTimeout = time.Minute

type request struct {
	sig     tgot.MsgSignature
	query   string
	style   balaboba.Style
	user    int64
	timeout time.Time

	queue *queue.Queue[request]
}

func (request) Name() string {
	return "balaboba"
}

func (r *request) Handle(ctx tgot.MessageContext, q *tg.CallbackQuery) (tgot.CallbackAnswer, bool, error) {
	if q.From.ID != r.user {
		return tgot.CallbackAnswer{}, false, nil
	}

	if r.queue.Len() > 0 {
		_, err := ctx.EditText(tgot.EditText{
			Text: "request is added to the queue",
		}, emptyKeybaord)
		if err != nil {
			return tgot.CallbackAnswer{}, true, err
		}
	}

	u, _ := strconv.ParseUint(q.Data, 10, 8)
	r.style = balaboba.Style(u)
	r.queue.Push(*r)

	return tgot.CallbackAnswer{}, true, nil
}

func (r request) Timeout() time.Time {
	return r.timeout
}

func (r request) Cancel(ctx tgot.MessageContext) error {
	_, err := ctx.EditText(tgot.EditText{Text: "request timeout"}, emptyKeybaord)
	return err
}
