package porfirevich

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"tghwbot/common"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/commands"
)

var mut sync.Mutex

var client = Client{http: &http.Client{Timeout: 15 * time.Second}}

var CMD = commands.SimpleCommand{
	Command: "porfirevich",
	Desc:    "text continuation",
	Func: func(ctx *tgot.Message, msg *tg.Message, args []string) error {
		prompt := strings.Join(args, " ")
		if prompt == "" {
			return ctx.ReplyText("Think of the beginning of the story")
		}
		logger := common.Log(ctx)

		mut.Lock()
		defer mut.Unlock()
		ctx.Chat().SendChatAction(tg.ActionTyping)
		replies, err := client.Generate(Request{
			Prompt:      prompt,
			Length:      150,
			Temperature: 0.3,
			Model:       ModelLawa,
		})
		if err != nil {
			logger.Err(err).Msg("text generation failed")
			return ctx.ReplyText(err.Error())
		}

		var text string
		for _, r := range replies {
			text += prompt + r + "\n\n"
		}
		return ctx.ReplyText(text)
	},
}
