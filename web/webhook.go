package web

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/karalef/tgot/api/tg"
)

// OnWebhook handles webhook request.
func (w *Webservice) OnWebhook(ctx *fiber.Ctx) error {
	if err := w.verifySecret(ctx); err != nil {
		return err
	}

	if ctx.Get("Content-Type") != "application/json" {
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}

	var upd tg.Update
	if err := ctx.BodyParser(&upd); err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}

	return w.bot.Handle(&upd)
}

// length 1-256 (default: 128)
func generateSecret(length uint) (string, error) {
	if length > 256 {
		return "", errors.New("secret length must be in range 1-256")
	}
	if length < 1 {
		length = 128
	}
	buf := bytes.NewBuffer(make([]byte, 0, length+length%2))
	_, err := io.CopyN(hex.NewEncoder(buf), rand.Reader, int64(length/2+length%2))
	if err != nil {
		return "", err
	}
	return buf.String()[:length], nil
}
