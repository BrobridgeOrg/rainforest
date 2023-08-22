# Workshop

# Start your Rainforest Hub server
```
cd cmd/rainforest_hub
goreman start
```
* Client connections:   4222
* Leafnode connections: 7422

# Start your own Rainforest Leaf server
```
go run cmd/rainforest_leaf/rainforest_leaf.go --port=4111 --domain=tachun --hub-urls=localhost:7422 --kv-path=./data/badger/sts-0 --stream-path=./data/stream/sts-0
```
# Create a State Data Product
``` bash
nats request '$RAINFOREST.API.DP.CREATE.*' --server=localhost:4111 \
'{
  "product": {
    "name": "OrdersState",
    "domain": "ProductTeam",
    "type": "DATA_PRODUCT_TYPE_STATE",
    "description": "This is a State Data Product"
  }
}'
```
# Publish some State
```
nats publish '$RAINFOREST.DP.STATE.OrdersState.0' value_0 --server=localhost:4111
nats publish '$RAINFOREST.DP.STATE.OrdersState.1' value_1 --server=localhost:4111
nats publish '$RAINFOREST.DP.STATE.OrdersState.2' value_2 --server=localhost:4111
nats publish '$RAINFOREST.DP.STATE.OrdersState.3' value_3 --server=localhost:4111
nats publish '$RAINFOREST.DP.STATE.OrdersState.4' value_4 --server=localhost:4111
nats publish '$RAINFOREST.DP.STATE.OrdersState.5' value_5 --server=localhost:4111
nats publish '$RAINFOREST.DP.STATE.OrdersState.6' value_6 --server=localhost:4111
nats publish '$RAINFOREST.DP.STATE.OrdersState.7' value_7 --server=localhost:4111
nats publish '$RAINFOREST.DP.STATE.OrdersState.8' value_8 --server=localhost:4111
```
# Read it!
```
nats subscribe '$RAINFOREST.DP.STATE.OrdersState.3' --last

nats request '$RAINFOREST.API.KV.*' \
'{
  "scan": {
    "limit": 10,
    "reverse": false,
    "start_key": "OrdersState.3",
    "end_key": "OrdersState.8"
  }
}
'
```

# Create another Rainforest server
```
go run cmd/rainforest_leaf/rainforest_leaf.go --port=4112 --domain=prod --hub-urls=localhost:7422 --kv-path=./data/badger/sts-1 --stream-path=./data/stream/sts-1
```

# Create a Data Product. but ... a "Source" Data Product!
```
nats request '$RAINFOREST.API.DP.CREATE.*' --server=localhost:4112 \
'{
  "product": {
    "name": "SecondaryDataProduct",
    "domain": "tachun",
    "type": "DATA_PRODUCT_TYPE_SOURCE",
    "description": "This is a Source Data Product",
    "source_data_products": [
      {
        "name": "STATE_OrdersState",
        "domain": "tachun"
      }
    ]
  }
}'
```
# I Stop my Data Product


# Read it on your own replica!


# I start my Rainforest server

# Write some data to Data Product

# Read from your Data Product! 