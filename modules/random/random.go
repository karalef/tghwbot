package random

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"tghwbot/bot"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Rand(max int) int {
	return randInt(max)
}

func RandP(max int, power float64) int {
	return randP(max, power)
}

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
	Run: func(ctx *bot.Context, msg *tgbotapi.Message, args []string) {
		var max int64 = 100
		if len(args) > 0 {
			num := args[0]
			var err error
			max, err = strconv.ParseInt(num, 10, 64)
			if err != nil || max <= 0 {
				ctx.ReplyText("Specify a number between 1 and MaxInt64")
			}
		}
		ctx.ReplyText(strconv.FormatInt(myRand.Int63n(max), 10))
	},
}

var Flip = bot.Command{
	Cmd:         "flip",
	Description: "flip a coin",
	Run: func(ctx *bot.Context, msg *tgbotapi.Message, args []string) {
		r := "Tails"
		if myRand.Intn(2) == 1 {
			r = "Heads"
		}
		ctx.ReplyText(r)
	},
}

var Info = bot.Command{
	Cmd:         "info",
	Description: "random event probability",
	Run: func(ctx *bot.Context, msg *tgbotapi.Message, args []string) {
		if len(args) == 0 {
			ctx.ReplyText("Specify the event")
		}
		p := myRand.Intn(101)
		e := strings.Join(args, " ")
		ctx.ReplyText("The probability that " + e + " — " + strconv.Itoa(p) + "%")
	},
}

var When = bot.Command{
	Cmd:         "when",
	Description: "random date of the event",
	Run: func(ctx *bot.Context, msg *tgbotapi.Message, args []string) {
		if len(args) == 0 {
			ctx.ReplyText("Specify the event")
		}
		t := time.Now().AddDate(randP(51, 1.5), randInt(12), randInt(31))
		e := strings.Join(args, " ")
		ctx.ReplyText(e + " " + t.Format("02 Jan 2006"))
	},
}
