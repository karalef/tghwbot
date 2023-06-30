package main

import (
	"os"
	"os/signal"
	"syscall"
	"tghwbot/modules/citgen"
	"tghwbot/modules/craiyon"
	"tghwbot/modules/debug"
	"tghwbot/modules/random"
	"tghwbot/modules/search"
	"tghwbot/modules/text"
	"tghwbot/web"
	"time"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/commands"
	"github.com/karalef/tgot/router"
	"github.com/karalef/tgot/updates"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	time.Local = time.UTC
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Printf("starting bot (PID: %d)", os.Getpid())

	b, err := tgot.NewWithToken(os.Getenv("TOKEN"))
	if err != nil {
		log.Error().Err(err).Msg("bot initialization failed")
		return
	}

	modsCtx := b.MakeContext("modules")
	cbRouter := router.NewCallbacks()

	var cmds commands.List
	cmds = commands.List{
		commands.MakeHelp(&cmds),
		debug.CMD,
		random.Info,
		random.Number,
		random.When,
		text.Gen,
		citgen.CMD,
		search.CMD,
		craiyon.CMD(modsCtx),
	}
	cmds.Setup(b)

	b.OnInlineQuery = search.OnInline
	b.OnCallbackQuery = cbRouter.Route
	b.OnMessage = (&commands.MessageHandler{
		Username: b.Me().Username,
		Command:  cmds.Command,
	}).Handle

	run(b)
}

func run(b *tgot.Bot) {
	var start func() error
	var cancel func()
	if u, ok := os.LookupEnv("WEBHOOK_URL"); ok {
		ws := web.New(b)
		addr := os.Getenv("PORT")
		if addr != "" {
			addr = ":" + addr
		}
		start = func() error {
			return ws.Start(addr, u)
		}
		cancel = ws.Stop
	} else {
		lp := updates.NewLongPoller(30, 0, 0)
		start = func() error {
			return lp.Run(b)
		}
		cancel = lp.Close
	}

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		s := <-sig
		cancel()
		log.Info().Msgf("bot stopped by %s signal", s.String())
	}()

	log.Info().Msg("bot started")
	err := start()
	if err != nil {
		log.Err(err).Msg("bot stopped with error")
	}
}
