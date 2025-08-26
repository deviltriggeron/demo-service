module service

go 1.24.2

require (
	postgres/postgres v0.0.0
	cache/cache v0.0.0
	model/model v0.0.0
)

replace (
	postgres/postgres => ./postgres
	cache/cache => ./cache
	model/model => ./model
)
