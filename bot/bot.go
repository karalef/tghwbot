package bot

import (
	"context"
	"net/http"
	"strings"
	"tghwbot/logger"
	"time"

	"gopkg.in/telebot.v3"
)

// New creates new bot.
func New(token string, log *logger.Logger, cmds ...*Command) (*Bot, error) {
	api, err := telebot.NewBot(telebot.Settings{
		Token:   token,
		Updates: 5,
		Poller: &telebot.LongPoller{
			Timeout: time.Second * 20,
		},
		Client: &http.Client{},
	})
	if err != nil {
		return nil, err
	}
	b := &Bot{
		api:  api,
		log:  log,
		cmds: cmds,
	}
	b.cmds = append(b.cmds, &ping, makeHelp(b))
	return b, nil
}

// Bot type.
type Bot struct {
	api  *telebot.Bot
	log  *logger.Logger
	cmds []*Command
}

func (b *Bot) setupCommands() error {
	commands := make([]telebot.Command, len(b.cmds))
	for i := range b.cmds {
		commands[i] = telebot.Command{
			Text:        b.cmds[i].Cmd,
			Description: b.cmds[i].Description,
		}
	}
	return b.api.SetCommands(commands)
}

// Run starts bot.
func (b *Bot) Run(ctx context.Context) error {
	err := b.setupCommands()
	if err != nil {
		return err
	}

	stop := make(chan struct{})
	b.api.Poller.Poll(b.api, b.api.Updates, stop)

	for {
		var upd telebot.Update
		select {
		case <-ctx.Done():
			close(stop)
			return context.Canceled
		case upd = <-b.api.Updates:
		}
		switch {
		case upd.Message != nil:
			b.onMessage(upd.Message)
		}
	}
}

func (b *Bot) onMessage(upd *telebot.Message) {
	if upd.FromChannel() {
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

	ctx := b.makeContext(cmd, upd)
	go cmd.Run(ctx, upd, args)
}

func (b *Bot) parseCommand(c string) (cmd *Command, args []string) {
	if len(c) == 0 || c[0] != Prefix {
		return nil, nil
	}
	split := strings.Split(c[1:], " ")
	c, args = split[0], split[1:]
	if i := strings.Index(c, "@"); i != -1 && len(c) > i+1 {
		if b.api.Me.Username != c[i+1:] {
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
