package records

import (
	"fmt"
	"net"
	"strings"
)

func NomeReverso(ip net.IP) string {
	ipv4 := ip.To4()
	if ipv4 != nil {
		octetos := strings.Split(ipv4.String(), ".")
		if len(octetos) != 4 {
			return ""
		}

		return fmt.Sprintf("%s.%s.%s.%s.in-addr.arpa.", octetos[3], octetos[2], octetos[1], octetos[0])
	}

	ipv6 := ip.To16()
	if ipv6 == nil {
		return ""
	}

	var hex strings.Builder

	for i := 0; i < 16; i += 2 {
		fmt.Fprintf(&hex, "%04x", uint16(ipv6[i])<<8|uint16(ipv6[i+1]))
	}

	texto := hex.String()
	nibbles := make([]string, 0, len(texto))

	for i := len(texto) - 1; i >= 0; i-- {
		nibbles = append(nibbles, string(texto[i]))
	}

	return strings.Join(nibbles, ".") + ".ip6.arpa."
}
