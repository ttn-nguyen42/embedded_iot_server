package events

import (
	"context"
	"crypto/tls"
	"labs/htmx-blog/internal/configs"
	custerror "labs/htmx-blog/internal/error"
	"log"
	"time"

	"github.com/nats-io/nats-server/v2/server"
)

type EmbeddedNats struct {
	configs *Options

	server *server.Server
	name   string
}

func New(options ...Optioner) *EmbeddedNats {
	opts := Options{}
	for _, opt := range options {
		opt(&opts)
	}
	serverConfigs := opts.configs

	server := buildServer(serverConfigs)

	if server == nil {
		return nil
	}

	return &EmbeddedNats{
		configs: &opts,
		server:  server,
		name:    serverConfigs.Name,
	}
}

func buildServer(configs *configs.EventStoreConfigs) *server.Server {
	if configs == nil {
		return nil
	}
	if !configs.Enabled {
		return nil
	}

	log.Printf("buildServer: build config for NATS server")

	serverOptions := server.Options{
		Host:                   configs.Host,
		Port:                   configs.Port,
		ServerName:             configs.Name,
		Username:               configs.Username,
		Password:               configs.Password,
		JetStream:              true,
		DisableJetStreamBanner: true,
	}

	if configs.Tls.Enabled() {
		serverTls, err := buildTlsConfigs(&configs.Tls)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		serverOptions.TLSConfig = serverTls
		serverOptions.TLS = true
	}

	server, err := server.NewServer(&serverOptions)
	if err != nil {
		log.Fatalf("buildServer: build server err = %s", err)
		return nil
	}

	server.ConfigureLogger()
	return server
}

func buildTlsConfigs(tlsConfigs *configs.TlsConfig) (*tls.Config, error) {
	configs, err := server.GenTLSConfig(&server.TLSConfigOpts{
		CertFile: tlsConfigs.Cert,
		KeyFile:  tlsConfigs.Key,
		CaFile:   tlsConfigs.Authority,
		Verify:   false,
		Insecure: true,
	})
	if err != nil {
		return nil, custerror.FormatInternalError("buildTlsConfigs: err = %s", err)
	}
	return configs, nil

}

type Options struct {
	configs *configs.EventStoreConfigs
}

type Optioner func(opts *Options)

func WithNatsConfigs(configs *configs.EventStoreConfigs) Optioner {
	return func(opts *Options) {
		opts.configs = configs
	}
}

func (n *EmbeddedNats) Start() error {
	n.server.Start()
	if !n.server.ReadyForConnections(time.Second * 3) {
		return custerror.FormatInternalError("EmbeddedNats.Start: connection not ready")
	}
	return nil
}

func (n *EmbeddedNats) Stop(ctx context.Context) error {
	n.server.Shutdown()
	n.server.WaitForShutdown()
	return nil
}

func (n *EmbeddedNats) Name() string {
	return n.name
}
