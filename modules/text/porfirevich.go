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
