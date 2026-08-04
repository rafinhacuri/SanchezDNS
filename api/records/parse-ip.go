package records

import (
	"net"
	"strings"
)

func ParseIP(tipo, vl string) net.IP {
	ip := net.ParseIP(strings.TrimSpace(strings.TrimSuffix(vl, ".")))
	if ip == nil {
		return nil
	}

	if tipo == "A" {
		return ip.To4()
	}

	if ip.To4() != nil {
		return nil
	}

	return ip.To16()
}
