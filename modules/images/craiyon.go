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
	"tghwbot/queue"
	"time"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/commands"
)

func CraiyonCommand(modulesCtx tgot.Context) *commands.Command {
	d := net.Dialer{
		Timeout: 2 * time.Minute,
	}
	craiyon := &Craiyon{
		client: http.Client{
			Timeout: d.Timeout,
			Transport: &http.Transport{
				DialContext:         d.DialContext,
				TLSHandshakeTimeout: d.Timeout,
			},
		},
		api: modulesCtx.Child("craiyon"),
	}
	craiyon.queue = queue.New(craiyon.complete, 5)

	return &commands.Command{
		Cmd:         "dalle",
		Description: "Craiyon, formerly DALL·E mini, is an AI model that can draw images from any text prompt!",
		Args: []commands.Arg{
			{
				Required: true,
				Name:     "prompt",
			},
		},
		Func: func(ctx tgot.ChatContext, msg *tg.Message, args []string) error {
			prompt := strings.Join(args, " ")
			if len(prompt) == 0 {
				return ctx.ReplyE(msg.ID, tgot.NewMessage("write a prompt"))
			}
			craiyon.queue.Push(craiyonRequest{
				chat:   msg.Chat.ID,
				orig:   msg.ID,
				prompt: prompt,
			})
			return ctx.ReplyE(msg.ID, tgot.NewMessage("request is added to the queue"))
		},
	}
}

type craiyonRequest struct {
	chat   int64
	orig   int
	prompt string
}

// Craiyon is a client for craiyon api.
type Craiyon struct {
	client http.Client
	queue  *queue.Queue[craiyonRequest]

	api tgot.Context
}

func (c *Craiyon) complete(req craiyonRequest) {
	chat := c.api.OpenChat(tgot.ChatID(req.chat))

	imgs, err := c.Generate(req.prompt)
	if err != nil {
		c.api.Logger().Error(err.Error())
		err = chat.ReplyE(req.orig, tgot.NewMessage("Generation error"))
		if err != nil {
			c.api.Logger().Error(err.Error())
		}
		return
	}

	mediaGroup := tgot.MediaGroup{
		Media: make([]tg.MediaInputter, len(imgs)),
	}
	for i := range imgs {
		mediaGroup.Media[i] = tg.NewInputMediaPhoto(tg.FileReader(strconv.Itoa(i), imgs[i]))
	}

	_, err = chat.SendMediaGroup(mediaGroup, tgot.SendOptions{
		ReplyTo: req.orig,
	})
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
