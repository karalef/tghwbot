package debug

import (
	"encoding/json"
	"fmt"
	"runtime"
	"tghwbot/bot"
	"tghwbot/common/format"
	"tghwbot/common/rt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var DebugCmd = bot.Command{
	Cmd:         "debug",
	Description: "debug",
	Args: []bot.Arg{
		{
			Consts: []string{"obj", "stat", "gc"},
		},
	},
	Run: func(ctx *bot.Context, msg *tgbotapi.Message, args []string) {
		var out string
		if len(args) == 0 {
			out = memStats(false)
		} else {
			switch args[0] {
			case "obj":
				out = marshalObj(msg)
			case "stat":
				out = stat()
			case "gc":
				out = memStats(true)
			}
		}
		ctx.ReplyText(out)
	},
}

func marshalObj(obj interface{}) string {
	o, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		return err.Error()
	}
	return string(o)
}

var start = time.Now()

func stat() string {
	t := start.Format("02 Jan 2006 15:04:05")
	uptime := time.Now().Sub(start).Truncate(time.Second)
	return fmt.Sprintf("Stats\n\nStart time: %s\nUptime:%s", t, uptime)
}

func memStats(gc bool) string {
	ms := rt.GetMemStats(gc)

	var str string
	if gc {
		str = fmt.Sprintf("\n\nGarbage collection is done in %dns\nTotal stop-the-world time %dns", ms.GCTime, ms.PauseTotal)
	}
	str = fmt.Sprint(
		"Allocated: ", format.FmtBytes(ms.Allocated, false, false),
		"\nHeap objects: ", ms.Objects,
		"\nHeap inuse: ", format.FmtBytes(ms.InUse, false, false),
		"\nGoroutines: ", runtime.NumGoroutine(),
		str)
	return str
}
