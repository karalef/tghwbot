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

type params url.Values

func (p params) set(key string, value interface{}) params {
	if value == nil {
		return p
	}
	vals := url.Values(p)
	switch v := value.(type) {
	case string:
		if v != "" {
			vals.Set(key, v)
		}
	case int:
		if v != 0 {
			vals.Set(key, strconv.Itoa(v))
		}
	case int64:
		if v != 0 {
			vals.Set(key, strconv.FormatInt(v, 10))
		}
	case bool:
		if v {
			vals.Set(key, strconv.FormatBool(v))
		}
	case float64:
		if v != 0 {
			vals.Set(key, strconv.FormatFloat(v, 'f', 6, 64))
		}
	default:
		b, _ := json.Marshal(value)
		vals.Set(key, string(b))
	}
	return p
}

func performRequest[T any](b *Bot, method string, p params) (T, error) {
	return performRequestContext[T](context.Background(), b, method, p)
}

func performRequestEmpty(b *Bot, method string, p params) error {
	_, err := performRequest[internal.Empty](b, method, p)
	return err
}

func performRequestContext[T any](ctx context.Context, b *Bot, method string, p params) (T, error) {
	u := b.apiURL + b.token + "/" + method
	data := url.Values(p)
	resp, err := internal.PostFormContext(ctx, b.client, u, data)

	var nilResult T
	switch err {
	case nil:
	case context.Canceled, context.DeadlineExceeded:
		return nilResult, err
	default:
		return nilResult, &Error{Err: err}
	}
	defer resp.Body.Close()

	r, raw, err := internal.DecodeJSON[tg.APIResponse[T]](resp.Body)
	if err != nil {
		return nilResult, &Error{
			Err:      err,
			Method:   method,
			Data:     data,
			Response: raw,
		}
	}
	if r.APIError != nil {
		return nilResult, r.APIError
	}

	return r.Result, nil
}

func (b *Bot) downloadFile(path string) ([]byte, error) {
	resp, err := b.client.Get(b.fileURL + b.token + "/" + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (b *Bot) getMe() (*tg.User, error) {
	return performRequest[*tg.User](b, "getMe", nil)
}

func (b *Bot) getUpdates(ctx context.Context, offset, timeout int, allowed ...string) ([]tg.Update, error) {
	p := params{}
	p.set("offset", offset)
	//p.set("limit", limit)
	p.set("timeout", timeout)
	p.set("allowed_updates", allowed)

	return performRequestContext[[]tg.Update](ctx, b, "getUpdates", p)
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
	v.set("language_code", p.Lang)
	v.set("scope", p.Scope)
	v.set("commands", p.Commands)
	return v
}

func (b *Bot) getCommands(s *commandScope, lang string) ([]tg.Command, error) {
	p := commandParams{
		Scope: s,
		Lang:  lang,
	}
	return performRequest[[]tg.Command](b, "getMyCommands", p.params())
}

func (b *Bot) setCommands(p *commandParams) error {
	return performRequestEmpty(b, "setMyCommands", p.params())
}

func (b *Bot) deleteCommands(s *commandScope, lang string) error {
	p := commandParams{
		Scope: s,
		Lang:  lang,
	}
	return performRequestEmpty(b, "deleteMyCommands", p.params())
}
