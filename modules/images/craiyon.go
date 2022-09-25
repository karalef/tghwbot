package images

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"tghwbot/modules"
	"tghwbot/queue"
	"time"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/tg"
)

var craiyon *Craiyon

func InitCraiyon() {
	d := net.Dialer{
		Timeout: 2 * time.Minute,
	}
	craiyon = &Craiyon{
		client: http.Client{
			Timeout: d.Timeout,
			Transport: &http.Transport{
				DialContext:         d.DialContext,
				TLSHandshakeTimeout: d.Timeout,
			},
		},
		api: modules.API.Child("craiyon"),
	}
	craiyon.queue = queue.New(craiyon.complete)
}

var CraiyonCmd = tgot.Command{
	Cmd:         "dalle",
	Description: "Craiyon, formerly DALL·E mini, is an AI model that can draw images from any text prompt!",
	Args: []tgot.Arg{
		{
			Required: true,
			Name:     "prompt",
		},
	},
	Run: func(ctx tgot.MessageContext, msg *tg.Message, args []string) error {
		prompt := strings.Join(args, " ")
		if len(prompt) == 0 {
			return ctx.ReplyText("write a prompt")
		}
		sent, err := ctx.Chat.Send(tgot.NewMessage("request is added to the queue"), tgot.SendOptions[tg.ReplyMarkup]{
			BaseSendOptions: tgot.BaseSendOptions{
				ReplyTo: msg.ID,
			},
		})
		if err == nil {
			craiyon.queue.Push(craiyonRequest{
				sig:    ctx.MessageSignature(sent),
				chat:   msg.Chat.ID,
				orig:   msg.ID,
				sent:   sent.ID,
				prompt: prompt,
			})
		}
		return err
	},
}

type craiyonRequest struct {
	sig    tgot.MessageSignature
	chat   int64
	orig   int
	sent   int
	prompt string
}

// Craiyon is a client for craiyon api.
type Craiyon struct {
	client http.Client
	mut    sync.Mutex
	queue  *queue.Queue[craiyonRequest]

	api tgot.Context
}

func (c *Craiyon) complete(req craiyonRequest) {
	c.mut.Lock()
	defer c.mut.Unlock()
	_, err := c.api.EditText(req.sig, tgot.EditText{Text: "wait up to 2 minutes..."})
	if err != nil {
		c.api.Logger().Error(err.Error())
	}

	chat := c.api.OpenChat(req.chat)

	imgs, err := c.Generate(req.prompt)
	if err != nil {
		c.api.Logger().Error(err.Error())
		err = nil
		chat.Send(tgot.NewMessage("Generation error"), tgot.SendOptions[tg.ReplyMarkup]{
			BaseSendOptions: tgot.BaseSendOptions{
				ReplyTo: req.orig,
			},
		})
		if err != nil {
			c.api.Logger().Error(err.Error())
		}
	}

	mediaGroup := make(tgot.MediaGroup, len(imgs))
	for i := range imgs {
		mediaGroup[i] = tg.NewInputMediaPhoto(tg.FileReader(strconv.Itoa(i), imgs[i]))
	}

	_, err = chat.SendMediaGroup(mediaGroup, tgot.MediaGroupSendOptions{
		ReplyTo: req.orig,
	})
	if err != nil {
		c.api.Logger().Error(err.Error())
	}

	err = chat.DeleteMessage(req.sent)
	if err != nil {
		c.api.Logger().Error(err.Error())
	}
}

func (c *Craiyon) request(ctx context.Context, prompt string) ([]string, error) {
	method := http.MethodOptions
	var body io.Reader

	if prompt != "" {
		body = strings.NewReader("{\"prompt\":\"" + prompt + "\"}")
		method = http.MethodPost
	}

	req, err := http.NewRequestWithContext(ctx, method, "https://backend.craiyon.com/generate", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status %s (%d)", resp.Status, resp.StatusCode)
	}
	if prompt == "" {
		return nil, nil
	}

	var r struct {
		Images  []string `json:"images"`
		Version string   `json:"version"`
	}
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&r); err != nil {
		raw, _ := io.ReadAll(io.MultiReader(dec.Buffered(), resp.Body))
		return nil, fmt.Errorf("%s\nresponse: %s", err.Error(), string(raw))
	}
	return r.Images, err
}

// IsAvailable checks the api for availability.
func (c *Craiyon) IsAvailable() bool {
	_, err := c.request(context.Background(), "")
	return err != nil
}

// Generate generates images from prompt.
func (c *Craiyon) Generate(prompt string) ([]io.Reader, error) {
	if prompt == "" {
		return nil, errors.New("prompt must be non-empty")
	}
	resp, err := c.request(context.Background(), prompt)
	if err != nil {
		return nil, err
	}

	imgs := make([]io.Reader, len(resp))
	for i := range resp {
		imgs[i] = base64.NewDecoder(base64.StdEncoding, strings.NewReader(resp[i]))
	}
	return imgs, nil
}
