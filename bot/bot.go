package bot

import (
	"context"
	"net/http"
	"strings"
	"tghwbot/bot/logger"
	"tghwbot/bot/tg"
)

// New creates new bot.
func New(token string, log *logger.Logger, cmds ...*Command) (*Bot, error) {
	b := Bot{
		token:  token,
		apiURL: tg.DefaultAPIURL,
		client: &http.Client{},
		log:    log,
	}

	me, err := b.getMe()
	if err != nil {
		return nil, err
	}
	b.Me = me

	b.cmds = append(cmds, &ping, makeHelp(&b))
	return &b, nil
}

// Bot type.
type Bot struct {
	token  string
	apiURL string
	client *http.Client
	log    *logger.Logger

	sync bool
	cmds []*Command

	Me *tg.User
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

// Run starts bot.
func (b *Bot) Run(ctx context.Context, lastUpdate int) error {
	err := b.setupCommands()
	if err != nil {
		return err
	}

	for {
		upds, err := b.getUpdates(ctx, lastUpdate+1, 30, nil)
		if err != nil {
			//TODO
			return err
		}
		for _, upd := range upds {
			switch {
			case upd.Message != nil:
				b.onMessage(upd.Message)
			}
			lastUpdate = upd.ID
		}
	}
}

func (b *Bot) onMessage(msg *tg.Message) {
	if msg.Chat.IsChannel() {
		return
	}

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
