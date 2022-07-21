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

// RandInt returns random int in range 0-max.
func RandInt(max int) int {
	return myRand.Intn(max)
}

// RandP returns random int in range 0-max with probability
// controlled by power.
func RandP(max int, power float64) int {
	r := myRand.Float64()
	return int(math.Floor(float64(max+1) * math.Pow(r, power)))
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
			if strings.IndexByte(num, '-') > 0 {
				s := strings.SplitN(num, "-", 2)
				offset, err = strconv.ParseInt(s[0], 10, 64)
				if err != nil || offset < 0 {
					ctx.Reply("Укажите числа от 1 до MaxInt64")
				}
				num = s[2]
			}
			max, err = strconv.ParseInt(num, 10, 64)
			if err != nil || max <= 0 {
				ctx.Reply("Укажите число от 1 до MaxInt64")
			}
		}
		max -= offset
		ctx.Reply("Выпало число " + strconv.FormatInt(offset+myRand.Int63n(max), 10))
	},
}

var Flip = bot.Command{
	Cmd:         "flip",
	Description: "flip a coin",
	Run: func(ctx *bot.Context, msg *tg.Message, args []string) {
		r := "Выпала решка"
		if myRand.Intn(2) == 1 {
			r = "Выпал орел"
		}
		ctx.Reply(r)
	},
}

var Info = bot.Command{
	Cmd:         "info",
	Description: "event probability",
	Run: func(ctx *bot.Context, msg *tg.Message, args []string) {
		if len(args) == 0 {
			ctx.Reply("Укажите событие")
		}
		p := myRand.Intn(101)
		e := strings.Join(args, " ")
		ctx.Reply("Вероятность того, что " + e + " — " + strconv.Itoa(p) + "%")
	},
}

var When = bot.Command{
	Cmd:         "when",
	Description: "Когда произойдет событие",
	Run: func(ctx *bot.Context, msg *tg.Message, args []string) {
		if len(args) == 0 {
			ctx.Reply("Укажите событие")
		}
		t := time.Now().AddDate(RandP(51, 1.5), RandInt(12), RandInt(31))
		e := strings.Join(args, " ")
		ctx.Reply(e + " " + t.Format("02 Jan 2006"))
	},
}
