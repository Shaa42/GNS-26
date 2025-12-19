package writeconf

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func WriteConfig(routerName string, data map[string]map[string]string, internalProtocol string) {
	rId := routerName[1:]
	FILENAME := routerName + "_configs_i" + rId + "_startup-config" + ".cfg"
	ripStr := ""

	switch internalProtocol {
	case "RIP":
		ripStr = "ipv6 rip rip_process enable"

	default:
		panic("unrecognized internal routing protocol (atm only RIP can be used)")
	}

	file, err := os.Create(FILENAME)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	links := data[routerName]

	var interfacesStr strings.Builder
	// For each interface
	for interfaceName, ip := range links {

		fmt.Fprintf(&interfacesStr, "%s %s\n%s %s\n%s\n%s",
			"interface",
			interfaceName,
			"no ip address\n negotiation auto\n ipv6 address",
			ip,
			"ipv6 enable",
			ripStr)
	}

	header := "!\n!\n!"
	header += "\nservice timestamps debug datetime msec"
	header += "\nservice timestamps log datetime msec"
	header += "\nno service password-encryption"
	header += "\n!"
	header += "\nhostname"

	tail := "\n!"
	tail += "\nip cef"
	tail += "\nno ip domain-lookup"
	tail += "\nno ip icmp rate-limit unreachable"
	tail += "\nip tcp synwait 5"
	tail += "\nno cdp log mismatch duplex"
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

	content := fmt.Sprintf("%s %s\n%s \n%s",
		header,
		routerName,
		interfacesStr.String(),
		tail)

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(".cfg écrit sous", FILENAME)
}
