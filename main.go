package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"tghwbot/modules"
	"tghwbot/modules/debug"
	"tghwbot/modules/images"
	"tghwbot/modules/random"
	"tghwbot/modules/text"
	"time"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/logger"
)

var color = flag.Bool("color", false, "use colored log")

func init() {
	loc, err := time.LoadLocation(os.Getenv("LOG_TIME_ZONE"))
	if err != nil {
		panic(err)
	}
	time.Local = loc
	flag.Parse()
}

func main() {
	log := logger.Default("HwBot", *color)
	log.Info("PID: %d", os.Getpid())

	b, err := tgot.New(os.Getenv("TOKEN"), tgot.Config{
		Logger:   log,
		MakeHelp: true,
		Commands: []*tgot.Command{
			&debug.DebugCmd,
			&random.Info,
			&random.Number,
			&random.Roll,
			&random.When,
			&text.Gen,
			&text.BalabobaCmd,
			&images.CitgenCmd,
			&images.Search,
			&images.CraiyonCmd,
		},
		Handler: tgot.Handler{
			OnInlineQuery:   images.OnInline,
			OnCallbackQuery: modules.CallbackRouter.Route,
		},
	})
	if err != nil {
		panic(err)
	}

	modules.API = b.MakeContext("modules")
	text.InitBalaboba()
	images.InitCraiyon()

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	err = b.RunContext(ctx)
	if err != nil {
		log.Error("bot finished with an error: %s", err)
		return
	}

	select {
	case <-ctx.Done():
		log.Info("stopped by os signal")
	default:
		log.Info("stopping without error")
	}
}
