package main

import (
	"log"
	"runtime"
	"strings"

	"github.com/Awareness-Labs/rainforest/pkg/proto/consumer"
	"github.com/Awareness-Labs/rainforest/pkg/server"
	"github.com/Awareness-Labs/rainforest/pkg/stream"
	"github.com/nats-io/nats.go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {

	pflag.String("port", "4222", "port to serve on")
	pflag.String("domain", "", "domain of the rainforest server (required)")
	pflag.StringArray("hub-urls", []string{}, "remote connection hub URLs")
	pflag.Int("leaf-port", 7422, "leaf port to start")
	pflag.String("stream-path", "./data/stream", "directory to store stream data")
	pflag.String("kv-path", "", "directory to store key value data")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)
}

func main() {
	// Start embedded Stream Server
	cfg := stream.StreamServerConfig{}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	if strings.TrimSpace(cfg.Domain) == "" {
		log.Fatal("domain is required and cannot be empty")
	}

	// Start NATS, JetStream embedded server
	strServer := stream.NewStreamServer(cfg)
	strServer.Start()

	// Connect to NATS
	nc, err := nats.Connect("localhost:" + cfg.Port)
	if err != nil {
		log.Fatal(err)
	}

	// Start Rainforest server
	rfServer := server.NewServer(nc)
	rfServer.Start()

	// Start KV consumer
	kv := consumer.NewKeyValueConsumer(nc, cfg.KVPath)
	go kv.Start()

	// Wait to stop
	runtime.Goexit()
	defer nc.Close()
}
