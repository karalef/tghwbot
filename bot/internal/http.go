package internal

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// PostFormContext issues a POST to the specified URL, with data's keys and
// values URL-encoded as the request body.
// Is also unwraps errors of context.
func PostFormContext(ctx context.Context, c *http.Client, url string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Do(req)
	return resp, unwrapCtxErr(err)
}

func unwrapCtxErr(err error) error {
	switch {
	case errors.Is(err, context.Canceled):
		return context.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return context.DeadlineExceeded
	}
	return err
}
