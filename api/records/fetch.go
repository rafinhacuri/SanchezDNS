package records

import (
	"context"

	"resty.dev/v3"
)

func Fetch(ctx context.Context, httpc *resty.Client, zona string) ([]Simplified, *Soa, error) {
	zone, err := FetchZone(ctx, httpc, zona)
	if err != nil {
		return nil, nil, err
	}

	lista, soa := Simplificar(zone)

	return lista, soa, nil
}
