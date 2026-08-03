package statistics

import (
	"context"
	"errors"
	"fmt"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func FetchDetalhes(ctx context.Context, httpc *resty.Client, zoneID string) (ZoneDetails, error) {
	var detalhes ZoneDetails

	resp, err := httpc.R().
		SetContext(ctx).
		SetResult(&detalhes).
		Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, zoneID))
	if err != nil {
		return ZoneDetails{}, err
	}

	if resp.IsStatusFailure() {
		return ZoneDetails{}, errors.New("pdns: " + resp.Status() + " " + resp.String())
	}

	return detalhes, nil
}
