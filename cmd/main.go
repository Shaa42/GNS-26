package main

import (
	"gns-26/internal/applycfg"
	"gns-26/internal/parseintent"
	"gns-26/internal/writeconf"
)

func main() {

	// Parse json intent file and print the AS' info
	asMap, err := parseintent.NewAS("networks/network_intent_2.json")
	if err != nil {
		panic(err)
	}

	for _, asVal := range asMap {
		// Create a .cfg file for every AS' router
		writeconf.WriteConfig(asVal)
	}

	for _, asVal := range asMap {
		// asVal.LogAS()
		for _, router := range asVal.Routers {
			// Apply the created .cfg to every router
			applycfg.ApplyCfg(router.Name[1:], "big-nw")
		}
	}
}
