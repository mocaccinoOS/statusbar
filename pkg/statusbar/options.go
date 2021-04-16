package statusbar

import "github.com/MocaccinoOS/statusbar/pkg/block"

type Options struct {
	Blocks                                              []block.Block
	Title, Tooltip, NotificatorAppName, NotificatorIcon string
}

type Option func(cfg *Options) error

func (cfg *Options) Apply(opts ...Option) error {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(cfg); err != nil {
			return err
		}
	}
	return nil
}

func WithBlocks(b ...block.Block) func(cfg *Options) error {
	return func(cfg *Options) error {
		cfg.Blocks = b
		return nil
	}
}

func WithTitle(s string) func(cfg *Options) error {
	return func(cfg *Options) error {
		cfg.Title = s
		return nil
	}
}

func WithTooltip(s string) func(cfg *Options) error {
	return func(cfg *Options) error {
		cfg.Tooltip = s
		return nil
	}
}

func WithAppName(s string) func(cfg *Options) error {
	return func(cfg *Options) error {
		cfg.NotificatorAppName = s
		return nil
	}
}

func WithNotificationIcon(s string) func(cfg *Options) error {
	return func(cfg *Options) error {
		cfg.NotificatorIcon = s
		return nil
	}
}
