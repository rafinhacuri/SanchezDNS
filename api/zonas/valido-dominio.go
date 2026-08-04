package zonas

import "strings"

func ValidoDominio(domain, tipo string) bool {
	switch tipo {
	case "reverse":
		return strings.HasSuffix(domain, ".in-addr.arpa.")
	case "reverse-ipv6":
		return strings.HasSuffix(domain, ".ip6.arpa.")
	default:
		return true
	}
}
