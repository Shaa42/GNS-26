package main

import (
	"fmt"
	"gns-26/applycfg"
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
	as112 := asMap["AS112"]
	// Create a .cfg file for router named R1 with data dict
	writeconf.WriteConfig(as111)
	as111.LogAS()
	writeconf.WriteConfig(as112)
	as112.LogAS()

	// Apply ths just created .cfg to every router
	for _, router := range as111.Routers {
		applycfg.ApplyCfg(router.Name[1:], "GNS-project")
	}

	for _, router := range as112.Routers {
		applycfg.ApplyCfg(router.Name[1:], "GNS-project")
	}

	fmt.Println("-- All .cfg applied.")
}
