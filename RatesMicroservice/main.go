package main

import (
	"RatesMicroservice/server"
	_ "github.com/lib/pq"
)

func main() {
	server.Start()
}
