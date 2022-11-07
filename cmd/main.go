package main

import "github.com/s1ovac/wb-service/internal/publish"

func main() {
	pb := publish.New()
	sc := publish.InitConnect(pb)
	pb.DropMessage(sc)

}
