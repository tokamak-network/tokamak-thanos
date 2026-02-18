package rpc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"

	gethrpc "github.com/ethereum/go-ethereum/rpc"
	opmetrics "github.com/tokamak-network/tokamak-thanos/op-service/metrics"
)

// WithWebsocketEnabled enables WebSocket support on the server.
func WithWebsocketEnabled() ServerOption {
	return func(b *Server) {
		// WebSocket support is always enabled in this build
	}
}

// WithRPCRecorder sets the RPC recorder for metrics.
func WithRPCRecorder(_ opmetrics.HTTPRecorder) ServerOption {
	return func(b *Server) {
		// RPC recording stub
	}
}

// Port returns the listening port of the server.
func (s *Server) Port() (int, error) {
	if s.listener == nil {
		return 0, fmt.Errorf("server not started")
	}
	addr := s.listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// NewStream creates a typed event stream. This is a stub for interop compatibility.
func NewStream[T any](log interface{}, bufSize int) *Stream[T] {
	return &Stream[T]{bufSize: bufSize, ch: make(chan *T, bufSize)}
}

// Stream is a typed event stream stub.
type Stream[T any] struct {
	bufSize int
	ch      chan *T
}

// Send sends an event to the stream.
func (s *Stream[T]) Send(event *T) {
	// best-effort send
	select {
	case s.ch <- event:
	default:
	}
}

// Serve blocks until an event is available and returns it.
func (s *Stream[T]) Serve() (*T, error) {
	ev, ok := <-s.ch
	if !ok {
		return nil, fmt.Errorf("stream closed")
	}
	return ev, nil
}

// Subscribe creates an RPC subscription for the stream.
func (s *Stream[T]) Subscribe(ctx context.Context) (*gethrpc.Subscription, error) {
	return nil, fmt.Errorf("stream subscriptions not supported in this build")
}

// ObtainJWTSecret reads or generates a JWT secret for RPC authentication.
func ObtainJWTSecret(log interface{}, path string, generate bool) ([32]byte, error) {
	if path == "" {
		return [32]byte{}, fmt.Errorf("JWT secret path is empty")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if generate {
			// Generate a random secret
			var secret [32]byte
			if _, err := rand.Read(secret[:]); err != nil {
				return [32]byte{}, fmt.Errorf("generating JWT secret: %w", err)
			}
			if err := os.WriteFile(path, []byte(fmt.Sprintf("%x", secret)), 0600); err != nil {
				return [32]byte{}, fmt.Errorf("writing JWT secret: %w", err)
			}
			return secret, nil
		}
		return [32]byte{}, fmt.Errorf("reading JWT secret: %w", err)
	}
	// Parse hex-encoded secret
	var secret [32]byte
	cleaned := strings.TrimSpace(strings.TrimPrefix(string(data), "0x"))
	b, err := hex.DecodeString(cleaned)
	if err != nil {
		return [32]byte{}, fmt.Errorf("decoding JWT secret: %w", err)
	}
	copy(secret[:], b)
	return secret, nil
}
