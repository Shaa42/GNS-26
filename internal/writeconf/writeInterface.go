package writeconf

func ConfIPv6(ip string, interf string) string {
	/*
	 * Configure an interface with an IPv6 address
	 */
	str := "interface " + interf + "\n"
	str += "    ipv6 address " + ip + "\n"
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

func ConfBGP(asNum string, routerID string) string {
	str := "router bgp " + asNum + "\n"
	str += "no bgp default ipv4-unicast\n"
	str += "bgp router-id " + routerID + "\n"
	return str
}
