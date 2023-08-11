package main

import (
	"log"
	"time"

	api "github.com/Awareness-Labs/rainforest/pkg/proto/api/v1"
	core "github.com/Awareness-Labs/rainforest/pkg/proto/core/v1"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func main() {

	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	CreateDataProduct(nc)
	ListDataProducs(nc)
	// Feature: Create source Data Product
}

func CreateDataProduct(nc *nats.Conn) {
	// Feature: Create KV Data Product
	req := &api.CreateDataProductRequest{
		Product: &core.DataProduct{
			Name: "ExampleKeyValueDataProduct",
			Type: core.DataProductType_KEY_VALUE,
		},
	}
	data, _ := proto.Marshal(req)
	msg := &nats.Msg{
		Subject: "$RAINFOREST.API.DP.CREATE",
		Data:    data,
	}
	res, _ := nc.RequestMsg(msg, 1*time.Second)

	log.Println(string(res.Data))

	// Feature: Create Stream Data Product
	req = &api.CreateDataProductRequest{
		Product: &core.DataProduct{
			Name: "ExampleStreamDataProduct",
			Type: core.DataProductType_STREAM,
		},
	}
	data, _ = proto.Marshal(req)
	msg = &nats.Msg{
		Subject: "$RAINFOREST.API.DP.CREATE",
		Data:    data,
	}
	res, _ = nc.RequestMsg(msg, 1*time.Second)

	log.Println(string(res.Data))

	// Feature: Create Object Data Product
	req = &api.CreateDataProductRequest{
		Product: &core.DataProduct{
			Name: "ExampleObjectDataProduct",
			Type: core.DataProductType_OBJECT,
		},
	}
	data, _ = proto.Marshal(req)
	msg = &nats.Msg{
		Subject: "$RAINFOREST.API.DP.CREATE",
		Data:    data,
	}
	res, _ = nc.RequestMsg(msg, 1*time.Second)

	log.Println(string(res.Data))

}

func ListDataProducs(nc *nats.Conn) {
	// Feature: List all Data Product
	req := &api.ListDataProductsRequest{}
	data, _ := proto.Marshal(req)
	msg := &nats.Msg{
		Subject: "$RAINFOREST.API.DP.LIST",
		Data:    data,
	}
	msg, _ = nc.RequestMsg(msg, 1*time.Second)
	res := &api.ListDataProductsResponse{}
	proto.Unmarshal(msg.Data, res)
	for _, product := range res.Products {
		log.Println("Product Name:", product.GetName(), "Product Type:", product.GetType())
	}

}
