package main

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var DataProductName string = "StreamDataProductExample"
var DataProductDescription string = "This is a demo Stream-based Data Product"
var DataProductSubject string = "$RAINFOREST.DATAPRODUCT." + DataProductName

func main() {
	nc, _ := nats.Connect("rainforest_user:password@localhost:4111")

	js, _ := jetstream.New(nc)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	str, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:        DataProductName,
		Description: DataProductDescription,
		Subjects:    []string{DataProductSubject},
	})

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(str.Info(ctx))
}
