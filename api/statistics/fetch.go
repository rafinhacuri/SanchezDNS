package statistics

import (
	"context"
	"errors"
	"fmt"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func Fetch(ctx context.Context, httpc *resty.Client) ([]Stat, error) {
	var stats []Stat

	resp, err := httpc.R().
		SetContext(ctx).
		SetResult(&stats).
		Get(fmt.Sprintf("/api/v1/servers/%s/statistics", env.C.DnsServerId))
	if err != nil {
		return nil, err
	}

	if resp.IsStatusFailure() {
		return nil, errors.New("pdns: " + resp.Status() + " " + resp.String())
	}

	return stats, nil
}
