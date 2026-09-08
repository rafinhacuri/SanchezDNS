package zonas

import (
	"context"

	"resty.dev/v3"
)

func FetchDnssec(ctx context.Context, httpc *resty.Client, zone ZonePdns) (bool, error) {
	if zone.Dnssec != nil {
		return *zone.Dnssec, nil
	}

	detalhe, err := FetchDetalhe(ctx, httpc, zone.Name)
	if err != nil {
		return false, err
	}

	return detalhe.Dnssec, nil
}
