module main

go 1.22

require (
	cache/cache v0.0.0
	kafka/kafka v0.0.0
	postgres/postgres v0.0.0
	server/server v0.0.0
	service/service v0.0.0
)

require (
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/segmentio/kafka-go v0.4.48 // indirect
	model/model v0.0.0 // indirect
)

replace (
	cache/cache => ./cache
	kafka/kafka => ./kafka
	model/model => ./model
	postgres/postgres => ./postgres
	server/server => ./server
	service/service => ./service
)
