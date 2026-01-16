package parseintent

import (
	"encoding/json"
	"fmt"
	"os"
)

type InfoAS struct {
	Name          string            `json:"name"`
	Protocol      string            `json:"protocol"`
	AddressFamily string            `json:"address_family"`
	NetworkSubnet string            `json:"network_subnet"`
	RemoteAS      map[string]string `json:"remote_as"`
	Routers       []InfoRouter      `json:"routers"`
	Links         []InfoLink        `json:"links"`
}

type InfoRouter struct {
	Name       string                   `json:"name"`
	RouterID   string                   `json:"router_id"`
	Interfaces map[string]InfoInterface `json:"interfaces"`
}

type InfoInterface struct {
	Role string `json:"role"`
}

type InfoLink struct {
	Endpoints map[string]string `json:"endpoints"`
	Role      string            `json:"role"`
	Info      map[string]string `json:"info"`
}

func NewAS(path string) (map[string]InfoAS, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("test 1")
		return nil, err
	}

	var payload map[string]InfoAS
	err = json.Unmarshal(content, &payload)
	if err != nil {
		fmt.Println("test 2")
		return nil, err
	}

	return payload, nil
}

func (as InfoAS) LogAS() {
	fmt.Printf("Name: %s\n", as.Name)
	fmt.Printf("Protocol: %s\n", as.Protocol)
	fmt.Printf("Address Family: %s\n", as.AddressFamily)
	fmt.Printf("Network Subnet: %s\n", as.NetworkSubnet)

	fmt.Println("\n--- Remote AS ---")
	for asName, relation := range as.RemoteAS {
		fmt.Printf("  %s: %s\n", asName, relation)
	}

	fmt.Println("\n--- Routers ---")
	for _, router := range as.Routers {
		fmt.Printf("  Router: %s (ID: %s)\n", router.Name, router.RouterID)
		for ifName, ifInfo := range router.Interfaces {
			fmt.Printf("    - %s: role=%s\n", ifName, ifInfo.Role)
		}
	}

	fmt.Println("\n--- Links ---")
	for i, link := range as.Links {
		fmt.Printf("  Link %d: role=%s\n", i+1, link.Role)
		fmt.Printf("    Endpoints:\n")
		for router, iface := range link.Endpoints {
			fmt.Printf("      %s: %s\n", router, iface)
		}
		fmt.Printf("    Info: %v\n", link.Info)
	}
}
