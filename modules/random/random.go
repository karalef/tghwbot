package random

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/commands"
)

var myRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandInt returns random int in range [0,max).
func RandInt(max int) int {
	return myRand.Intn(max)
}

// RandP returns random int in range [0,max) with probability
// controlled by power.
func RandP(max int, power float64) int {
	r := myRand.Float64()
	return int(math.Floor(float64(max) * math.Pow(r, power)))
}

var Number = commands.SimpleCommand{
	Command: "rand",
	Desc:    "random number",
	Func: func(ctx *tgot.Message, msg *tg.Message, args []string) error {
		var max int64 = 100
		var offset int64 = 0
		if len(args) > 0 {
			num := args[0]
			var err error
			if i := strings.IndexByte(num, '-'); i > 0 && i < len(num)-1 {
				offset, err = strconv.ParseInt(num[:i], 10, 64)
				if err != nil || offset < 0 {
					return ctx.ReplyText("Specify the numbers in range 0 - MaxInt64")
				}
				num = num[i+1:]
			}
			max, err = strconv.ParseInt(num, 10, 64)
			if err != nil || max <= 0 {
				return ctx.ReplyText("Specify the numbers in range 1 - MaxInt64")
			}
		}
		max -= offset
		return ctx.ReplyText(strconv.FormatInt(offset+myRand.Int63n(max), 10))
	},
}

var Info = commands.SimpleCommand{
	Command: "info",
	Desc:    "event probability",
	Func: func(ctx *tgot.Message, msg *tg.Message, args []string) error {
		if len(args) == 0 {
			return ctx.ReplyText("Specify the event")
		}
		p := myRand.Intn(101)
		e := strings.Join(args, " ")
		return ctx.ReplyText("The probability that " + e + " — " + strconv.Itoa(p) + "%")
	},
}

var When = commands.SimpleCommand{
	Command: "when",
	Desc:    "random date of event",
	Func: func(ctx *tgot.Message, msg *tg.Message, args []string) error {
		if len(args) == 0 {
			return ctx.ReplyText("Specify the event")
		}
		t := time.Now().AddDate(RandP(51, 1.5), RandInt(12), RandInt(31))
		e := strings.Join(args, " ")
		return ctx.ReplyText(e + " " + t.Format("02 Jan 2006"))
	},
}
