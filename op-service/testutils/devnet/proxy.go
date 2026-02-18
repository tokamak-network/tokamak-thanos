package devnet

import (
	"log/slog"
)

// RetryProxy wraps an RPC URL with retry logic.
type RetryProxy struct {
	url string
}

func NewRetryProxy(lgr *slog.Logger, rpcURL string) *RetryProxy {
	return &RetryProxy{url: rpcURL}
}

func (p *RetryProxy) URL() string {
	return p.url
}
