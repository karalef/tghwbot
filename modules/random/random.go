package random

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"tghwbot/bot"
	"tghwbot/bot/tg"
	"time"
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

var Number = bot.Command{
	Cmd:         "rand",
	Description: "random number",
	Run: func(ctx *bot.Context, msg *tg.Message, args []string) {
		var max int64 = 100
		var offset int64 = 0
		if len(args) > 0 {
			num := args[0]
			var err error
			if i := strings.IndexByte(num, '-'); i > 0 && i < len(num)-1 {
				offset, err = strconv.ParseInt(num[:i], 10, 64)
				if err != nil || offset < 0 {
					ctx.Reply("Specify the numbers in range 0 - MaxInt64")
				}
				num = num[i+1:]
			}
			max, err = strconv.ParseInt(num, 10, 64)
			if err != nil || max <= 0 {
				ctx.Reply("Specify the number in range 1 - MaxInt64")
			}
		}
		max -= offset
		ctx.Reply(strconv.FormatInt(offset+myRand.Int63n(max), 10))
	},
}

var Info = bot.Command{
	Cmd:         "info",
	Description: "event probability",
	Run: func(ctx *bot.Context, msg *tg.Message, args []string) {
		if len(args) == 0 {
			ctx.Reply("Specify the event")
		}
		p := myRand.Intn(101)
		e := strings.Join(args, " ")
		ctx.Reply("The probability that " + e + " — " + strconv.Itoa(p) + "%")
	},
}

var When = bot.Command{
	Cmd:         "when",
	Description: "random date of event",
	Run: func(ctx *bot.Context, msg *tg.Message, args []string) {
		if len(args) == 0 {
			ctx.Reply("Provide the event")
		}
		t := time.Now().AddDate(RandP(51, 1.5), RandInt(12), RandInt(31))
		e := strings.Join(args, " ")
		ctx.Reply(e + " " + t.Format("02 Jan 2006"))
	},
}
