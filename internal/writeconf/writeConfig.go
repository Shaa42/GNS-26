package writeconf

import (
	"fmt"
	"log"
	"os"
)

func getRouterID(rN string) string {

	rID := fmt.Sprintf("%d.%d.%d.%s",
		1,
		1,
		1,
		rN)

	return rID
}
func WriteConfig(routerName string, routerID string, data map[string]map[string]string, internalProtocol string) {
	rN := routerName[1:]
	FILENAME := routerName + "_configs_i" + rN + "_startup-config" + ".cfg"
	confInterfaceStr := ""
	ospfConfStr := "!\n!"

	// Depending on the protocol, will add strings related specifically to that protocol
	switch internalProtocol {
	case "RIP":
		confInterfaceStr = "ipv6 rip rip_process enable"
	case "OSPF":
		confInterfaceStr = "ipv6 ospf 1 area 0" // WARNING: currently CANNOT handle multiple areas
		ospfConfStr += "\nipv6 router ospf 1"
		ospfConfStr += "\n router-id "
		ospfConfStr += getRouterID(rN)

	default:
		panic("unrecognized internal routing protocol (atm only RIP and OSPF can be used)")
	}

	file, err := os.Create(FILENAME)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	links := data[routerName]

	interfacesStr := ""
	// For each interface
	for interfaceName, ip := range links {
		if ip == "-1" {
			continue
		}
		interfaceStr := "interface"
		interfaceStr += " "
		interfaceStr += interfaceName
		interfacesStr += ConfIPv6(ip, interfaceName)
		interfacesStr += " "
		interfacesStr += confInterfaceStr
		interfacesStr += "\n!\n"
	}
	interfacesStr += "!"

	header := "!"
	header += "\nservice timestamps debug datetime msec"
	header += "\nservice timestamps log datetime msec"
	header += "\nno service password-encryption"
	header += "\n!"
	header += "\nhostname"

	subHeader := "!\n!"
	subHeader += "\nip cef"
	subHeader += "\nno ip domain-lookup"
	subHeader += "\nipv6 unicast-routing"
	subHeader += "\nipv6 cef"
	subHeader += "\nip tcp synwait 5"
	subHeader += "\nmultilink bundle-name authenticated"

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

	content := fmt.Sprintf("%s %s\n%s \n%s\n%s\n%s",
		header,
		routerName,
		subHeader,
		interfacesStr,
		ospfConfStr,
		tail)

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(".cfg écrit sous", FILENAME)
}
