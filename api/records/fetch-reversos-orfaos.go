package records

import (
	"context"

	"resty.dev/v3"
)

func FetchReversosOrfaos(ctx context.Context, httpc *resty.Client, zonaReversa string) ([]ReversoOrfao, error) {
	zone, err := FetchZone(ctx, httpc, zonaReversa)
	if err != nil {
		return nil, err
	}

	zonas, err := FetchZonas(ctx, httpc)
	if err != nil {
		return nil, err
	}

	ips, err := IPsForward(ctx, httpc, zonas)
	if err != nil {
		return nil, err
	}

	orfaos := []ReversoOrfao{}

	for _, rr := range zone.RRSets {
		if rr.Type != "PTR" {
			continue
		}

		ip := IPDoReverso(rr.Name)
		if ip == nil || ips[ip.String()] {
			continue
		}

		alvo := ""
		if len(rr.Records) > 0 {
			alvo = rr.Records[0].Content
		}

		orfaos = append(orfaos, ReversoOrfao{NomeReverso: rr.Name, IP: ip.String(), Alvo: alvo})
	}

	return orfaos, nil
}
