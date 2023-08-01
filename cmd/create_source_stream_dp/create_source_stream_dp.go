package main

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var DataProductName string = "SourceStreamDataProductExample"
var DataProductDescription string = "This is a demo Source-Stream-based Data Product"
var DataProductSubject string = "$RAINFOREST.DATAPRODUCT." + DataProductName

var SourceDomain string = "rainforest_leaf_0"
var SourceDataProductName = "StreamDataProductExample"

func main() {
	nc, _ := nats.Connect("rainforest_user:password@localhost:4112")

	js, _ := jetstream.New(nc)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	str, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:        DataProductName,
		Description: DataProductDescription,
		Sources: []*jetstream.StreamSource{
			{
				Domain: SourceDomain,
				Name:   SourceDataProductName,
			},
		},
	})

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(str.Info(ctx))
}
