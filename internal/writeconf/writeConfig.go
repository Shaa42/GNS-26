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

	var ibgpNeighbors []string
	var ebgpNeighbors [][2]string
	asProviders := data.Providers

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

						ibgpNeighbors = append(ibgpNeighbors, interfaceIP)
					}

				}
			}
		}
	}

	// Setup the .cfg file for every router
	for _, router := range data.Routers {
		is_eBGP_router = false
		ebgpNeighbors = nil
		policyStr := ""

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
				// Turns ON the loopbacks only if eBGP (thus iBGP) is activated
				if eBGP_activated {
					interfaceIP = generateLoopbackIPv6(data.NetworkSubnet, router.Name)
					selfLoopbackIP = interfaceIP
				}
			case "internal":
				interfaceIP = generateIPv6(data.NetworkSubnet, interfaceInfo.Subnet, interfaceInfo.HostID)
			case "ebgp":
				// If it's an eBGP connection, the specified IP is static
				is_eBGP_router = true
				if interfaceInfo.LocalIPv6 == "" {
					panic(fmt.Sprintf(
						"eBGP interface %s on router %s has no local_ipv6",
						interfaceName,
						router.Name,
					))
				}
				interfaceIP = interfaceInfo.LocalIPv6

				var eBGP_peer [2]string // [REMOTE_AS_NUM, REMOTE_AS_IP]
				eBGP_peer[0] = interfaceInfo.PeerAS
				eBGP_peer[1] = interfaceInfo.PeerIPv6

				// Adds the remote-AS router as neighbor
				ebgpNeighbors = append(ebgpNeighbors, eBGP_peer)

			}

			interfacesStr.WriteString(ConfIPv6(interfaceIP, interfaceName))
			interfacesStr.WriteString("\n")
			interfacesStr.WriteString(" ")
			interfacesStr.WriteString(confInterfaceStr)
			interfacesStr.WriteString("\n")

			// Sets up the OSPF link cost
			if data.Protocol == "OSPF" && interfaceInfo.Role == "internal" {
				interfacesStr.WriteString(fmt.Sprintf(" ipv6 ospf cost %d\n", interfaceInfo.Cost))
			}
			interfacesStr.WriteString("\n!\n")
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

			//set a name for each policy
			prefixListName := "PL-EXPORT-" + router.Name
			routerMapName := "RM-EXPORT-" + router.Name

			preferredProvider := ""
			isMultiHomed := len(asProviders) > 1
			if isMultiHomed {
				preferredProvider = asProviders[0]
				for _, p := range asProviders[1:] {
					if asNumber(p) < asNumber(preferredProvider) {
						preferredProvider = p
					}
				}
			}

			//add ipv6 prefix-list and route-map to policyStr
			policyStr += "ipv6 prefix-list " + prefixListName + " seq 10 permit " + loopbackSubnet + "/48\n"
			policyStr += "route-map " + routerMapName + " permit 10\n"
			policyStr += " match ipv6 address prefix-list " + prefixListName + "\n!\n"
			if preferredProvider != "" {
				policyStr += "route-map RM-IN-LOCAL-PREF " + router.Name + " permit 10\n"
				policyStr += " set local-preference 200\n!\n"
			}

			// Advertise the loopback subnet
			bgpNeighborActivate += "  network " + loopbackSubnet + "/48"
			bgpNeighborActivate += "\n"

			// Sets up a static route in order to allow the loopback subnet advertisement
			bgpSelfStaticRoute = "ipv6 route " + loopbackSubnet + "/48" + " Null0\n"

			for _, remotePeer := range ebgpNeighbors {

				remotePeerASName := remotePeer[0]
				remotePeerAS := remotePeerASName[2:] // remotePeer[0] is like AS111
				remotePeerIP := remotePeer[1]

				isProvider := contains(asProviders, remotePeerASName)
				isCustomer := !isProvider
				if isCustomer {
					bgpConfStr += "  neighbor " + remotePeerIP + " default-originate\n"
				}

				if remotePeerASName == preferredProvider {
					bgpConfStr += " neighbor " + remotePeerIP + " route-map RM-IN-LOCAL-PREF " + router.Name + " in\n"
				}

				bgpConfStr += " neighbor " + remotePeerIP + " remote-as " + remotePeerAS
				bgpConfStr += " \n"

				//apply route-map OUT to eBGP neighbor
				bgpConfStr += " neighbor " + remotePeerIP + " route-map " + routerMapName + " out"
				bgpConfStr += " \n"

				bgpNeighborActivate += "  neighbor " + remotePeerIP + " activate"
				bgpNeighborActivate += "  \n"
			}
		}

		for _, neighborIP := range ibgpNeighbors {
			// for each neighbor, neighbor activate + update-source Loopback0
			if neighborIP == selfLoopbackIP {
				continue
			}

			neighborAddr := strings.Split(neighborIP, "/")[0]

			// Sets up the neighbours and use their Loopback as source
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

		content := fmt.Sprintf("%s %s\n%s \n%s\n%s\n%s\n%s\n%s\n%s\n%s",
			header,
			router.Name,
			subHeader,
			interfacesStr.String(),
			policyStr,
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
