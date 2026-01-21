package writeconf

import (
	"fmt"
	"gns-26/internal/parseintent"
	"log"
	"os"
	"strconv"
	"strings"
)

func WriteConfig(data parseintent.InfoAS) {
	for _, router := range data.Routers {
		rN := router.Name[1:]
		FILENAME := router.Name + "_configs_i" + rN + "_startup-config" + ".cfg"
		confInterfaceStr := ""
		ripConfStr := ""
		ospfConfStr := "!\n!"
		bgpConf := buildEBGPConfig(data.Name, router, data.Links)

		// Depending on the protocol, will add strings related specifically to that protocol
		switch data.Protocol {
		case "RIPng":
			confInterfaceStr = "ipv6 rip rip_process enable"
			ripConfStr += "\nipv6 router rip rip_process"
			ripConfStr += "\n redistribute bgp " + data.Name[2:]
		case "OSPF":
			confInterfaceStr = "ipv6 ospf 1 area 0" // WARNING: currently CANNOT handle multiple areas
			ospfConfStr += "\nipv6 router ospf 1"
			ospfConfStr += "\n router-id "
			ospfConfStr += router.RouterID
			ospfConfStr += "\n redistribute bgp " + data.Name[2:]

		default:
			panic("unrecognized internal routing protocol (atm only RIP and OSPF can be used)")
		}

		file, err := os.Create(FILENAME)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		//links := data[routerName]
		interfaces := router.Interfaces

		var interfacesStr strings.Builder
		// For each interface
		hostID := 1
		for interfaceName, interfaceInfo := range interfaces {
			if interfaceInfo.Role == "none" {
				continue
			}
			interfaceStr := "interface"
			interfaceStr += " "
			interfaceStr += interfaceName
			interfaceIP := ""

			if interfaceInfo.Role == "loopback" {
				interfaceIP = generateLoopback(data.NetworkSubnet, rN)
			} else {
				interfaceIP = generateIPv6(data.NetworkSubnet, interfaceInfo.Subnet, interfaceInfo.HostID)
			}

			interfacesStr.WriteString(ConfIPv6(interfaceIP, interfaceName))
			interfacesStr.WriteString("\n")
			if interfaceInfo.Role != "loopback" {
				interfacesStr.WriteString(" no shutdown\n")
			}
			if interfaceInfo.Role == "internal" || interfaceInfo.Role == "loopback" {
				interfacesStr.WriteString(" ")
				interfacesStr.WriteString(confInterfaceStr)
				interfacesStr.WriteString("\n!\n")
			}
			hostID++
		}
		interfacesStr.WriteString("!")

		header := strHeader()
		subHeader := strSubHeader()
		tail := strTail()

		content := fmt.Sprintf("%s %s\n%s \n%s\n%s\n%s\n%s",
			header,
			router.Name,
			subHeader,
			interfacesStr.String(),
			ospfConfStr,
			ripConfStr,
			bgpConf,
			tail)

		_, err = file.WriteString(content)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(".cfg écrit sous", FILENAME)
	}
}

func buildEBGPConfig(
	asName string,
	router parseintent.InfoRouter,
	links []parseintent.InfoLink,
) string {
	var sb strings.Builder

	localAS, _ := strconv.Atoi(asName[2:])

	for _, link := range links {
		if link.Role != "eBGP" {
			continue
		}

		peerAS := link.Info["as"]
		peerIPv6 := link.Info["ipv6"]

		remoteAS, _ := strconv.Atoi(peerAS[2:])

		if _, ok := link.Endpoints[router.Name]; !ok {
			continue
		}

		sb.WriteString("router bgp ")
		sb.WriteString(strconv.Itoa(localAS))
		sb.WriteString("\n")

		sb.WriteString(" bgp router-id ")
		sb.WriteString(router.RouterID)
		sb.WriteString("\n")

		sb.WriteString(" neighbor ")
		sb.WriteString(peerIPv6)
		sb.WriteString(" remote-as ")
		sb.WriteString(strconv.Itoa(remoteAS))
		sb.WriteString("\n!\n")
	}

	return sb.String()
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
