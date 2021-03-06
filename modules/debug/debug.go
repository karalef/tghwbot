package debug

import (
	"encoding/json"
	"fmt"
	"runtime"
	"tghwbot/bot"
	"tghwbot/bot/tg"
	"tghwbot/common/format"
	"tghwbot/common/rt"
	"time"
)

var DebugCmd = bot.Command{
	Cmd:         "debug",
	Description: "debug",
	Args: []bot.Arg{
		{
			Consts: []string{"obj", "mem"},
		},
		{
			Consts: []string{"gc"},
		},
	},
	Run: func(ctx *bot.Context, msg *tg.Message, args []string) {
		var out string
		if len(args) == 0 {
			out = stat()
		} else {
			switch args[0] {
			case "obj":
				out = marshalObj(msg)
			case "mem":
				out = memStats(len(args) > 1 && args[1] == "gc")
			}
		}
		ctx.Reply(out)
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
	return fmt.Sprintf("Stats\n\nStart time: %s\nUptime: %s", t, uptime)
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
