package main

import (
	"log"

	"github.com/s1ovac/wb-service/internal/publish"
)

func main() {
	pb := publish.New()
	if err := pb.DropMessage(); err != nil {
		log.Fatal(err)
	}
}
