package radius

import (
	"github.com/hashicorp/go-hclog"
)

// WithLogger sets the logging implementation for the server.
func WithLogger(l hclog.Logger) Option {
	return func(s *Server) error {
		s.log = l.Named("radius")
		return nil
	}
}

// WithNetAuth configures the netauth client used by the server.
func WithNetAuth(n netauth) Option {
	return func(s *Server) error {
		s.n = n
		return nil
	}
}
