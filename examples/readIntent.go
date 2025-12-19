package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func generateIPv6(prefix string, subnetID int, interfID int) string {
	base := strings.TrimSuffix(prefix, "::")
	ipv6 := fmt.Sprintf("%s:%d::%d", base, subnetID, interfID)
	return ipv6
}

func main() {
	data, err := os.ReadFile("examples/json/network_intent_example.json")
	if err != nil {
		panic(err)
	}

	var intent map[string]interface{}
	err = json.Unmarshal(data, &intent)
	if err != nil {
		panic(err)
	}

	for asName, asData := range intent {
		fmt.Println(" AS:", asName)
		asMap := asData.(map[string]interface{})
		protocol := asMap["protocol"].(string)
		subnet := asMap["network_subnet"].(string)
		fmt.Println(" Protocol:", protocol)
		fmt.Println(" Subnet:", subnet)

		links := asMap["links"].(map[string]interface{})
		subnetSeenCount := make(map[int]int)
		assigned := make(map[string]map[string]int)

		for routerName, routerData := range links {
			fmt.Println(" Router:", routerName)
			routerMap := routerData.(map[string]interface{})

			if assigned[routerName] == nil {
				assigned[routerName] = make(map[string]int)
			}

			for interfName, subnetID := range routerMap {
				idStr := subnetID.(string)
				id, err := strconv.Atoi(idStr)
				if err != nil {
					panic(err)
				}

				if id != -1 {
					subnetSeenCount[id]++
					interfID := subnetSeenCount[id]
					assigned[routerName][interfName] = interfID
					ipv6 := generateIPv6(subnet, id, interfID)
					fmt.Println(" Interface:", interfName, " SubnetID:", id, " IPv6:", ipv6)
				}
			}
		}
	}
}
