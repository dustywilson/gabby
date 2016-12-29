package main

import (
	"fmt"
	"log"

	"github.com/dustywilson/gabby/peer"
)

func main() {
	// g := gabber.New(domain.New("wilson.farm"), nil)
	// n, err := g.Write([]byte("Fred!"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("We sent %d bytes.\n", n)
	p, err := peer.Parse("dusty@wilson.farm")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(p.String())
	fmt.Println(p.FetchKey())
	fmt.Println(p.Key().Expired())
	fmt.Println(p.Key().ShouldRefresh())
}
