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
go run cmd/create_source_kv_dp/create_source_kv_dp.go 
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
