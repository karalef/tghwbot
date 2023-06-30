package common

import (
	"github.com/karalef/tgot"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Log[Ctx tgot.BaseContext[Ctx]](ctx Ctx) zerolog.Logger {
	return log.With().Str("context", ctx.Name()).Logger()
}
