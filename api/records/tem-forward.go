package records

import (
	"context"
	"net"
	"strings"

	"resty.dev/v3"
)

func TemForward(ctx context.Context, httpc *resty.Client, zonas []ZoneInfo, alvo net.IP) bool {
	if alvo == nil {
		return false
	}

	tipo := "AAAA"
	if alvo.To4() != nil {
		tipo = "A"
	}

	for _, info := range zonas {
		nome := strings.TrimSuffix(info.Name, ".")
		if strings.HasSuffix(nome, ".in-addr.arpa") || strings.HasSuffix(nome, ".ip6.arpa") {
			continue
		}

		zone, err := FetchZone(ctx, httpc, info.Name)
		if err != nil {
			continue
		}

		if zoneTemIP(zone, tipo, alvo) {
			return true
		}
	}

	return false
}

func zoneTemIP(zone Zone, tipo string, alvo net.IP) bool {
	for _, rr := range zone.RRSets {
		if rr.Type != tipo {
			continue
		}

		for _, rec := range rr.Records {
			ip := net.ParseIP(strings.TrimSpace(strings.TrimSuffix(rec.Content, ".")))
			if ip != nil && ip.Equal(alvo) {
				return true
			}
		}
	}

	return false
}
