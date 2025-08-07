#!/bin/bash

brew services start postgresql
psql -U postgres -d postgres -h localhost -f postgres/model.sql
go build .
./main