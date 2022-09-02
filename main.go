package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"tghwbot/bot"
	"tghwbot/bot/logger"
	"tghwbot/modules/debug"
	"tghwbot/modules/images"
	"tghwbot/modules/random"
	"tghwbot/modules/text"
	"time"
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

	b, err := bot.New(os.Getenv("TOKEN"), bot.Config{
		Logger:   log,
		MakeHelp: true,
		Commands: []*bot.Command{
			&debug.DebugCmd,
			&random.Info,
			&random.Number,
			&random.Roll,
			&random.When,
			&text.Gen,
			&images.CitgenCmd,
			&images.Search,
		},
		Handler: bot.Handler{
			OnInlineQuery: images.OnInline,
		},
	})
	if err != nil {
		panic(err)
	}

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
