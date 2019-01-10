package main

import (
	"fmt"
	"github.com/alexisvisco/ovh-domain-api/domain"
	"github.com/alexisvisco/ovh-domain-api/domain/subsidiary"
)

func main() {
	client, e := domain.NewClient(subsidiary.EU)
	if e != nil {
		panic(e)
	}
	results, e := client.DomainInfo("google.fr")
	if e != nil {
		if e == domain.UnknownExtension {
			fmt.Println("UnknownExtension")
		} else {
			panic(e)
		}
	}
	fmt.Println(results.IsTaken())
}
