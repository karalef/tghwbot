package bot

import (
	"context"
	"errors"
	"net/http"
	"runtime"
	"strings"
	"tghwbot/bot/logger"
	"tghwbot/bot/tg"
)

// Config contains bot configuration.
type Config struct {
	Token    string
	Logger   *logger.Logger // logger.DefaultWriter if empty
	Commands []*Command
	MakeHelp bool
}

// New creates new bot.
func New(config Config) (*Bot, error) {
	if config.Token == "" {
		return nil, errors.New("no token provided")
	}
	if config.Logger == nil {
		config.Logger = logger.New(logger.DefaultWriter, "")
	}
	b := Bot{
		token:   config.Token,
		apiURL:  tg.DefaultAPIURL,
		fileURL: tg.DefaultFileURL,
		client:  &http.Client{},
		log:     config.Logger,
		cmds:    config.Commands,
	}
	if config.MakeHelp {
		b.cmds = append(b.cmds, makeHelp(&b))
	}

	me, err := b.getMe()
	if err != nil {
		return nil, err
	}
	b.Me = me

	return &b, nil
}

// Bot type.
type Bot struct {
	token   string
	apiURL  string
	fileURL string
	client  *http.Client
	log     *logger.Logger

	stop context.CancelFunc
	cmds []*Command

	Me *tg.User
}

func (b *Bot) closeExecution() {
	runtime.Goexit()
}

func (b *Bot) setupCommands() error {
	commands := make([]tg.Command, len(b.cmds))
	for i := range b.cmds {
		commands[i] = tg.Command{
			Command:     b.cmds[i].Cmd,
			Description: b.cmds[i].Description,
		}
	}
	return b.setCommands(&commandParams{Commands: commands})
}

// Stop stops polling for updates.
func (b *Bot) Stop() {
	if b.stop != nil {
		b.stop()
	}
}

// Run starts bot.
// Returns nil if context is closed.
func (b *Bot) Run(ctx context.Context, lastUpdate int) error {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, b.stop = context.WithCancel(ctx)

	err := b.setupCommands()
	if err != nil {
		return err
	}

	for {
		upds, err := b.getUpdates(ctx, lastUpdate+1, 30, "message")
		switch err {
		case nil:
		case context.Canceled, context.DeadlineExceeded:
			return nil
		default:
			return err
		}
		for i := range upds {
			b.handle(&upds[i])
			lastUpdate = upds[i].ID
		}
	}
}

func (b *Bot) handle(upd *tg.Update) {
	switch {
	case upd.Message != nil:
		b.onMessage(upd.Message)
	}
}

func (b *Bot) onMessage(msg *tg.Message) {
	text := msg.Text
	if text == "" {
		text = msg.Caption
	}
	cmd, args := b.parseCommand(text)
	if cmd == nil {
		return
	}

	ctx := b.makeContext(cmd, msg)
	go cmd.Run(ctx, msg, args)
}

func (b *Bot) parseCommand(c string) (cmd *Command, args []string) {
	if len(c) == 0 || c[0] != Prefix {
		return nil, nil
	}
	split := strings.Split(c[1:], " ")
	c, args = split[0], split[1:]
	if i := strings.Index(c, "@"); i != -1 && len(c) > i+1 {
		if b.Me.Username != c[i+1:] {
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
