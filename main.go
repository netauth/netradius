package main

import (
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/netauth/netauth/pkg/netauth"
	"github.com/netauth/netradius/radius"
	"github.com/spf13/viper"
)

func main() {
	var appLogger hclog.Logger

	llevel := os.Getenv("NETAUTH_LOGLEVEL")
	if llevel != "" {
		appLogger = hclog.New(&hclog.LoggerOptions{
			Name:  "netradius",
			Level: hclog.LevelFromString(llevel),
		})
	} else {
		appLogger = hclog.NewNullLogger()
	}

	// Take over the built in logger and set it up for Trace level
	// priority.  The only thing that logs at this priority are
	// protocol messages from the underlying ldap server mux.
	log.SetOutput(appLogger.Named("radius.protocol").
		StandardWriter(
			&hclog.StandardLoggerOptions{
				ForceLevel: hclog.Trace,
			},
		),
	)

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/netauth/")
	viper.AddConfigPath("$HOME/.netauth/")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("NETAUTH")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		appLogger.Error("Error loading config", "error", err)
		os.Exit(5)
	}

	nacl, err := netauth.NewWithLog(appLogger.Named("netauth"))
	if err != nil {
		appLogger.Error("Error initializing client", "error", err)
		os.Exit(2)
	}
	nacl.SetServiceName("netradius")

	srvr, err := radius.New(radius.WithLogger(appLogger), radius.WithNetAuth(nacl))
	if err != nil {
		appLogger.Error("Error initializing", "error", err)
		os.Exit(1)
	}

	if err := srvr.Serve(); err != nil {
		appLogger.Error("Error serving", "error", err)
		os.Exit(1)
	}
}
