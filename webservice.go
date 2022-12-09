package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"os"

	"github.com/karalef/tgot/updates"
)

type webservice struct {
	serv   *updates.Server
	secret string
}

func initWebservice() (*webservice, error) {
	url := os.Getenv("WEBHOOK_URL")
	if url == "" {
		return nil, errors.New("WEBHOOK_URL is required for starting as webservice")
	}

	addr := ":8080"
	if port, ok := os.LookupEnv("PORT"); ok {
		addr = ":" + port
	}

	var secret [32]byte
	_, err := rand.Read(secret[:])
	if err != nil {
		return nil, errors.New("rand reader: " + err.Error())
	}

	ws := &webservice{
		secret: hex.EncodeToString(secret[:]),
	}
	ws.serv = updates.NewWebhookServer(addr, updates.ServerConfig{
		Path:        "/tgWebhook",
		URL:         url,
		SecretToken: ws.secret,
	})
	ws.serv.Mux.HandleFunc("/health", ws.health)
	ws.serv.Mux.HandleFunc("/stop", ws.stop)

	return ws, nil
}

func (ws *webservice) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (ws *webservice) stop(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	secret := query.Get("secret")
	if secret == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if secret != ws.secret {
		w.WriteHeader(http.StatusForbidden)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("server is shutting down..."))
	go ws.serv.Close()
}
