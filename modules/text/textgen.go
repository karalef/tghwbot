package text

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"tghwbot/common"
	"time"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/commands"
)

var textgenMut sync.Mutex

var Gen = commands.SimpleCommand{
	Command: "textgen",
	Desc:    "text generation",
	Func: func(ctx *tgot.Message, msg *tg.Message, args []string) error {
		logger := common.Log(ctx)
		query := strings.Join(args, " ")
		if query == "" {
			return ctx.ReplyText("Think of the beginning of the story")
		}

		textgenMut.Lock()
		defer textgenMut.Unlock()
		ctx.Chat().SendChatAction(tg.ActionTyping)
		replies, err := porfirevich(query, 30)
		if err != nil {
			logger.Err(err).Msg("text generation failed")
			return ctx.ReplyText(err.Error())
		}

		var text string
		for _, r := range replies {
			text += query + r + "\n\n"
		}
		return ctx.ReplyText(text)
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
