package main

import (
	"context"
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

func main() {
	println("PID", os.Getpid())

	log := logger.New(logger.DefaultWriter, "HwBot")

	b, err := bot.New(bot.Config{
		Token:    os.Getenv("TOKEN"),
		Logger:   log,
		MakeHelp: true,
		Commands: []*bot.Command{
			&debug.DebugCmd,
			&random.Flip,
			&random.Info,
			&random.Number,
			&random.When,
			&text.Gen,
			&images.CitgenCmd,
		},
	})
	if err != nil {
		panic(err)
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	err = b.Run(ctx, 0)

	switch err {
	case nil:
		log.Info("stopping without error")
	case context.Canceled:
		log.Info("stopping by os signal")
	default:
		log.Error("bot finished with an error: %s", err)
	}
}
