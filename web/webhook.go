package web

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/karalef/tgot/api/tg"
)

// OnWebhook handles webhook request.
func (ws *Service) OnWebhook(ctx *fiber.Ctx) error {
	if err := ws.verifySecret(ctx); err != nil {
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

	return ws.bot.Handle(&upd)
}

// length 1-128 (default: 64)
// the size of the output will be twice the length
func generateSecret(length uint8) string {
	if length == 0 {
		length = 64
	}
	buf := bytes.NewBuffer(make([]byte, 0, length))
	_, err := io.CopyN(hex.NewEncoder(buf), rand.Reader, int64(length))
	if err != nil {
		panic("error while generating secret: " + err.Error())
	}
	return buf.String()[:length]
}
