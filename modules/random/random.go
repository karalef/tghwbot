package random

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"tghwbot/bot"
	"time"

	"gopkg.in/telebot.v3"
)

var myRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randInt(max int) int {
	return myRand.Intn(max)
}

func randP(max int, power float64) int {
	r := myRand.Float64()
	return int(math.Floor(float64(max+1) * math.Pow(r, power)))
}

var Number = bot.Command{
	Cmd:         "rand",
	Description: "random number",
	Run: func(ctx *bot.Context, msg *telebot.Message, args []string) {
		var max int64 = 100
		if len(args) > 0 {
			num := args[0]
			var err error
			max, err = strconv.ParseInt(num, 10, 64)
			if err != nil || max <= 0 {
				ctx.ReplyClose("Укажите число от 1 до MaxInt64")
			}
		}
		ctx.Reply("Выпало число " + strconv.FormatInt(myRand.Int63n(max), 10))
	},
}

var Flip = bot.Command{
	Cmd:         "flip",
	Description: "flip a coin",
	Run: func(ctx *bot.Context, msg *telebot.Message, args []string) {
		r := "Выпала решка"
		if myRand.Intn(2) == 1 {
			r = "Выпал орёл"
		}
		ctx.Reply(r)
	},
}

var Info = bot.Command{
	Cmd:         "info",
	Description: "event probability",
	Run: func(ctx *bot.Context, msg *telebot.Message, args []string) {
		if len(args) == 0 {
			ctx.ReplyClose("Укажите событие")
		}
		p := myRand.Intn(101)
		e := strings.Join(args, " ")
		ctx.Reply("Вероятность того, что " + e + " — " + strconv.Itoa(p) + "%")
	},
}

var When = bot.Command{
	Cmd:         "when",
	Description: "Когда произойдет событие",
	Run: func(ctx *bot.Context, msg *telebot.Message, args []string) {
		if len(args) == 0 {
			ctx.ReplyClose("Укажите событие")
		}
		t := time.Now().AddDate(randP(51, 1.5), randInt(12), randInt(31))
		e := strings.Join(args, " ")
		ctx.Reply(e + " " + t.Format("02 Jan 2006"))
	},
}
