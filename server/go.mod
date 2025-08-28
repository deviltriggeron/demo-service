module server

go 1.22

require (
	service/service v0.0.0
	model/model v0.0.0
)

replace (
	service/service => ./service
	model/model => ./model
)

require github.com/gorilla/mux v1.8.1 // indirect
