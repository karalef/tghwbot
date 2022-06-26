package bot

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// New creates new bot.
func New(token string, cmds ...*Command) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	b := &Bot{
		api:  api,
		cmds: cmds,
	}
	b.cmds = append(b.cmds, &ping, makeHelp(b))
	return b, nil
}

// Bot type.
type Bot struct {
	api  *tgbotapi.BotAPI
	cmds []*Command
}

func (b *Bot) close() {
	b.api.StopReceivingUpdates()
	b.cmds = nil
}

func (b *Bot) setupCommands() error {
	cfg := tgbotapi.SetMyCommandsConfig{
		Commands: make([]tgbotapi.BotCommand, len(b.cmds)),
	}
	for i := range b.cmds {
		cfg.Commands[i] = tgbotapi.BotCommand{
			Command:     b.cmds[i].Cmd,
			Description: b.cmds[i].Description,
		}
	}
	_, err := b.api.Request(cfg)
	return err
}

// Run starts bot.
func (b *Bot) Run(ctx context.Context) error {
	err := b.setupCommands()
	if err != nil {
		return err
	}
	ch := b.api.GetUpdatesChan(tgbotapi.UpdateConfig{
		Timeout: 60,
	})

	for {
		var upd tgbotapi.Update
		select {
		case <-ctx.Done():
			b.close()
			return context.Canceled
		case upd = <-ch:
		}
		switch {
		case upd.Message != nil:
			b.onMessage(upd.Message)
		case upd.CallbackQuery != nil:
			b.onCallbackQuery(upd.CallbackQuery)
		}
	}
}

func (b *Bot) onCallbackQuery(upd *tgbotapi.CallbackQuery) {}

func (b *Bot) onMessage(upd *tgbotapi.Message) {
	if upd.Chat.IsChannel() {
		return
	}

	text := upd.Text
	if text == "" {
		text = upd.Caption
	}
	cmd, args := b.parseCommand(text)
	if cmd == nil {
		return
	}

	ctx := b.makeContext(cmd, upd.Chat)
	go cmd.Run(ctx, upd, args)
}

func (b *Bot) parseCommand(c string) (cmd *Command, args []string) {
	if len(c) == 0 || c[0] != Prefix {
		return nil, nil
	}
	split := strings.Split(c[1:], " ")
	c, args = split[0], split[1:]
	if i := strings.Index(c, "@"); i != -1 && len(c) > i+1 {
		if b.api.Self.UserName != c[i+1:] {
			return nil, nil
		}
		c = c[:i]
	}
	args = split[1:]
	for _, cmd := range b.cmds {
		if c == cmd.Cmd {
			return cmd, args
		}
	}
	return nil, nil
}
