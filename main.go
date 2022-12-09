package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"tghwbot/modules"
	"tghwbot/modules/debug"
	"tghwbot/modules/images"
	"tghwbot/modules/random"
	"tghwbot/modules/text"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/commands"
	"github.com/karalef/tgot/logger"
	"github.com/karalef/tgot/updates"
)

var color = flag.Bool("color", false, "use colored log")
var wh = flag.Bool("webservice", false, "start as webservice")

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

	a, err := api.New(os.Getenv("TOKEN"), "", "", nil)
	if err != nil {
		log.Error("api initialization failed: %s", err.Error())
		return
	}

	var poller updates.Poller
	if *wh {
		ws, err := initWebservice()
		if err != nil {
			log.Error(err.Error())
			return
		}
		poller = ws.serv
	} else {
		poller = updates.NewLongPoller(30, 0, 0)
	}

	b, err := tgot.New(a, poller, tgot.Config{
		Logger: log,
		Commands: commands.List{
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

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		s := <-sig
		b.Stop()
		log.Info("bot stopped by %s", s.String())
	}()

	log.Info("bot started")
	err = b.Run()
	if err != nil {
		log.Error("bot stopped with error: %s", err.Error())
	}
	log.Info("exit")
}
