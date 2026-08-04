package zonas

import "strings"

func MatchTipo(nome, tipo string) bool {
	v4 := strings.HasSuffix(nome, ".in-addr.arpa.")
	v6 := strings.HasSuffix(nome, ".ip6.arpa.")

	switch tipo {
	case "reverse":
		return v4
	case "reverse-ipv6":
		return v6
	default:
		return !v4 && !v6
	}
}
