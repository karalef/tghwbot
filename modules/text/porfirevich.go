package text

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var porfirevichClient = &http.Client{
	Timeout: time.Second * 10,
}

type porfirevichResult struct {
	Replies []string
	Error   error
}

func porfirevichAsync(start string, length int) <-chan *porfirevichResult {
	ch := make(chan *porfirevichResult)
	go func() {
		r, e := porfirevich(start, length)
		ch <- &porfirevichResult{
			Replies: r,
			Error:   e,
		}
	}()
	return ch
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
