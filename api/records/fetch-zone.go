package records

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func FetchZone(ctx context.Context, httpc *resty.Client, zona string) (Zone, error) {
	resp, err := httpc.R().
		SetContext(ctx).
		Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, url.PathEscape(zona)))
	if err != nil {
		return Zone{}, err
	}

	if resp.IsStatusFailure() {
		return Zone{}, errors.New("pdns: " + resp.Status() + " " + resp.String())
	}

	var zone Zone

	err = json.Unmarshal(resp.Bytes(), &zone)
	if err != nil {
		return Zone{}, err
	}

	return zone, nil
}
