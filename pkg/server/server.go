package server

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	api "github.com/Awareness-Labs/rainforest/pkg/proto/api/v1"
	core "github.com/Awareness-Labs/rainforest/pkg/proto/core/v1"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	conn          *nats.Conn
	streamManager jetstream.JetStream
	kvManager     nats.JetStreamContext
	objManager    nats.JetStreamContext
}

func NewServer(nc *nats.Conn) *Server {

	js, err := jetstream.New(nc)
	if err != nil {
		log.Printf("Failed to initialize JetStream: %v", err)
		return nil
	}

	kv, err := nc.JetStream()
	if err != nil {
		log.Printf("Failed to initialize Key-Value Manager: %v", err)
		return nil
	}

	obj, err := nc.JetStream()
	if err != nil {
		log.Printf("Failed to initialize Object Manager: %v", err)
		return nil
	}

	return &Server{
		conn:          nc,
		streamManager: js,
		kvManager:     kv,
		objManager:    obj,
	}
}

func (s *Server) Start() {
	// Subscribe to DataProduct topics
	_, err := s.conn.Subscribe("$RAINFOREST.API.DP.CREATE", s.CreateDataProduct)
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	_, err = s.conn.Subscribe("$RAINFOREST.API.DP.READ", s.ReadDataProduct)
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	_, err = s.conn.Subscribe("$RAINFOREST.API.DP.UPDATE", s.UpdateDataProduct)
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	_, err = s.conn.Subscribe("$RAINFOREST.API.DP.DELETE", s.DeleteDataProduct)
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	_, err = s.conn.Subscribe("$RAINFOREST.API.DP.LIST", s.ListDataProducts)
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}
}

func (s *Server) CreateDataProduct(m *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Serialize
	req := &api.CreateDataProductRequest{}
	if err := proto.Unmarshal(m.Data, req); err != nil {
		sendErrorResponse(m, "Error unmarshaling request")
		return
	}

	// Validate req (Original Data Product or Secondary Data Product)
	if req.GetProduct().Type == core.DataProductType_UNDEFINED {
		sendErrorResponse(m, "Error that not define data product type")
		return
	}

	// Implement logic to create a DataProduct.
	switch req.Product.Type {
	case core.DataProductType_KEY_VALUE:
		sources := []*nats.StreamSource{}
		if len(req.Product.GetSourceDataProducts()) != 0 {
			for _, source := range req.GetProduct().SourceDataProducts {
				// Must align with parents type, or nats will not process
				if source.GetType() != req.Product.Type {
					sendErrorResponse(m, "Error that not align with parent type")
					return
				}
				sources = append(sources, &nats.StreamSource{
					Name:   source.GetName(),
					Domain: source.GetDomain(),
				})
			}
		}
		s.kvManager.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:  req.Product.GetName(),
			Sources: sources,
		})

	case core.DataProductType_STREAM:
		sources := []*jetstream.StreamSource{}
		if len(req.Product.GetSourceDataProducts()) != 0 {
			for _, source := range req.GetProduct().SourceDataProducts {
				// Must align with parents type, or nats will not process
				if source.GetType() != req.Product.Type {
					sendErrorResponse(m, "Error that not align with parent type")
					return
				}
				sources = append(sources, &jetstream.StreamSource{
					Name:   source.GetName(),
					Domain: source.GetDomain(),
				})
			}
		}
		s.streamManager.CreateStream(ctx, jetstream.StreamConfig{
			Name:     req.Product.GetName(),
			Subjects: []string{"$FOREST.DP." + req.Product.GetName()},
		})

	case core.DataProductType_OBJECT:
		// Object do not have "source"
		s.objManager.CreateObjectStore(&nats.ObjectStoreConfig{
			Bucket: req.Product.GetName(),
		})

	case core.DataProductType_UNDEFINED:
		sendErrorResponse(m, "Error not define data product type")
		return
	}

	resp := &api.CreateDataProductResponse{
		Status: "OK",
	}

	data, err := proto.Marshal(resp)
	if err != nil {
		sendErrorResponse(m, "Error marshaling response")
		return
	}

	responseMsg := &nats.Msg{
		Subject: m.Reply,
		Data:    data,
	}

	m.RespondMsg(responseMsg)
}

// TODO:
func (s *Server) ReadDataProduct(m *nats.Msg) {
	req := &api.ReadDataProductRequest{}
	if err := proto.Unmarshal(m.Data, req); err != nil {
		sendErrorResponse(m, "Error unmarshaling request")
		return
	}

	// Implement logic to read a DataProduct by its name.
	// For this example, returning a dummy DataProduct.
	resp := &api.ReadDataProductResponse{
		Product: &core.DataProduct{
			Name:   req.Name,
			Domain: "Sample Domain",
			Type:   core.DataProductType_KEY_VALUE,
		},
	}
	data, err := proto.Marshal(resp)
	if err != nil {
		sendErrorResponse(m, "Error marshaling response")
		return
	}
	m.Respond(data)
}

// TODO:
func (s *Server) UpdateDataProduct(m *nats.Msg) {
	req := &api.UpdateDataProductRequest{}
	if err := proto.Unmarshal(m.Data, req); err != nil {
		sendErrorResponse(m, "Error unmarshaling request")
		return
	}

	// Implement logic to update a DataProduct.
	// For this example, returning an "OK" status.
	resp := &api.UpdateDataProductResponse{
		Status: "OK",
	}
	data, err := proto.Marshal(resp)
	if err != nil {
		sendErrorResponse(m, "Error marshaling response")
		return
	}
	m.Respond(data)
}

// TODO:
func (s *Server) DeleteDataProduct(m *nats.Msg) {
	req := &api.DeleteDataProductRequest{}
	if err := proto.Unmarshal(m.Data, req); err != nil {
		sendErrorResponse(m, "Error unmarshaling request")
		return
	}

	// Implement logic to delete a DataProduct by its name.
	// For this example, returning an "OK" status.
	resp := &api.DeleteDataProductResponse{
		Status: "OK",
	}
	data, err := proto.Marshal(resp)
	if err != nil {
		sendErrorResponse(m, "Error marshaling response")
		return
	}
	m.Respond(data)
}

func (s *Server) ListDataProducts(m *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &api.ListDataProductsRequest{}
	if err := proto.Unmarshal(m.Data, req); err != nil {
		sendErrorResponse(m, "Error unmarshaling request")
		return
	}

	// Read Data Product from NATS
	streams := s.streamManager.StreamNames(ctx)
	products := []*core.DataProduct{}

	// Stream-based data product
	for s := range streams.Name() {

		if strings.HasPrefix(s, "KV_") {
			products = append(products, &core.DataProduct{
				Name: strings.TrimPrefix(s, "KV_"),
				Type: core.DataProductType_KEY_VALUE,
			})
		} else if strings.HasPrefix(s, "OBJ_") {
			products = append(products, &core.DataProduct{
				Name: strings.TrimPrefix(s, "OBJ_"),
				Type: core.DataProductType_OBJECT,
			})
		} else {
			products = append(products, &core.DataProduct{
				Name: s,
				Type: core.DataProductType_STREAM,
			})
		}
	}
	if streams.Err() != nil {
		fmt.Println("Unexpected error ocurred")
	}

	// Process resp
	resp := &api.ListDataProductsResponse{
		Products: products,
	}
	data, err := proto.Marshal(resp)

	if err != nil {
		sendErrorResponse(m, "Error marshaling response")
		return
	}
	m.Respond(data)
}

func sendErrorResponse(m *nats.Msg, errorMsg string) {
	// Customize error responses based on your application's needs.
	m.Respond([]byte(errorMsg))
}
