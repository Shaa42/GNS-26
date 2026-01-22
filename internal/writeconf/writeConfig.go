package writeconf

import (
	"fmt"
	"gns-26/internal/parseintent"
	"log"
	"os"
	"strings"
)

func WriteConfig(data parseintent.InfoAS) {
	eBGP_activated := false // If eBGP is activated, thus activating iBGP
	is_eBGP_router := false
	
	var loopbackNeighbours []string
	if len(data.RemoteAS) > 0 {
		// More than one AS in the network
		eBGP_activated = true
	}
	if eBGP_activated {
		// This first loop is for finding the iBGP neighbors
		for _, router := range data.Routers {
			interfaces := router.Interfaces

			for interfaceName, interfaceInfo := range interfaces {
				if interfaceInfo.Role == "none" {
					continue
				}
				interfaceStr := "interface"
				interfaceStr += " "
				interfaceStr += interfaceName
				interfaceIP := ""
				
				if interfaceInfo.Role == "loopback" {
					// If the interface is a loopback and eBGP is activated, 
					// set up the loopback and add it to the internal_neighbor slice
					if eBGP_activated {
						interfaceIP = generateLoopbackIPv6(data.NetworkSubnet, router.Name)
						
						loopbackNeighbours = append(loopbackNeighbours, interfaceIP)
					}
					
				}
			}
		}
	}

	

	// Setup the .cfg file for every router
	for _, router := range data.Routers {
		is_eBGP_router = false
		// R1 => rN = "1"
		rN := router.Name[1:]
		FILENAME := router.Name + "_configs_i" + rN + "_startup-config" + ".cfg"
		confInterfaceStr := ""
		ripConfStr := ""
		ospfConfStr := "!\n!"

		// Depending on the protocol, will add strings related specifically to that protocol
		switch data.Protocol {
		case "RIPng":
			confInterfaceStr = "ipv6 rip rip_process enable"
			ripConfStr += "\nipv6 router rip rip_process"
		case "OSPF":
			confInterfaceStr = "ipv6 ospf 1 area 0" // WARNING: currently CANNOT handle multiple areas
			ospfConfStr += "\nipv6 router ospf 1"
			ospfConfStr += "\n router-id "
			ospfConfStr += router.RouterID

		default:
			panic("unrecognized internal routing protocol (atm only RIP and OSPF can be used)")
		}
		
		file, err := os.Create(FILENAME)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		interfaces := router.Interfaces

		var interfacesStr strings.Builder
		// For each interface
		hostID := 1
		selfLoopbackIP := ""
		for interfaceName, interfaceInfo := range interfaces {
			if interfaceInfo.Role == "none" {
				continue
			}
			interfaceStr := "interface"
			interfaceStr += " "
			interfaceStr += interfaceName
			interfaceIP := ""

			switch interfaceInfo.Role {
			case "loopback":
				interfaceIP = generateLoopback(data.NetworkSubnet, router.Name)
        if eBGP_activated {
					interfaceIP = generateLoopbackIPv6(data.NetworkSubnet, router.Name)
					selfLoopbackIP = interfaceIP
				}
			case "internal":
				interfaceIP = generateIPv6(data.NetworkSubnet, interfaceInfo.Subnet, interfaceInfo.HostID)
			case "ebgp":
				if interfaceInfo.LocalIPv6 == "" {
					panic(fmt.Sprintf(
						"eBGP interface %s on router %s has no local_ipv6",
						interfaceName,
						router.Name,
					))
				}
				interfaceIP = interfaceInfo.LocalIPv6
			}

			interfacesStr.WriteString(ConfIPv6(interfaceIP, interfaceName))
			interfacesStr.WriteString("\n")
			if interfaceInfo.Role == "internal" || interfaceInfo.Role == "loopback" {
				interfacesStr.WriteString(" ")
				interfacesStr.WriteString(confInterfaceStr)
				interfacesStr.WriteString("\n")
				if data.Protocol == "OSPF" && interfaceInfo.Role == "internal" {
					interfacesStr.WriteString(fmt.Sprintf(" ipv6 ospf cost %d\n", interfaceInfo.Cost))
				}
				interfacesStr.WriteString("\n!\n")
			}
			if interfaceInfo.Role == "eBGP" {
				is_eBGP_router = true
			}
			hostID++
		}
		interfacesStr.WriteString("!")

		header := strHeader()
		subHeader := strSubHeader()
		tail := strTail()

		localAS := data.Name[2:]
		bgpConfStr := "router bgp " + localAS
		bgpConfStr += "\n"
		bgpConfStr += " bgp router-id " + router.RouterID
		bgpConfStr += " \n"

		bgpNeighborActivate := ""
		bgpSelfStaticRoute := ""
		if is_eBGP_router {
			bgpConfStr += "\n no bgp default ipv4-unicast"
			bgpConfStr += "\n"

			loopbackIP := generateLoopbackIPv6(data.NetworkSubnet, router.Name)

			idx := strings.LastIndex(loopbackIP, ":")
			loopbackSubnet := loopbackIP[:idx+1]
			
			bgpNeighborActivate += "  network " + loopbackSubnet + "/48"
			bgpNeighborActivate += "\n"

			bgpSelfStaticRoute = "ipv6 route " + loopbackSubnet + "/48" + " Null0"
		}
		

		for _, neighborIP := range loopbackNeighbours {
			// for each neighbor, neighbor activate + update-source Loopback0
			if neighborIP == selfLoopbackIP {
				continue
			}


			neighborAddr := strings.Split(neighborIP, "/")[0]

			bgpConfStr += " neighbor " + neighborAddr + " remote-as " + localAS
			bgpConfStr += " \n"
			bgpConfStr += " neighbor " + neighborAddr + " update-source " + "Loopback0"
			bgpConfStr += " \n"
			
			bgpNeighborActivate += "  neighbor " + neighborAddr + " activate"
			bgpNeighborActivate += "  \n"
		}
		bgpConfStr += " !"
		bgpConfStr += "\n address-family ipv4"
		bgpConfStr += "\n exit-address-family"
		bgpConfStr += "\n"
		bgpConfStr += " !"
		bgpConfStr += "\n"

		bgpConfStr += " address-family ipv6"
		bgpConfStr += "\n"
		bgpConfStr += bgpNeighborActivate
		bgpConfStr += "\n exit-address-family"

		

		content := fmt.Sprintf("%s %s\n%s \n%s\n%s\n%s\n%s\n%s\n%s",
			header,
			router.Name,
			subHeader,
			interfacesStr.String(),
			bgpConfStr,
			bgpSelfStaticRoute,
			ospfConfStr,
			ripConfStr,
			tail)

		_, err = file.WriteString(content)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(".cfg écrit sous", FILENAME)
	}
}

func strHeader() string {
	header := "!"
	header += "\nservice timestamps debug datetime msec"
	header += "\nservice timestamps log datetime msec"
	header += "\nno service password-encryption"
	header += "\n!"
	header += "\nhostname"

	return header
}

func strSubHeader() string {
	subHeader := "!\n!"
	subHeader += "\nip cef"
	subHeader += "\nno ip domain-lookup"
	subHeader += "\nipv6 unicast-routing"
	subHeader += "\nipv6 cef"
	subHeader += "\nip tcp synwait 5"
	subHeader += "\nmultilink bundle-name authenticated"

	return subHeader
}

func strTail() string {
	tail := "!\n!"
	tail += "\nno ip icmp rate-limit unreachable"
	tail += "\n!"
	tail += "\nline con 0"
	tail += "\nexec-timeout 0 0"
	tail += "\nlogging synchronous"
	tail += "\nprivilege level 15"
	tail += "\nno login"
	tail += "\nline aux 0"
	tail += "\nexec-timeout 0 0"
	tail += "\nlogging synchronous"
	tail += "\nprivilege level 15"
	tail += "\nno login"
	tail += "\n!"
	tail += "\n!"
	tail += "\nend"

	return tail
}
