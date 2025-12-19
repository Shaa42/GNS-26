package main

import (
    "fmt"
    "log"
    "os"
)

func writeConfig(routerName string, data map[string]map[string]string) {
    FILENAME := routerName + "_configs_i" + routerName[1:] + "_startup-config" + ".cfg"

    file, err := os.Create(FILENAME)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

	links := data[routerName]

	interfacesStr := ""
	// For each interface
	for interfaceName, ip := range links {

	  interfacesStr += fmt.Sprintf("%s %s\n%s %s\n%s",
	  "interface",
	  interfaceName,
	  "no ip address\n negotiation auto\n ipv6 address",
	  ip,
	  "ipv6 enable\n !\n")
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

	content := fmt.Sprintf("%s %s\n%s %s", header, routerName, interfacesStr, tail)
	
    _, err = file.WriteString(content)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(".cfg écrit sous", FILENAME)
}