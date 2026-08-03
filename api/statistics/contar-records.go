package statistics

import (
	"context"

	"resty.dev/v3"
)

func ContarRecords(ctx context.Context, httpc *resty.Client, zones []PdnsZone) int {
	total := 0

	for _, zone := range zones {
		detalhes, err := FetchDetalhes(ctx, httpc, zone.ID)
		if err != nil {
			continue
		}

		for _, rrset := range detalhes.RRsets {
			total += len(rrset.Records)
		}
	}

	return total
}
