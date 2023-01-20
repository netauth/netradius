package radius

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

// New initializes the server and prepares it for serving
func New(opts ...Option) (*Server, error) {
	s := &Server{
		log: hclog.NewNullLogger(),
	}

	for _, o := range opts {
		if err := o(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Server) handler(w radius.ResponseWriter, r *radius.Request) {
	s.log.Debug("handler inbound", "request", fmt.Sprintf("%#v", r.Packet))

	username := rfc2865.UserName_GetString(r.Packet)
	password := rfc2865.UserPassword_GetString(r.Packet)

	if err := s.n.AuthEntity(r.Context(), username, password); err != nil {
		s.log.Warn("Error authenticating user", "error", err)
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}
	s.log.Debug("Writing response", "code", radius.CodeAccessAccept, "client", r.RemoteAddr)
	w.Write(r.Response(radius.CodeAccessAccept))
}

// Serve serves the RADIUS services.
func (s *Server) Serve() error {
	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(s.handler),
		SecretSource: radius.StaticSecretSource([]byte("secret")),
	}

	s.log.Info("Serving Radius on :1812")
	return server.ListenAndServe()
}
