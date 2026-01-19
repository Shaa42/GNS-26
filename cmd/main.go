package main

import (
	"gns-26/internal/parseintent"
	"gns-26/internal/writeconf"
)

func main() {

	// Parse json intent file and print the as' info
	asMap, err := parseintent.NewAS("examples/json/network_intent_template_v4.json")
	if err != nil {
		panic(err)
	}
	as111 := asMap["AS111"]
	// Create a .cfg file for router named R1 with data dict
	writeconf.WriteConfig(as111)
	as111.LogAS()
}
