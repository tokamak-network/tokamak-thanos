package endpoint

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
)

type MustRPC struct {
	Value RPC
}

func (u *MustRPC) UnmarshalText(data []byte) error {
	if u == nil {
		return fmt.Errorf("cannot unmarshal %q into nil MustRPC", string(data))
	}
	v := URL(data)
	if v == "" {
		return errors.New("empty RPC URL")
	}
	u.Value = v
	return nil
}

func (u MustRPC) MarshalText() ([]byte, error) {
	if u.Value == nil {
		return nil, errors.New("missing RPC")
	}
	out := u.Value.RPC()
	if out == "" {
		return nil, errors.New("missing RPC")
	}
	return []byte(out), nil
}

type OptionalRPC struct {
	Value RPC
}

func (u *OptionalRPC) UnmarshalText(data []byte) error {
	if u == nil {
		return fmt.Errorf("cannot unmarshal %q into nil OptionalRPC", string(data))
	}
	u.Value = URL(data)
	return nil
}

func (u OptionalRPC) MarshalText() ([]byte, error) {
	if u.Value == nil {
		return []byte{}, nil
	}
	out := u.Value.RPC()
	return []byte(out), nil
}

type RPC interface {
	RPC() string
}

type URL string

var _ RPC = URL("")

func (u URL) RPC() string {
	return string(u)
}

type WsRPC interface {
	RPC
	WsRPC() string
}

type HttpRPC interface {
	RPC
	HttpRPC() string
}

type ClientRPC interface {
	RPC
	ClientRPC() *rpc.Client
}

type HttpURL string

func (url HttpURL) RPC() string    { return string(url) }
func (url HttpURL) HttpRPC() string { return string(url) }

type WsURL string

func (url WsURL) RPC() string   { return string(url) }
func (url WsURL) WsRPC() string { return string(url) }

type WsOrHttpRPC struct {
	WsURL   string
	HttpURL string
}

func (r WsOrHttpRPC) RPC() string     { return r.WsURL }
func (r WsOrHttpRPC) WsRPC() string   { return r.WsURL }
func (r WsOrHttpRPC) HttpRPC() string { return r.HttpURL }

type ServerRPC struct {
	Fallback WsOrHttpRPC
	Server   *rpc.Server
}

func (e *ServerRPC) RPC() string     { return e.Fallback.RPC() }
func (e *ServerRPC) WsRPC() string   { return e.Fallback.WsRPC() }
func (e *ServerRPC) HttpRPC() string { return e.Fallback.HttpRPC() }
func (e *ServerRPC) ClientRPC() *rpc.Client { return rpc.DialInProc(e.Server) }

type Dialer func(v string) *rpc.Client

type RPCPreference int

const (
	PreferAnyRPC RPCPreference = iota
	PreferHttpRPC
	PreferWSRPC
)

func DialRPC(preference RPCPreference, r RPC, dialer Dialer) *rpc.Client {
	if v, ok := r.(HttpRPC); preference == PreferHttpRPC && ok {
		return dialer(v.HttpRPC())
	}
	if v, ok := r.(WsRPC); preference == PreferWSRPC && ok {
		return dialer(v.WsRPC())
	}
	if v, ok := r.(ClientRPC); ok {
		return v.ClientRPC()
	}
	return dialer(r.RPC())
}

func SelectRPC(preference RPCPreference, r RPC) string {
	if v, ok := r.(HttpRPC); preference == PreferHttpRPC && ok {
		return v.HttpRPC()
	}
	if v, ok := r.(WsRPC); preference == PreferWSRPC && ok {
		return v.WsRPC()
	}
	return r.RPC()
}
