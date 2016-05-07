#!/bin/bash

export PIKA_PORT=3000
export PIKA_DATABASE_URL="postgres://postgres:postgres@localhost:5432/pika?sslmode=disable"

go run main.go
