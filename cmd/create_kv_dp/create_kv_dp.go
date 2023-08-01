package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

var DataProductName string = "KVDataProductExample"
var DataProductDescription string = "This is a demo KV-based Data Product"

func main() {
	nc, _ := nats.Connect("rainforest_user:password@localhost:4111")

	js, err := nc.JetStream()
	if err != nil {
		log.Println(err)
		return
	}

	kv, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:      DataProductName,
		Description: DataProductDescription,
		// MaxValueSize int32
		// History      uint8
		// TTL          time.Duration
		// MaxBytes     int64
		// Storage      StorageType
		// Replicas     int
		// Placement    *Placement
		// RePublish    *RePublish
		// Mirror       *StreamSource
		// Sources      []*StreamSource
	})
	if err != nil {
		log.Println(err)
		return
	}

	kv.Put("TestRecord", []byte("Value"))

}
