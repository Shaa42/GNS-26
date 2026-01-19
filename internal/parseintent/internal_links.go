package parseintent

import (
	"strconv"
)

type InterfacePlan struct {
	Name      string
	IPv6      string
	EnableIGP bool
}

func planInternalLinks(
	links []interface{},
	prefix string,
) map[string][]InterfacePlan {

	res := make(map[string][]InterfacePlan)

	for _, l := range links {
		link := l.(map[string]interface{})

		role, ok := link["role"].(string)
		if !ok || role != "internal" {
			// traiter seulement les interfaces "internal"
			continue
		}

		info := link["info"].(map[string]interface{})

		subnetStr := info["subnet"].(string)
		subnetID, _ := strconv.Atoi(subnetStr)

		endpoints := link["endpoints"].(map[string]interface{})

		ifaceID := 1
		for router, ifaceRaw := range endpoints {
			iface := ifaceRaw.(string)

			ipv6 := generateIPv6(prefix, subnetID, ifaceID)

			res[router] = append(res[router], InterfacePlan{
				Name:      iface,
				IPv6:      ipv6,
				EnableIGP: true,
			})

			ifaceID++
		}
	}

	return res
}

func planLoopbacks(prefix string, routers []interface{}) map[string]InterfacePlan {

	res := make(map[string]InterfacePlan)

	for _, r := range routers {
		router := r.(map[string]interface{})

		name := router["name"].(string)

		res[name] = InterfacePlan{
			Name:      "Loopback0",
			IPv6:      generateLoopback(prefix, name),
			EnableIGP: true,
		}
	}

	return res
}
