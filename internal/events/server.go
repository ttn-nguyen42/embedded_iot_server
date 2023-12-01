package events

import (
	"context"
	"crypto/tls"
	"labs/htmx-blog/internal/configs"
	custerror "labs/htmx-blog/internal/error"
	"log"

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
	mqttConfigs := opts.mqttConfigs

	server := buildServer(serverConfigs, mqttConfigs)

	if server == nil {
		return nil
	}

	return &EmbeddedNats{
		configs: &opts,
		server:  server,
		name:    serverConfigs.Name,
	}
}

func buildServer(configs *configs.EventStoreConfigs, mqttConfigs *configs.EventStoreConfigs) *server.Server {
	if configs == nil || !configs.Enabled {
		return nil
	}
	var mqttOpts server.MQTTOpts

	if mqttConfigs != nil && mqttConfigs.Enabled {
		log.Printf("buildServer: build config for MQTT server")
		mqttTls, err := buildTlsConfigs(&mqttConfigs.Tls)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		mqttOpts = server.MQTTOpts{
			Host:      mqttConfigs.Host,
			Port:      mqttConfigs.Port,
			TLSConfig: mqttTls,
			Username:  mqttConfigs.Username,
			Password:  mqttConfigs.Password,
		}
	}

	log.Printf("buildServer: build config for NATS server")
	serverTls, err := buildTlsConfigs(&configs.Tls)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	serverOptions := server.Options{
		Host:       configs.Host,
		Port:       configs.Port,
		ServerName: configs.Name,
		MQTT:       mqttOpts,
		TLSConfig:  serverTls,
		Username:   configs.Username,
		Password:   configs.Password,
		JetStream:  true,
		// 1GB
		JetStreamMaxMemory: 1073741824,
		// 5GB
		JetStreamMaxStore: 5737418240,
	}

	server, err := server.NewServer(&serverOptions)
	if err != nil {
		log.Fatalf("buildServer: build server err = %s", err)
		return nil
	}

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
	configs     *configs.EventStoreConfigs
	mqttConfigs *configs.EventStoreConfigs
}

type Optioner func(opts *Options)

func WithNatsConfigs(configs *configs.EventStoreConfigs) Optioner {
	return func(opts *Options) {
		opts.configs = configs
	}
}

func WithMqttConfigs(configs *configs.EventStoreConfigs) Optioner {
	return func(opts *Options) {
		opts.mqttConfigs = configs
	}
}

func (n *EmbeddedNats) Start() error {
	n.server.Start()
	return nil
}

func (n *EmbeddedNats) Stop(ctx context.Context) error {
	n.server.Shutdown()
	return nil
}

func (n *EmbeddedNats) Name() string {
	return n.name
}
