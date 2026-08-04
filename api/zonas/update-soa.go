package zonas

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
	"github.com/rafinhacuri/SanchezDNS/api/records"
)

func UpdateSoa(ctx context.Context, httpc *resty.Client, zoneID string, soa Soa) error {
	ttl := 3600

	content := fmt.Sprintf(
		"%s. %s. 1 %d %d %d %d",
		strings.TrimSuffix(soa.StartOfAuthority, "."),
		strings.TrimSuffix(soa.Email, "."),
		soa.Refresh,
		soa.Retry,
		soa.Expire,
		soa.NegativeCacheTtl,
	)

	body := records.PDNSZonePatchRequest{
		RRSets: []records.PDNSRRSetChange{
			{
				Name:       zoneID,
				Type:       "SOA",
				TTL:        &ttl,
				ChangeType: "REPLACE",
				Records:    []records.Record{{Content: content, Disabled: false}},
			},
		},
	}

	resp, err := httpc.R().
		SetContext(ctx).
		SetBody(body).
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, url.PathEscape(zoneID)))
	if err != nil {
		return err
	}

	if resp.IsStatusFailure() {
		return errors.New("pdns: " + resp.Status() + " " + resp.String())
	}

	return nil
}
