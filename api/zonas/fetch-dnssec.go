package zonas

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func FetchDnssec(ctx context.Context, httpc *resty.Client, zone ZonePdns) (bool, error) {
	if zone.Dnssec != nil {
		return *zone.Dnssec, nil
	}

	var detail zoneDetail

	resp, err := httpc.R().
		SetContext(ctx).
		SetResult(&detail).
		Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, url.PathEscape(zone.Name)))
	if err != nil {
		return false, err
	}

	if resp.IsStatusFailure() {
		return false, errors.New("pdns: " + resp.Status() + " " + resp.String())
	}

	return detail.Dnssec, nil
}
