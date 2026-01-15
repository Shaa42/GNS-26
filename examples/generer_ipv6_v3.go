package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func generateIPv6(prefix string, subnetID int, ifaceID int) string {
	base := strings.TrimSuffix(prefix, "::")
	return fmt.Sprintf("%s:%d::%d", base, subnetID, ifaceID)
}

func main() {
	// lecture du intent
	data, err := os.ReadFile("examples/json/network_intent_new.json")
	if err != nil {
		panic(err)
	}

	var intent map[string]interface{}
	if err := json.Unmarshal(data, &intent); err != nil {
		panic(err)
	}

	// parcourrir l'AS
	for asName, asData := range intent {
		fmt.Println("AS:", asName)

		asMap := asData.(map[string]interface{})
		prefix := asMap["network_subnet"].(string)

		links := asMap["links"].([]interface{})

		// links
		for _, linkRaw := range links {
			link := linkRaw.(map[string]interface{})

			subnetRaw, ok := link["subnet"]
			if !ok {
				continue
			}
			subnetID := int(subnetRaw.(float64))

			endpoints := link["endpoints"].(map[string]interface{})

			ifaceID := 1
			for router, ifaceRaw := range endpoints {
				iface := ifaceRaw.(string)

				ipv6 := generateIPv6(prefix, subnetID, ifaceID)

				fmt.Printf(
					"Router: %s  Interface: %s  Subnet: %d  IPv6: %s\n",
					router, iface, subnetID, ipv6,
				)

				ifaceID++
			}
		}
	}
}
