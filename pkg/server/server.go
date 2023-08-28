package server

import (
	"context"
	"errors"
	"log"
	"time"

	apiv1 "github.com/Awareness-Labs/rainforest/pkg/proto/api/v1"
	corev1 "github.com/Awareness-Labs/rainforest/pkg/proto/core/v1"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	StateDataProductSubjectPrefix  = "$RAINFOREST.DP.STATE."
	EventDataProductSubjectPrefix  = "$RAINFOREST.DP.EVENT."
	SourceDataProductSubjectPrefix = "$RAINFOREST.DP.SOURCE."

	StateDataProductPrefix  = "STATE_"
	EventDataProductPrefix  = "EVENT_"
	SourceDataProductPrefix = "SOURCE_"
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

	// API Subscribe to DataProduct topics
	_, err := s.conn.Subscribe("$RAINFOREST.API.DP.CREATE.*", s.CreateDataProduct)
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	_, err = s.conn.Subscribe("$RAINFOREST.API.DP.INFO.*", s.InfoDataProduct)
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	_, err = s.conn.Subscribe("$RAINFOREST.API.DP.UPDATE.*", s.UpdateDataProduct)
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	_, err = s.conn.Subscribe("$RAINFOREST.API.DP.DELETE.*", s.DeleteDataProduct)
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

	// Unmarshal
	req := &apiv1.CreateDataProductRequest{}

	err := protojson.Unmarshal(m.Data, req)

	if err != nil {
		sendErrorResponse(m, err)
		return
	}

	// TODO: Validate req
	// if err != nil {
	// 	sendErrorResponse(m, err)
	// 	return
	// }

	// Implement create a DataProduct.
	newDp := req.GetProduct()

	switch newDp.GetType() {
	case corev1.DataProductType_DATA_PRODUCT_TYPE_STATE:
		// if this DP has source
		sources := getSourceDataProduct(newDp)
		_, err := s.streamManager.CreateStream(ctx, jetstream.StreamConfig{
			Name:              StateDataProductPrefix + newDp.GetName(),
			MaxMsgsPerSubject: 1,
			Subjects:          []string{StateDataProductSubjectPrefix + newDp.GetName() + ".*"},
			Description:       newDp.GetDescription(),
			Discard:           jetstream.DiscardOld,
			Sources:           sources,
		})
		if err != nil {
			sendErrorResponse(m, err)
			return
		}
		sendResponse(m, "data product created: "+newDp.GetName())
	case corev1.DataProductType_DATA_PRODUCT_TYPE_EVENT:
		// if this DP has source
		sources := getSourceDataProduct(newDp)
		_, err := s.streamManager.CreateStream(ctx, jetstream.StreamConfig{
			Name:        EventDataProductPrefix + newDp.GetName(),
			Subjects:    []string{EventDataProductSubjectPrefix + newDp.GetName() + ".*"},
			Description: newDp.GetDescription(),
			Sources:     sources,
		})
		if err != nil {
			sendErrorResponse(m, err)
			return
		}
		sendResponse(m, "data product created: "+newDp.GetName())

	default:
		sendErrorResponse(m, errors.New("not support this type data product"))
		return
	}

}

func sendErrorResponse(m *nats.Msg, err error) {
	res := &apiv1.ErrorResponse{
		ErrorCode:    "400",
		ErrorMessage: err.Error(),
	}

	resData, _ := protojson.Marshal(res)
	m.Respond(resData)
}

func sendResponse(m *nats.Msg, status string) {
	// res := &apiv1.DataProductResponse{}
	// resData, _ := json.Marshal(res)
	m.Respond([]byte(status))
}

func getSourceDataProduct(newDp *corev1.DataProduct) []*jetstream.StreamSource {
	sources := []*jetstream.StreamSource{}
	for _, dp := range newDp.SourceDataProducts {
		// Must align with parents type, or nats will not process
		sources = append(sources, &jetstream.StreamSource{
			Name:   dp.Name,
			Domain: dp.Domain,
		})
		log.Println("SOURCE:", dp.Name, dp.Domain)
	}
	return sources
}

// Not MVP feature!

// TODO:
func (s *Server) InfoDataProduct(m *nats.Msg) {
}

// TODO:
func (s *Server) UpdateDataProduct(m *nats.Msg) {
}

// TODO:
func (s *Server) DeleteDataProduct(m *nats.Msg) {
}

// TODO
func (s *Server) ListDataProducts(m *nats.Msg) {
}
