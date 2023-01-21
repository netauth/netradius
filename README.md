# NetAuth RADIUS Server

The NetAuth RADIUS server acts as a RADIUS protocol bridge to allow
basic authentication against the NetAuth entity authentication flow
for services that understand the RADIUS protocol, sometimes referred
to as 802.1X.

The RADIUS server is extremely simplistic and only supports the PAP
profile, so use with care.  For this reason you should try to install
the server as close to your RADIUS clients as possible to avoid
information being transmitted insecurely.  There is no need to run a
single netradius server for an entire site, an arbitrarily many
servers may be deployed.

The NAS Secret implementation is very simplistic, expecting a single
shared key for all clients.  You may set this key via the
`NETAUTH_RADIUS_SECRET` environment variable.

