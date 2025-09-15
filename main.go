package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/salahfarzin/roja-shop/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	server.Run()
}
