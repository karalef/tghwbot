package text

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"tghwbot/modules"
	"time"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/commands"
)

var textgenMut sync.Mutex

var Gen = commands.Command{
	Cmd:         "textgen",
	Description: "text generation",
	Func: func(ctx tgot.ChatContext, msg *tg.Message, args []string) error {
		query := strings.Join(args, " ")
		if query == "" {
			return modules.ReplyText(ctx, msg, "Think of the beginning of the story")
		}

		textgenMut.Lock()
		defer textgenMut.Unlock()
		ctx.SendChatAction(tg.ActionTyping)
		replies, err := porfirevich(query, 30)
		if err != nil {
			ctx.Logger().Error(err.Error())
			return modules.ReplyText(ctx, msg, err.Error())
		}

		var text string
		for _, r := range replies {
			text += query + r + "\n\n"
		}
		return modules.ReplyText(ctx, msg, text)
	},
}

var porfirevichClient = &http.Client{
	Timeout: time.Second * 15,
}

func porfirevich(start string, length int) ([]string, error) {
	params := map[string]interface{}{
		"length": length,
		"prompt": start,
	}
	body, _ := json.Marshal(params)
	resp, err := porfirevichClient.Post("https://pelevin.gpt.dobro.ai/generate/", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("service is unavailable")
	}
	defer resp.Body.Close()

	var replies map[string][]string
	err = json.NewDecoder(resp.Body).Decode(&replies)
	if err != nil {
		return nil, err
	}
	return replies["replies"], nil
}
