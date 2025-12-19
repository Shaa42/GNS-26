package main

import (
    "fmt"
    "log"
    "os"
)

func writeConfig(routerName string, data map[string]map[string]string, internalProtocol string) {
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

	interfacesStr := ""
	// For each interface
	for interfaceName, ip := range links {

	  interfacesStr += fmt.Sprintf("%s %s\n%s %s\n%s\n%s",
	  "interface",
	  interfaceName,
	  "no ip address\n negotiation auto\n ipv6 address",
	  ip,
	  "ipv6 enable",
	  ripStr)
   	}

	header := `!
!
!
service timestamps debug datetime msec
service timestamps log datetime msec
no service password-encryption
!
hostname`

	tail := `!
ip cef
no ip domain-lookup
no ip icmp rate-limit unreachable
ip tcp synwait 5
no cdp log mismatch duplex
!
line con 0
exec-timeout 0 0
logging synchronous
privilege level 15
no login
line aux 0
exec-timeout 0 0
logging synchronous
privilege level 15
no login
!
!
end`

	content := fmt.Sprintf("%s %s\n%s \n%s",
	header,
	routerName,
	interfacesStr,
	tail)
	
    _, err = file.WriteString(content)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(".cfg écrit sous", FILENAME)
}