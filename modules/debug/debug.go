package debug

import (
	"encoding/json"
	"fmt"
	"runtime"
	"tghwbot/common/format"
	"time"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/commands"
)

// CMD is a debug command.
var CMD = commands.SimpleCommand{
	Command: "debug",
	Desc:    "debug",
	Args: []commands.Arg{
		{
			Consts: []string{"obj", "mem"},
		},
		{
			Consts: []string{"gc"},
		},
	},
	Func: func(ctx tgot.ChatContext, msg *tg.Message, args []string) error {
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
		return ctx.ReplyE(msg.ID, tgot.NewMessage(out))
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
	if gc {
		runtime.GC()
	}

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	var str string
	if gc {
		str = fmt.Sprintf("\n\nGarbage collection is done in %dns", ms.PauseNs[(ms.NumGC+255)%256])
	}
	str = fmt.Sprint(
		"Allocated: ", format.FmtBytes(ms.Alloc, false, false),
		"\nHeap objects: ", ms.Mallocs-ms.Frees,
		"\nHeap inuse: ", format.FmtBytes(ms.HeapInuse, false, false),
		"\nGoroutines: ", runtime.NumGoroutine(),
		str)
	return str
}
