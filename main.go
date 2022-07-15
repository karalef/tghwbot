package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tghwbot/bot"
	"tghwbot/modules/debug"
	"tghwbot/modules/images"
	"tghwbot/modules/random"
	"tghwbot/modules/text"
	"tghwbot/modules/wikihow"

	"github.com/Toffee-iZt/wfs"
)

func main() {
	println(wfs.ExecPath())
	println("PID", os.Getpid())

	cmds := []*bot.Command{
		&debug.DebugCmd,
		&random.Flip,
		&random.Info,
		&random.Number,
		&random.When,
		&text.Gen,
		&images.CitgenCmd,
		&wikihow.Wikihow,
	}
	b, err := bot.New(os.Getenv("TOKEN"), cmds...)
	if err != nil {
		log.Panicf("tg auth: %s", err.Error())
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	err = b.Run(ctx)

	switch err {
	case nil:
		log.Print("stopping without error")
	case context.Canceled:
		log.Print("stopping by os signal")
	default:
		log.Fatalf("bot finished with an error: %s", err)
	}
}
