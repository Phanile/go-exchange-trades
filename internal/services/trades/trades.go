package trades

import "log/slog"

type Trades struct {
	log *slog.Logger
}

func NewTradesService(log *slog.Logger) *Trades {
	return &Trades{
		log: log,
	}
}
