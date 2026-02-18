package logmods

import (
	"log/slog"
)

type HandlerMod func(slog.Handler) slog.Handler

type Handler interface {
	slog.Handler
	Unwrap() slog.Handler
}

func FindHandler[H slog.Handler](h slog.Handler) (out H, ok bool) {
	for {
		if h == nil {
			ok = false
			return
		}
		if found, tempOk := h.(H); tempOk {
			return found, true
		}
		unwrappable, tempOk := h.(Handler)
		if !tempOk {
			ok = false
			return
		}
		h = unwrappable.Unwrap()
	}
}
