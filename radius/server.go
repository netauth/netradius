package radius

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2868"
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

	vals, err := s.n.EntityKVGet(r.Context(), username, "radius::vlan")
	if err != nil {
		s.log.Error("Error retrieving vlan information for entity", "entity", username, "error", err)
	}

	resp := r.Response(radius.CodeAccessAccept)

	if vals != nil {
		vlan := vals["radius::vlan"][0]
		rfc2868.TunnelType_Add(resp, 0, 13) // Unifi uses an unspecified value :/
		rfc2868.TunnelMediumType_Add(resp, 0, rfc2868.TunnelMediumType_Value_IEEE802)
		rfc2868.TunnelPrivateGroupID_AddString(resp, 0, vlan)
	}
	s.log.Debug("Writing response", "code", radius.CodeAccessAccept, "client", r.RemoteAddr)
	w.Write(resp)
}

// Serve serves the RADIUS services.
func (s *Server) Serve() error {
	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(s.handler),
		SecretSource: radius.StaticSecretSource([]byte(s.secret)),
	}
	s.radsrv = server

	s.log.Info("Serving Radius on :1812")
	return server.ListenAndServe()
}

// Shutdown requests the underlying packetserver to shutdown smoothly.
func (s *Server) Shutdown() error {
	return s.radsrv.Shutdown(context.Background())
}
