package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	apiv1 "github.com/Awareness-Labs/rainforest/pkg/proto/api/v1"
	"github.com/Awareness-Labs/rainforest/pkg/server"
	"github.com/dgraph-io/badger/v4"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/encoding/protojson"
)

// Ref: https://dgraph.io/docs/badger/get-started/
// Reference subject: $RAINFOREST.API.DATAPRODUCT.<data_product>.*
// Consume subject:   $RAINFOREST.API.DATAPRODUCT.>
// Service subject:   $RAINFOREST.API.KV.<data_product>

type KeyValueConsumer struct {
	db            *badger.DB
	conn          *nats.Conn
	streamManager jetstream.JetStream
}

func NewKeyValueConsumer(nc *nats.Conn) *KeyValueConsumer {

	db, err := badger.Open(badger.DefaultOptions("./data/badger"))
	if err != nil {
		log.Fatal(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	return &KeyValueConsumer{
		db:            db,
		conn:          nc,
		streamManager: js,
	}
}

// Consume from StateStream, materialize state into key-value database
func (c *KeyValueConsumer) Start() {

	js, err := jetstream.New(c.conn)
	if err != nil {
		log.Printf("Failed to initialize JetStream: %v", err)
		return
	}
	// list stream names
	go PeriodSched(js, c)

	c.StartKeyValueService()

	// Enable Badger GC
	go badgerGC(c.db)

	if err != nil {
		log.Printf("Failed to initialize Consumer: %v", err)
		return
	}

}

func PeriodSched(js jetstream.JetStream, c *KeyValueConsumer) {
	deployedConsumer := map[string]bool{}
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		// log.Println("Sched is working bitch!")
		ctx := context.TODO() // Use the appropriate context

		names := c.streamManager.StreamNames(ctx)
		for name := range names.Name() {

			_, ok := deployedConsumer[name]
			if ok {
				continue
			}

			if strings.HasPrefix(name, server.StateDataProductPrefix) {

				cons, err := js.OrderedConsumer(ctx, name, jetstream.OrderedConsumerConfig{})

				if err != nil {
					log.Printf("Failed to initialize Consumer: %v", err)
					return
				}
				go ConsumerFunc(cons, c.db)
				deployedConsumer[name] = true
			}
		}
	}
}

func ConsumerFunc(cons jetstream.Consumer, db *badger.DB) {
	for {

		msgs, err := cons.Fetch(1)
		if msgs.Error() != nil {
			// handle error
			log.Println(err)
			continue
		}
		for msg := range msgs.Messages() {
			fmt.Printf("Received a JetStream message: %s\n", string(msg.Data()))

			msg.Headers().Get("")

			k := strings.TrimPrefix(msg.Subject(), "$RAINFOREST.DP.STATE.")
			k = strings.Replace(k, ".", "/", -1)
			fmt.Println("k:", []byte(k), "v:", msg.Data())

			// Materilaize
			txn := db.NewTransaction(true)
			txn.Set([]byte(k), msg.Data())

			err := txn.Commit()
			if err != nil {
				log.Println(err)
			}
			msg.Ack()
		}

	}
}

func (c *KeyValueConsumer) StartKeyValueService() {
	c.conn.Subscribe("$RAINFOREST.API.KV.*", func(msg *nats.Msg) {

		req := &apiv1.KeyValueRequest{}
		protojson.Unmarshal(msg.Data, req)
		log.Println(req)
		switch req.GetOperation().(type) {
		case *apiv1.KeyValueRequest_Scan:
			err := c.db.View(func(txn *badger.Txn) error {
				log.Println("scan")
				kvs := []*apiv1.KeyValue{}

				opts := badger.DefaultIteratorOptions
				opts.Reverse = req.GetScan().GetReverse()

				it := txn.NewIterator(opts)
				defer it.Close()
				log.Println("Scan Op", []byte(req.GetScan().GetStartKey()), []byte(req.GetScan().GetEndKey()))
				for it.Seek([]byte(req.GetScan().GetStartKey())); it.ValidForPrefix([]byte(req.GetScan().GetEndKey())); it.Next() {

					item := it.Item()
					item.Value(func(v []byte) error {
						kvs = append(kvs, &apiv1.KeyValue{
							Key:   string(item.Key()),
							Value: string(v),
						})
						return nil
					})
				}
				log.Println(kvs)
				res := &apiv1.KeyValueDataResponse{
					Kvs: kvs,
				}
				resData, err := protojson.Marshal(res)
				if err != nil {
					sendErrorResponse(msg, err)
					return err
				}

				msg.Respond(resData)

				return nil
			})

			if err != nil {
				sendErrorResponse(msg, err)
				return
			}
		case *apiv1.KeyValueRequest_Get:
			err := c.db.View(func(txn *badger.Txn) error {
				kvs := []*apiv1.KeyValue{}
				item, err := txn.Get([]byte(req.GetScan().StartKey))
				if err != nil {
					sendErrorResponse(msg, err)
					return err
				}
				item.Value(func(v []byte) error {
					kvs = append(kvs, &apiv1.KeyValue{
						Key:   string(item.Key()),
						Value: string(v),
					})
					return nil
				})
				res := &apiv1.KeyValueDataResponse{
					Kvs: kvs,
				}
				resData, err := protojson.Marshal(res)
				if err != nil {
					sendErrorResponse(msg, err)
					return err
				}

				msg.Respond(resData)
				return nil
			})

			if err != nil {
				sendErrorResponse(msg, err)
				return
			}
		}

	})
}

func sendErrorResponse(m *nats.Msg, err error) {
	res := &apiv1.KeyValueResponse{
		Response: &apiv1.KeyValueResponse_ErrorResponse{
			ErrorResponse: &apiv1.ErrorResponse{
				ErrorCode:    "400",
				ErrorMessage: err.Error(),
			},
		},
	}
	resData, _ := json.Marshal(res)
	m.Respond(resData)
}

func badgerGC(db *badger.DB) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
	again:
		err := db.RunValueLogGC(0.7)
		if err == nil {
			goto again
		}
	}
}
