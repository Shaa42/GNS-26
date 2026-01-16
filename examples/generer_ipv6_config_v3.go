package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func generateIPv6(prefix string, subnetID int, ifaceID int) string {
	parts := strings.Split(prefix, "/")
	base := strings.TrimSuffix(parts[0], "::")
	return fmt.Sprintf("%s:%d::%d/64", base, subnetID, ifaceID)
}

func generateLoopback(prefix string, routerName string) string {
	parts := strings.Split(prefix, "/")
	base := strings.TrimSuffix(parts[0], "::")
	id := routerName[1:]
	return fmt.Sprintf("%s:100::%s/128", base, id)
}

func renderInterfaceProtocol(protocol string) string {
	switch protocol {
	case "RIP":
		return " ipv6 rip RIPng enable\n"
	case "OSPF":
		return " ipv6 ospf 1 area 0\n"
	default:
		return ""
	}
}

func renderBGP(localAS int, routerID string, peerAS int, peerIPv6 string) string {
	cfg := ""

	cfg += fmt.Sprintf("router bgp %d\n", localAS)
	cfg += fmt.Sprintf(" bgp router-id %s\n", routerID)
	cfg += " bgp log-neighbor-changes\n"
	cfg += fmt.Sprintf(" neighbor %s remote-as %d\n", peerIPv6, peerAS)
	cfg += "address-family ipv6 unicast\n"
	cfg += fmt.Sprintf(" neighbor %s activate\n", peerIPv6)
	cfg += "exit-address-family"

	//cfg += fmt.Sprintf(" neighbor %s update-source Loopback0\n", peerIPv6)
	cfg += "\n"

	return cfg
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
		protocol := asMap["protocol"].(string)
		prefix := asMap["network_subnet"].(string)
		links := asMap["links"].([]interface{})
		routers := asMap["routers"].(map[string]interface{})
		routerIDs := make(map[string]string)

		// loopback ipv6
		for routerName, routerData := range routers {
			routerMap := routerData.(map[string]interface{})
			routerIDs[routerName] = routerMap["router_id"].(string)

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
			role := "internal"
			value := link["role"]
			if value != nil {
				role = value.(string)
			}

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

				if role == "internal" {
					cfg[router] += fmt.Sprintf(
						"interface %s\n ipv6 address %s\n%s no shutdown\n\n",
						iface,
						ipv6,
						renderInterfaceProtocol(protocol),
					)
				}
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

		//handle eBGP
		for _, linkRaw := range links {
			link := linkRaw.(map[string]interface{})

			roleValue := link["role"]
			if roleValue == nil || roleValue.(string) != "eBGP" {
				continue
			}

			peer := link["peer"].(map[string]interface{})
			peerASstr := peer["as"].(string)
			peerIPv6 := peer["ipv6"].(string)

			peerAS, _ := strconv.Atoi(peerASstr[2:])
			localAS, _ := strconv.Atoi(asName[2:])
			endpoints := link["endpoints"].(map[string]interface{})

			for routerName := range endpoints {
				routerID := routerIDs[routerName]
				cfg[routerName] += renderBGP(localAS, routerID, peerAS, peerIPv6)
			}
		}

		for routerName, content := range cfg {
			os.WriteFile(routerName+".cfg", []byte(content), 0644)
		}
	}
}
