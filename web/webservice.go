package web

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/karalef/tgot"
)

// New creates new web service.
func New(b *tgot.Bot) *Service {
	return &Service{
		app:    fiber.New(),
		bot:    b,
		secret: generateSecret(0),
	}
}

// Service is a web service for telegram bot.
type Service struct {
	app    *fiber.App
	bot    *tgot.Bot
	secret string
}

// Start starts web service.
func (ws *Service) Start(addr, webhookURL string) error {
	if webhookURL == "" {
		return errors.New("webhookURL is required")
	}
	ok, err := ws.bot.SetWebhook(tgot.WebhookData{
		URL:            webhookURL,
		AllowedUpdates: ws.bot.Allowed(),
		DropPending:    true,
		SecretToken:    ws.secret,
	})
	if !ok {
		return err
	}

	ws.app.Post("/tgwebhook", ws.OnWebhook)
	ws.app.All("/health", ws.health)
	ws.app.All("/stop", ws.stop)
	if addr == "" {
		addr = ":8080"
	}
	return ws.app.Listen(addr)
}

// Stop stops web service.
func (ws *Service) Stop() {
	ws.app.Shutdown()
}

// Secret returns secret token for webhook.
func (ws *Service) Secret() string {
	return ws.secret
}

func (ws *Service) verifySecret(ctx *fiber.Ctx) error {
	if ws.secret != "" &&
		ws.secret != ctx.Get("X-Telegram-Bot-Api-Secret-Token") {
		ctx.Status(fiber.StatusForbidden)
		return errors.New("forbidden " + ctx.Path() + " request from " + ctx.IP())
	}
	return nil
}

func (ws *Service) health(ctx *fiber.Ctx) error {
	ctx.Status(http.StatusOK)
	return nil
}

func (ws *Service) stop(ctx *fiber.Ctx) error {
	if err := ws.verifySecret(ctx); err != nil {
		return err
	}
	ctx.Status(fiber.StatusOK)
	ctx.WriteString("server is shutting down...")
	go ws.Stop()
	return nil
}
