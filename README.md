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
```

```
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
