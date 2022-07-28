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
)

var color = flag.Bool("color", false, "use colored log")

func init() {
	flag.Parse()
}

func main() {
	log := logger.New(logger.NewWriter(os.Stderr, *color), "HwBot")
	log.Info("PID: %d", os.Getpid())

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	b, err := bot.NewWithContext(ctx, bot.Config{
		Token:    os.Getenv("TOKEN"),
		Logger:   log,
		MakeHelp: true,
		Commands: []*bot.Command{
			&debug.DebugCmd,
			&random.Info,
			&random.Number,
			&random.When,
			&text.Gen,
			&images.CitgenCmd,
			&images.Search,
		},
	})
	if err != nil {
		panic(err)
	}

	err = b.Run(0)

	switch err {
	case nil:
		log.Info("stopping without error")
	case context.Canceled:
		log.Info("stopping by os signal")
	default:
		log.Error("bot finished with an error: %s", err)
	}
}
