package internal

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// PostFormContext issues a POST to the specified URL, with data's keys and
// values URL-encoded as the request body.
func PostFormContext(ctx context.Context, c *http.Client, url string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.Do(req)
}
