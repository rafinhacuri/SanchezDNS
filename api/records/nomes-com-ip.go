package records

import (
	"context"
	"net"
	"slices"
	"strings"

	"resty.dev/v3"
)

func NomesComIP(ctx context.Context, httpc *resty.Client, zonas []ZoneInfo, alvo net.IP) ([]string, error) {
	if alvo == nil {
		return nil, nil
	}

	tipo := "AAAA"
	if alvo.To4() != nil {
		tipo = "A"
	}

	var nomes []string

	for _, info := range zonas {
		nome := strings.TrimSuffix(info.Name, ".")
		if strings.HasSuffix(nome, ".in-addr.arpa") || strings.HasSuffix(nome, ".ip6.arpa") {
			continue
		}

		zone, err := FetchZone(ctx, httpc, info.Name)
		if err != nil {
			return nil, err
		}

		nomes = append(nomes, nomesComIPNaZona(zone, tipo, alvo)...)
	}

	slices.Sort(nomes)

	return slices.Compact(nomes), nil
}

func nomesComIPNaZona(zone Zone, tipo string, alvo net.IP) []string {
	var nomes []string

	for _, rr := range zone.RRSets {
		if rr.Type != tipo {
			continue
		}

		for _, rec := range rr.Records {
			ip := net.ParseIP(strings.TrimSpace(strings.TrimSuffix(rec.Content, ".")))
			if ip != nil && ip.Equal(alvo) {
				nomes = append(nomes, rr.Name)

				break
			}
		}
	}

	return nomes
}
