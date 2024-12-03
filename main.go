package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"tghwbot/modules/citgen"
	"tghwbot/modules/debug"
	"tghwbot/modules/porfirevich"
	"tghwbot/modules/random"
	"tghwbot/modules/search"
	"tghwbot/modules/text"
	"tghwbot/web"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/commands"
	"github.com/karalef/tgot/handler"
	"github.com/karalef/tgot/router"
	"github.com/karalef/tgot/updates"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	time.Local = time.UTC
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	log.Printf("starting bot (PID: %d)", os.Getpid())

	a, err := api.NewDefault(os.Getenv("TOKEN"))
	if err != nil {
		log.Error().Err(err).Msg("api initialization failed")
		return
	}

	cbRouter := router.NewCallbacks()
	var cmds commands.List
	cmds = commands.List{
		commands.MakeHelp(&cmds),
		debug.CMD,
		random.Info,
		random.Number,
		random.When,
		text.Gen,
		porfirevich.CMD,
		citgen.CMD,
		search.CMD,
	}
	filter := &commands.Filter{
		Command: cmds.Command,
	}

	b, err := tgot.New(a, &handler.Handler{
		OnInlineQuery:   search.OnInline,
		OnCallbackQuery: cbRouter.Route,
		OnMessage:       filter.Handle,
	})
	if err != nil {
		log.Error().Err(err).Msg("bot initialization failed")
		return
	}

	filter.Username = b.Me().Username
	err = cmds.Setup(b)
	if err != nil {
		log.Error().Err(err).Msg("commands setup failed")
		return
	}

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
		sig := make(chan os.Signal, 1)
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
