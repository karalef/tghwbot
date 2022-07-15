package modules

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

var rapidAPIKey = os.Getenv("RAPID_API_KEY")
var rapidAPIClient = http.Client{}

// IsRapidAPIReady returns true if key is provided.
func IsRapidAPIReady() bool {
	return rapidAPIKey != ""
}

// RapidAPIRequest type.
type RapidAPIRequest struct {
	Method string
	Host   string
	Path   string
	Query  url.Values
	Data   io.Reader
}

// RapidAPI sends http request to RapidAPI.
func RapidAPI[T any](ctx context.Context, r RapidAPIRequest) (*T, error) {
	if !IsRapidAPIReady() {
		return nil, errors.New("no RapidAPI key provided")
	}
	if r.Method == "" || r.Host == "" || r.Path == "" {
		return nil, errors.New("invalid RapidAPI request")
	}
	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), time.Second*15)
	}
	u := "https://" + r.Host + r.Path + "?" + r.Query.Encode()
	req, err := http.NewRequestWithContext(ctx, r.Method, u, r.Data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-RapidAPI-Host", r.Host)
	req.Header.Set("X-RapidAPI-Key", rapidAPIKey)
	resp, err := rapidAPIClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result T
	return &result, json.NewDecoder(resp.Body).Decode(&result)
}
