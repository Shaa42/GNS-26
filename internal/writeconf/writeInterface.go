package writeconf

import (
	"fmt"
	"strings"
)

func generateIPv6(prefix string, subnetID string, ifaceID string) string {
	// Example: AS 1 subnet 2 hostid 1
	// 2001:1:2::1/64
	parts := strings.Split(prefix, "/")
	base := strings.TrimSuffix(parts[0], "::")
	return fmt.Sprintf("%s:%s::%s/64", base, subnetID, ifaceID)
}
func generateLoopbackIPv6(prefix string, routerName string) string {
	// AS 1 rID 2
	// Example: 2001:1:100::2/128
	parts := strings.Split(prefix, "/")
	base := strings.TrimSuffix(parts[0], "::")
	id := routerName[1:]
	return fmt.Sprintf("%s:100::%s/128", base, id)
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