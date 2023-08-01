package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

var DataProductName string = "ObjDataProductExample"
var DataProductDescription string = "This is a demo Obj-based Data Product"

func main() {
	nc, _ := nats.Connect("rainforest_user:password@localhost:4111")

	js, err := nc.JetStream()
	if err != nil {
		log.Println(err)
		return
	}

	obj, err := js.CreateObjectStore(&nats.ObjectStoreConfig{
		Bucket:      DataProductName,
		Description: DataProductDescription,
		// TTL         time.Duration
		// MaxBytes    int64
		// Storage     StorageType
		// Replicas    int
		// Placement   *Placement
	})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(obj)
}
