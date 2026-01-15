package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func generateIPv6(prefix string, subnetID int, ifaceID int) string {
	base := strings.TrimSuffix(prefix, "::")
	return fmt.Sprintf("%s:%d::%d/64", base, subnetID, ifaceID)
}

func generateLoopback(prefix string, routerName string) string {
	base := strings.TrimSuffix(prefix, "::")
	id := routerName[1:]
	return fmt.Sprintf("%s:100::%s/128", base, id)
}

func main() {
	// lecture du intent
	data, err := os.ReadFile("json/network_intent_template_v3.json")
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
		cfg := make(map[string]string)

		asMap := asData.(map[string]interface{})
		prefix := asMap["network_subnet"].(string)

		links := asMap["links"].([]interface{})
		routers := asMap["routers"].(map[string]interface{})

		// loopback ipv6
		for routerName := range routers {
			cfg[routerName] = "ipv6 unicast-routing\n\n"
			loopbackIPv6 := generateLoopback(prefix, routerName)
			cfg[routerName] += fmt.Sprintf(
				"interface Loopback0\n ipv6 address %s\n\n",
				loopbackIPv6,
			)
		}

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

				cfg[router] += fmt.Sprintf(
					"interface %s\n ipv6 address %s\n ipv6 rip RIPng enable\n no shutdown\n\n",
					iface, ipv6,
				)

				ifaceID++
			}
		}

		for routerName, content := range cfg {
			filename := routerName + ".cfg"
			err := os.WriteFile(filename, []byte(content), 0644)
			if err != nil {
				panic(err)
			}
			fmt.Println("Generated", filename)
		}
	}
}
