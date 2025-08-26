module kafka

go 1.24.2

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/segmentio/kafka-go v0.4.48 // indirect
	model/model v0.0.0
	postgres/postgres v0.0.0
	cache/cache v0.0.0
	service/service v0.0.0
)

replace (
	model/model => ./model
	postgres/postgres => ./postgres
	cache/cache => ./cache
	service/service => ./service
)
