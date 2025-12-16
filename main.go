package main

import (
    "fmt"
    "log"
    "os"
)

func writeConfig(router_name string) {
    filename := "config" + ".cfg"

    file, err := os.Create(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

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

	content := fmt.Sprintf("%s %s\n%s", header, router_name, tail)
	
    _, err = file.WriteString(content)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(".cfg écrit sous", filename)
}

func main() {
	// Create the default .cfg file for router named R1
    writeConfig("R1")
}
