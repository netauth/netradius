package radius

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"layeh.com/radius"
)

// Server owns the radius server itself and its handlers.
type Server struct {
	log hclog.Logger
	n   netauth

	radsrv radius.PacketServer

	secret string
}

// Option enables passing of various options to the server on startup.
type Option func(*Server) error

type netauth interface {
	AuthEntity(context.Context, string, string) error
	EntityKVGet(context.Context, string, string) (map[string][]string, error)
}
