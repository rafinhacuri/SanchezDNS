package records

import (
	"context"
	"strings"

	"resty.dev/v3"
)

func IPsForward(ctx context.Context, httpc *resty.Client, zonas []ZoneInfo) (map[string]bool, error) {
	ips := map[string]bool{}

	for _, info := range zonas {
		nome := strings.TrimSuffix(info.Name, ".")
		if strings.HasSuffix(nome, ".in-addr.arpa") || strings.HasSuffix(nome, ".ip6.arpa") {
			continue
		}

		zone, err := FetchZone(ctx, httpc, info.Name)
		if err != nil {
			return nil, err
		}

		for _, registro := range RegistrosIP(zone) {
			ips[registro.IP.String()] = true
		}
	}

	return ips, nil
}
