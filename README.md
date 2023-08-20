# rainforest

## Archtecture


## Data Product is first class citizen
* Data Product can build another Data Product
* Application level design consume data product

## Data Product Service
* DataProduct -> Service

## Data Product Adapter
* Data Product Adapter -> Data Product

### Data Product Type
* Source-Aligned
    * KV-based  
    * Obj-based
    * Stream-based
* Aggregate-Aligned/Consumer-Aligned
    * Multi-KV
    * Multi-Stream


## Labs: Quickstart!
* Level-0: Basic Concept
    * Create KV DataProduct
    * Create Stream DataProdcut
    * Create Object DataProduct
* Level-1: Advance Concept
    * Create Source KV DataProduct
    * Create Source Stream DataProdcut
* Level-2: 
    * Create Pipeline DataProduct
    * Create Embedded Logic DataProduct

## Demo: Create Data Product(stream-based)
SQL Database/Application Behavier/Event -> CDC Event -> stream operation 
```
nats server ls --server nats://rainforest_user:password@localhost:4111

nats --server nats://rainforest_user:password@localhost:4111 stream add Conversations --subjects='$RAINFOREST.DATAPRODUCT.Conversations'
nats --server nats://rainforest_user:password@localhost:4112 stream add Orders --subjects='$RAINFOREST.DATAPRODUCT.Orders'
```

## Demo: Create Data Product(kv-based)
SQL Database/Application Behavier/Event -> CDC Event -> kv operation (There is only one kind of adapter!) 
SQL Database/Application Behavier/Event -> CDC Event -> transaction operation (There is only one kind of adapter!) 


## Demo: Create Data Product(object-based)
```
nats publish '$RAINFOREST.DATAPRODUCT.StreamDataProductExample' test --server rainforest_user:password@localhost:4111
```

## Demo check Data Product
```
nats stream ls --server rainforest_user:password@localhost:4111
nats object ls --server rainforest_user:password@localhost:4111
nats kv ls --server rainforest_user:password@localhost:4111
```

## Demo automation
User(Behavier) -> Operation -> Domain "State" -> Discovery & Replication

To another one consumer

User(Behavier) -> Operation -> Domain "State" -> Discovery & Replication

* Application example: Event -> Object -> OLAP

* Application example: KV -> Distributed Cache

* Application example: KV -> Secondary Index, GIS Index, Fulltext Index

* Application example: Event -> LLM(API Call)


# Demo for Brobridge 
## Key-Value Demo
### Start 1 Hub 2 Leaf (You can add Leaf-Server as you want)
* rainforest_leaf_0
* rainforest_leaf_1

```
cd cluster
goreman start
```

### Create (KV)Data Product (From 4111 host)
```
go run cmd/create_kv_dp/create_kv_dp.go
```

```
nats kv ls --server rainforest_user:password@localhost:4111
```
### Create Secondary (KV)Data Product (From 4112 host)
```


```

```
nats kv ls --server rainforest_user:password@localhost:4112
```
### Write Key/Value on (KV)Data Product (From 4111 host)
```
nats kv put KVDataProductExample key_0 value_0 --server rainforest_user:password@localhost:4111
nats kv put KVDataProductExample key_1 value_1 --server rainforest_user:password@localhost:4111
nats kv put KVDataProductExample key_2 value_2 --server rainforest_user:password@localhost:4111
```

### Read Key/Value from Secondary (KV)Data Product (From 4112 host)
```
nats kv ls --server rainforest_user:password@localhost:4112
```

## Stream Demo
### Start 1 Hub 2 Leaf (You can add Leaf-Server as you want)
* rainforest_leaf_0
* rainforest_leaf_1

```
cd cluster
goreman start
```
### Create (Stream)Data Product (From 4111 host)
```
go run cmd/create_stream_dp/create_stream_dp.go
```

```
nats s ls --server rainforest_user:password@localhost:4111
```
### Create Secondary (Stream)Data Product (From 4112 host)
```
go run cmd/create_source_stream_dp/create_source_stream_dp.go
```

```
nats stream ls --server rainforest_user:password@localhost:4112
```
### Write Event on (Stream)Data Product (From 4111 host)
```
nats publish '$RAINFOREST.DATAPRODUCT.StreamDataProductExample' event_0 --server rainforest_user:password@localhost:4111
nats publish '$RAINFOREST.DATAPRODUCT.StreamDataProductExample' event_1 --server rainforest_user:password@localhost:4111
nats publish '$RAINFOREST.DATAPRODUCT.StreamDataProductExample' event_2 --server rainforest_user:password@localhost:4111
```

