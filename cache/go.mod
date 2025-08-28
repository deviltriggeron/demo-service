module cache

go 1.22

require (
	model/model v0.0.0
	postgres/postgres v0.0.0
)

replace (
	model/model => ./model
	postgres/postgres => ./postgres
)
