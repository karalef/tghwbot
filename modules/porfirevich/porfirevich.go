package porfirevich

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// Model is a text generation model.
type Model uint8

// models.
const (
	ModelLawa Model = iota
	ModelMig
	ModelXLarge
	ModelGPT3
	ModelFrida
)

func (m Model) String() string {
	switch m {
	case ModelLawa:
		return "lawa"
	case ModelMig:
		return "mig"
	case ModelXLarge:
		return "xlarge"
	case ModelGPT3:
		return "gpt3"
	case ModelFrida:
		return "frida"
	default:
		return ""
	}
}

func (m Model) MarshalJSON() ([]byte, error) {
	str := m.String()
	if str == "" {
		return nil, errors.New("invalid model")
	}
	return []byte(`"` + str + `"`), nil
}

// Request is a request for text generation.
type Request struct {
	Prompt      string  `json:"prompt"`
	Length      uint    `json:"length"`
	Model       Model   `json:"model"`
	Temperature float64 `json:"temperature"`
}

func (r Request) Validate() error {
	if r.Prompt == "" {
		return errors.New("prompt is empty")
	}
	if r.Length == 0 {
		return errors.New("length must be non-zero")
	}
	if r.Model.String() == "" {
		return errors.New("invalid model")
	}
	if r.Temperature < 0 || r.Temperature > 10 {
		return errors.New("invalid temperature")
	}
	return nil
}

type Client struct {
	http *http.Client
}

func (c Client) Generate(req Request) ([]string, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Post("https://api.porfirevich.com/generate/", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("service is unavailable")
	}
	defer resp.Body.Close()

	var response struct {
		Replies []string `json:"replies"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response.Replies, nil
}
