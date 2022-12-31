package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"tghwbot/modules/debug"
	"tghwbot/modules/images"
	"tghwbot/modules/random"
	"tghwbot/modules/text"
	"tghwbot/modules/text/balaboba"
	"tghwbot/web"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/commands"
	"github.com/karalef/tgot/logger"
	"github.com/karalef/tgot/router"
	"github.com/karalef/tgot/updates"
)

var color = flag.Bool("color", false, "use colored log")

func init() {
	flag.Parse()
}

func main() {
	var colorConf logger.ColorConfig
	if *color {
		colorConf = logger.DefaultColorConfig
	}
	log := logger.Default("HwBot", colorConf)
	log.Info("starting bot (PID: %d)", os.Getpid())

	b, err := tgot.NewWithToken(os.Getenv("TOKEN"), log)
	if err != nil {
		log.Error("bot initialization failed: %s", err.Error())
		return
	}

	modsCtx := b.MakeContext("modules")
	cbRouter := router.NewCallbacks()

	var cmds commands.List
	cmds = commands.List{
		commands.MakeHelp(&cmds),
		&debug.DebugCmd,
		&random.Info,
		&random.Number,
		&random.When,
		&text.Gen,
		balaboba.Command(modsCtx, cbRouter),
		&images.CitgenCmd,
		&images.Search,
		images.CraiyonCommand(modsCtx),
	}
	cmds.Setup(b)

	b.OnInlineQuery = images.OnInline
	b.OnCallbackQuery = cbRouter.Route
	b.OnMessage = (&commands.MessageHandler{
		Username: b.Me().Username,
		Command:  cmds.Command,
	}).Handle

	run(b, log)
	log.Info("exit")
}

func run(b *tgot.Bot, log *logger.Logger) {
	var start func() error
	var cancel func()
	if u, ok := os.LookupEnv("WEBHOOK_URL"); ok {
		ws, err := web.NewWebservice(b)
		if err != nil {
			log.Error("failed to start webservice: %s", err.Error())
			return
		}
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
		log.Info("bot stopped by %s", s.String())
	}()

	log.Info("bot started")
	err := start()
	if err != nil {
		log.Error("bot stopped with error: %s", err.Error())
	}
}
