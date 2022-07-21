package bot

import (
	"tghwbot/bot/internal"
	"tghwbot/bot/tg"
)

// InlineContext type.
type InlineContext struct {
	bot           *Bot
	inlineQueryID string
}

func (c *InlineContext) getBot() *Bot {
	return c.bot
}

func (c *InlineContext) caller() string {
	return "bot::InlineHandler"
}

// InlineAnswer represents answer to inline query.
type InlineAnswer struct {
	Results           []tg.InlineQueryResult
	CacheTime         int
	IsPersonal        bool
	NextOffset        string
	SwitchPMText      string
	SwitchPMParameter string
}

// Answer answers to inline query.
func (c *InlineContext) Answer(answer *InlineAnswer) {
	p := params{}.set("inline_query_id", c.inlineQueryID)
	p.set("result", answer.Results)
	p.set("cache_time", answer.CacheTime)
	p.set("is_personal", answer.IsPersonal)
	p.set("next_offset", answer.NextOffset)
	p.set("switch_pm_text", answer.SwitchPMText)
	p.set("switch_pm_parameter", answer.SwitchPMParameter)
	api[internal.Empty](c, "answerInlineQuery", p)
	c.bot.closeExecution()
}
