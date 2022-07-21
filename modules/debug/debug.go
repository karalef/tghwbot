package debug

import (
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
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
			Consts: []string{"info", "obj", "mem"},
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
			case "info":
				out = info()
			case "obj":
				out = marshalObj(msg)
			case "mem":
				out = memStats(len(args) > 1 && args[1] == "gc")
			}
		}
		ctx.Reply(out)
	},
}

func info() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return runtime.Version()
	}
	var revision, time, modified string
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
			if len(revision) > 8 {
				revision = revision[:8]
			}
		case "vcs.time":
			time = s.Value
		case "vcs.modified":
			modified = s.Value
		}
	}
	var buf strings.Builder
	buf.WriteString("Go version: " + bi.GoVersion + "\n\n")
	buf.WriteString("Commit: " + revision + "\n")
	buf.WriteString("Commit time: " + time + "\n")
	buf.WriteString("Have uncommited changes: " + modified)
	return buf.String()
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
