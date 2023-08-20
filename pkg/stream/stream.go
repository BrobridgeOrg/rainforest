package stream

import (
	"log"
	"net/url"

	"github.com/nats-io/nats-server/v2/server"
)

type StreamServer struct {
	natsServer *server.Server
}

type StreamServerConfig struct {
	Port     int      `mapstructure:"port"`
	Domain   string   `mapstructure:"domain"`
	HubURLs  []string `mapstructure:"hub-urls"`
	LeafPort int      `mapstructure:"leaf-port"`
	StoreDir string   `mapstructure:"store-dir"`
}

func NewStreamServer(cfg StreamServerConfig) *StreamServer {
	routes := []*url.URL{}
	for _, hub := range cfg.HubURLs {
		routes = append(routes, &url.URL{
			Host: hub,
		})
	}

	s, err := server.NewServer(&server.Options{
		Port: cfg.Port,
		JetStream:       true,
		JetStreamDomain: cfg.Domain,
		StoreDir:        cfg.StoreDir,
		LeafNode: server.LeafNodeOpts{
			Host: "0.0.0.0",
			// Port: cfg.LeafPort, not use in lead server
		},
		Routes: routes,
	})

	if err != nil {
		log.Fatal(err)
	}

	s.ConfigureLogger()

	return &StreamServer{
		natsServer: s,
	}
}

func (s *StreamServer) Start() {
	s.natsServer.Start()
}
