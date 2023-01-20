package radius

import (
	"context"

	"github.com/hashicorp/go-hclog"
)

// Server owns the radius server itself and its handlers.
type Server struct {
	log hclog.Logger
	n   netauth
}

// Option enables passing of various options to the server on startup.
type Option func(*Server) error

type netauth interface {
	AuthEntity(context.Context, string, string) error
}
