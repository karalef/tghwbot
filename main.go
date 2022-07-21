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

	cmds := []*bot.Command{
		&debug.DebugCmd,
		&random.Flip,
		&random.Info,
		&random.Number,
		&random.When,
		&text.Gen,
		&images.CitgenCmd,
	}
	b, err := bot.New(os.Getenv("TOKEN"), log, cmds...)
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
