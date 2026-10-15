package records

import (
	"context"
	"fmt"
	"net/url"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func Patch(ctx context.Context, httpc *resty.Client, zona string, body PDNSZonePatchRequest) error {
	resp, err := httpc.R().
		SetContext(ctx).
		SetBody(body).
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, url.PathEscape(zona)))
	if err != nil {
		return err
	}

	if resp.StatusCode() != 204 && resp.StatusCode() != 201 {
		return &PdnsError{Status: resp.StatusCode(), Corpo: resp.String()}
	}

	return nil
}
