package bot

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"strconv"
	"tghwbot/bot/internal"
	"tghwbot/bot/tg"
)

// Error type.
type Error struct {
	Err      error
	Method   string
	Data     url.Values
	Response []byte
}

func (e *Error) Error() string {
	return e.Err.Error()
}

type params map[string]string

func (p params) build() url.Values {
	if p == nil {
		return nil
	}
	vals := url.Values{}
	for k, v := range p {
		vals.Set(k, v)
	}
	return vals
}

func (p params) add(key, value string) {
	if value != "" {
		p[key] = value
	}
}

func (p params) addInt(key string, value int) {
	if value != 0 {
		p[key] = strconv.Itoa(value)
	}
}

func (p params) addInt64(key string, value int64) {
	if value != 0 {
		p[key] = strconv.FormatInt(value, 10)
	}
}

func (p params) addBool(key string, value bool) {
	if value {
		p[key] = strconv.FormatBool(value)
	}
}

func (p params) addFloat(key string, value float64) {
	if value != 0 {
		p[key] = strconv.FormatFloat(value, 'f', 6, 64)
	}
}

func (p params) addJSON(key string, value interface{}) {
	if value == nil {
		return
	}

	b, _ := json.Marshal(value)
	p[key] = string(b)
}

func (b *Bot) performRequest(method string, p params, res interface{}) error {
	return b.performRequestContext(context.Background(), method, p, res)
}

func (b *Bot) performRequestContext(ctx context.Context, method string, p params, res interface{}) error {
	u := b.apiURL + b.token + "/" + method
	data := p.build()
	resp, err := internal.PostFormContext(ctx, b.client, u, data)
	switch err {
	case nil:
	case context.Canceled, context.DeadlineExceeded:
		return err
	default:
		return &Error{Err: err}
	}
	defer resp.Body.Close()

	var r tg.APIResponse
	raw, err := internal.DecodeJSON(resp.Body, &r)
	if err != nil {
		return &Error{
			Err:      err,
			Method:   method,
			Data:     data,
			Response: raw,
		}
	}
	if r.APIError != nil {
		return r.APIError
	}

	if res == nil {
		return nil
	}
	err = json.Unmarshal(r.Result, res)
	if err == nil {
		return nil
	}
	return &Error{
		Err:      err,
		Method:   method,
		Data:     data,
		Response: r.Result,
	}
}

func (b *Bot) downloadFile(path string) ([]byte, error) {
	resp, err := b.client.Get(tg.DefaultFileURL + b.token + "/" + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (b *Bot) getMe() (*tg.User, error) {
	var u tg.User
	return &u, b.performRequest("getMe", nil, &u)
}

func (b *Bot) getUpdates(ctx context.Context, offset, timeout int, allowed ...string) ([]tg.Update, error) {
	p := params{}
	p.addInt("offset", offset)
	//p.addInt("limit", limit)
	p.addInt("timeout", timeout)
	p.addJSON("allowed_updates", allowed)

	var updates []tg.Update
	return updates, b.performRequestContext(ctx, "getUpdates", p, &updates)
}

type commandScope struct {
	Type   string
	ChatID int64
	UserID int64
}

type commandParams struct {
	Commands []tg.Command
	Scope    *commandScope
	Lang     string
}

func (p *commandParams) params() params {
	if p == nil {
		return nil
	}
	v := params{}
	v.add("language_code", p.Lang)
	v.addJSON("scope", p.Scope)
	v.addJSON("commands", p.Commands)
	return v
}

func (b *Bot) getCommands(s *commandScope, lang string) ([]tg.Command, error) {
	p := commandParams{
		Scope: s,
		Lang:  lang,
	}
	var cmds []tg.Command
	return cmds, b.performRequest("getMyCommands", p.params(), &cmds)
}

func (b *Bot) setCommands(p *commandParams) error {
	return b.performRequest("setMyCommands", p.params(), nil)
}

func (b *Bot) deleteCommands(s *commandScope, lang string) error {
	p := commandParams{
		Scope: s,
		Lang:  lang,
	}
	return b.performRequest("deleteMyCommands", p.params(), nil)
}
