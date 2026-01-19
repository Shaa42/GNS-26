package writeconf

import (
	"fmt"
	"strings"
)

func generateIPv6(prefix string, subnetID string, ifaceID string) string {
	parts := strings.Split(prefix, "/")
	base := strings.TrimSuffix(parts[0], "::")
	return fmt.Sprintf("%s:%s::%s/64", base, subnetID, ifaceID)
}

func ConfIPv6(ip string, interf string) string {
	/*
	 * Configure an interface with an IPv6 address
	 */

	str := "interface " + interf + "\n"
	str += " no ip address" + "\n"
	str += " negotiation auto" + "\n"
	str += " ipv6 address " + ip + "\n"
	str += " ipv6 enable" + "\n"
	return str
}

func ConfIPv6UR() string {
	return "ipv6 unicast-routing\n"
}

func ConfRIP() string {
	return "ipv6 rip RIP6 enable\n"
}

func ConfNoSD() string {
	return "no shutdown\n"
}

func ConfBGP(localAS int, routerID string, peerAS int, peerIPv6 string) string {
	cfg := ""
	cfg += fmt.Sprintf("router bgp %d\n", localAS)
	cfg += fmt.Sprintf(" bgp router-id %s\n", routerID)
	cfg += " bgp log-neighbor-changes\n"
	cfg += fmt.Sprintf(" neighbor %s remote-as %d\n", peerIPv6, peerAS)
	cfg += "address-family ipv6 unicast\n"
	cfg += fmt.Sprintf(" neighbor %s activate\n", peerIPv6)
	cfg += "exit-address-family"
	//cfg += fmt.Sprintf(" neighbor %s update-source Loopback0\n", peerIPv6)
	cfg += "\n"

	return cfg
}
