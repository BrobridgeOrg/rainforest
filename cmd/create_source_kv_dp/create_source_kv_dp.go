package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

var DataProductName string = "SourceKVDataProductExample"
var DataProductDescription string = "This is a demo Source-KV-based Data Product"
var SourceDomain string = "rainforest_leaf_0"
var SourceDataProductName = "KVDataProductExample"

func main() {
	nc, _ := nats.Connect("rainforest_user:password@localhost:4112")

	js, err := nc.JetStream()
	if err != nil {
		log.Println(err)
		return
	}

	_, err = js.CreateKeyValue(&nats.KeyValueConfig{
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
		Sources: []*nats.StreamSource{
			{
				Domain: SourceDomain,
				Name:   SourceDataProductName,
				// OptStartSeq   uint64          `json:"opt_start_seq,omitempty"`
				// OptStartTime  *time.Time      `json:"opt_start_time,omitempty"`
				// FilterSubject string          `json:"filter_subject,omitempty"`
				// External      *ExternalStream `json:"external,omitempty"`
			},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

}