### Read Event from Secondary (Stream)Data Product (From 4112 host)
```
nats s ls --server rainforest_user:password@localhost:4112

nats s get --server rainforest_user:password@localhost:4112
```

# Art of Stream
## Event 
Just like normal event stream nothing special
## State Stream
``` bash
$ nats s add StateStream -s rainforest_user:password@localhost:4111
? Subjects $RAINFOREST.DataProduct.>
? Per Subject Messages Limit 1
```

### Publish 3 Event for entities
```
nats publish '$RAINFOREST.DataProduct.0' Recored_0 -s rainforest_user:password@localhost:4111
nats publish '$RAINFOREST.DataProduct.1' Recored_1 -s rainforest_user:password@localhost:4111
nats publish '$RAINFOREST.DataProduct.2' Recored_2 -s rainforest_user:password@localhost:4111
```

### Publish 1 After-Event 
```
nats publish '$RAINFOREST.DataProduct.0' Recored_0_new -s rainforest_user:password@localhost:4111
```

### Subscrible all of them (that will show Recored_0_new!!)
```
nats subscribe '$RAINFOREST.DataProduct.0' --start-sequence=1 -s rainforest_user:password@localhost:4111
```


# Build-In Worker
* Storage: State Stream -> BadgerDB 
* Fulltext: State Stream -> Bleve
* Sink: State Stream -> Object (Run DuckDB)

# API
$RAINFOREST.API.DATAPRODUCT.CREATE.*
$RAINFOREST.API.DATAPRODUCT.INFO.*
$RAINFOREST.API.DATAPRODUCT.UPDATE.*
$RAINFOREST.API.DATAPRODUCT.DELETE.*
$RAINFOREST.API.DATAPRODUCT.LIST

$RAINFOREST.API.DATAPRODUCT.<data_product>.<pk>
$RAINFOREST.API.DATAPRODUCT.<data_product>.<event_id>

$RAINFOREST.API.KV.<data_product> 
$RAINFOREST.API.FULLTEXT.<data_product>
$RAINFOREST.API.DUCK.<data_product>



# Presentation
1. Rainforest Intro
    * Span States, Events to Unlimited Team
    * Infra aligned
2. Demo I: Create State Data Product
3. Demo I: Create Source State Data Product
4. Demo I: KV API
5. Demo I: FULLTEXT API
5. Demo I: DUCK API

# Create State Data Prodcut (I just build for you)
```
nats request '$RAINFOREST.API.DATAPRODUCT.CREATE.*' \
'{ 
  "name": "StateDataProductExample", 
  "domain": "Team A", 
  "dataproduct_type": "StateDataProduct", 
  "description": "This is a sample state data product."
}'
```

# Write some State data
```
nats publish '$RAINFOREST.DP.STATE.StateDataProductExample.0' value_0
nats publish '$RAINFOREST.DP.STATE.StateDataProductExample.1' value_1
nats publish '$RAINFOREST.DP.STATE.StateDataProductExample.2' value_2
```
# Write some Event data
```
nats publish '$RAINFOREST.DP.EVENT.EventDataProductExample.0' click
nats publish '$RAINFOREST.DP.EVENT.EventDataProductExample.1' roll
nats publish '$RAINFOREST.DP.EVENT.EventDataProductExample.2' open
nats publish '$RAINFOREST.DP.EVENT.EventDataProductExample.2' download
nats publish '$RAINFOREST.DP.EVENT.EventDataProductExample.2' close
```

# Create Event Data Prodcut
```
nats request '$RAINFOREST.API.DATAPRODUCT.CREATE.*' 
'{ 
  "name": "EventDataProductExample", 
  "domain": "Team A", 
  "dataproduct_type": "EventDataProduct", 
  "description": "This is a sample event data product."
}'

nats publish '$RAINFOREST.API.DATAPRODUCT.EventDataProductExample.0' value_0
nats publish '$RAINFOREST.API.DATAPRODUCT.EventDataProductExample.0' value_1
nats publish '$RAINFOREST.API.DATAPRODUCT.EventDataProductExample.0' value_2
```

# Create Source Data Product
```
nats request '$RAINFOREST.API.DATAPRODUCT.CREATE.*' \
'{
  "name": "SourceStateDataProduct",
  "domain": "TeamB",
  "dataproduct_type": "StateDataProduct",
  "description": "This is a sample source state data product.",
  "sources": [
    {
      "name": "SourceStateDataProduct",
      "domain": "SourceDomain",
      "dataproduct_type": "StateDataProduct"
    }
  ]
}'


```


# Leaf server 
* leaf port: 4222

# Hub Server
* port: 4222
* leaf: 7422