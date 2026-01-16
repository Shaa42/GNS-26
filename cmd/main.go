package main

import (
	"gns-26/internal/parseintent"
	"gns-26/internal/writeconf"
)

func main() {
	as112 := map[string]map[string]string{
		"R4": {
			"GigabitEthernet 1/0": "-1",
			"GigabitEthernet 2/0": "2001:2:100::1/64",
			"GigabitEthernet 3/0": "-1",
		},

		"R5": {
			"GigabitEthernet 1/0": "2001:2:100::2/64",
			"GigabitEthernet 2/0": "2001:2:101::1/64",
			"GigabitEthernet 3/0": "-1",
		},

		"R6": {
			"GigabitEthernet 1/0": "2001:2:101::2/64",
			"GigabitEthernet 2/0": "-1",
			"GigabitEthernet 3/0": "-1",
		},
	}

	// Create a .cfg file for router named R1 with data dict
	writeconf.WriteConfig("R6", as112, "OSPF")

	// Parse json intent file and print the as' info
	asMap, err := parseintent.NewAS("examples/json/network_intent_template_v4.json")
	if err != nil {
		panic(err)
	}
	as111 := asMap["AS111"]
	as111.LogAS()
}
