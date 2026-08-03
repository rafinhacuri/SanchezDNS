package zonas

import (
	"context"
	"errors"
	"fmt"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func Fetch(ctx context.Context, httpc *resty.Client) ([]ZonePdns, error) {
	var zones []ZonePdns

	resp, err := httpc.R().
		SetContext(ctx).
		SetResult(&zones).
		Get(fmt.Sprintf("/api/v1/servers/%s/zones", env.C.DnsServerId))
	if err != nil {
		return nil, err
	}

	if resp.IsStatusFailure() {
		return nil, errors.New("pdns: " + resp.Status() + " " + resp.String())
	}

	return zones, nil
}
