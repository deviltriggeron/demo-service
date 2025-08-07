module main

go 1.24.2

require postgres/postgres v0.0.0

require github.com/lib/pq v1.10.9 // indirect

replace postgres/postgres => ./postgres
