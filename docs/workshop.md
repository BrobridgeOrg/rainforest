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
./rainforest_leaf --port=4111 --domain=tachun --hub-urls=localhost:7422
```
# Create a State Data Product
``` bash
nats request '$RAINFOREST.API.DP.CREATE.*' \
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
nats publish '$RAINFOREST.DP.STATE.OrdersState.0' value_0
nats publish '$RAINFOREST.DP.STATE.OrdersState.1' value_1
nats publish '$RAINFOREST.DP.STATE.OrdersState.2' value_2
nats publish '$RAINFOREST.DP.STATE.OrdersState.3' value_3
nats publish '$RAINFOREST.DP.STATE.OrdersState.4' value_4
nats publish '$RAINFOREST.DP.STATE.OrdersState.5' value_5
nats publish '$RAINFOREST.DP.STATE.OrdersState.6' value_6
nats publish '$RAINFOREST.DP.STATE.OrdersState.7' value_7
nats publish '$RAINFOREST.DP.STATE.OrdersState.8' value_8
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
./rainforest_leaf --hub-server=URL --user=USER --password=PASSWORD --domain= 
```

# Create a Data Product. but ... a "Source" Data Product!
```

```
# I Stop my Data Product


# Read it on your own replica!


# I start my Rainforest server

# Write some data to Data Product

# Read from your Data Product! 