package records

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func FetchZonas(ctx context.Context, httpc *resty.Client) ([]ZoneInfo, error) {
	resp, err := httpc.R().
		SetContext(ctx).
		Get(fmt.Sprintf("/api/v1/servers/%s/zones", env.C.DnsServerId))
	if err != nil {
		return nil, err
	}

	if resp.IsStatusFailure() {
		return nil, errors.New("pdns: " + resp.Status() + " " + resp.String())
	}

	var zonas []ZoneInfo

	err = json.Unmarshal(resp.Bytes(), &zonas)
	if err != nil {
		return nil, err
	}

	return zonas, nil
}
