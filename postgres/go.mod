module postgres

go 1.24.2

require (
	github.com/lib/pq v1.10.9
	model/model v0.0.0
)

replace model/model => ./model
