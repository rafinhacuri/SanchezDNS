package records

import (
	"context"
	"errors"
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
		return errors.New("pdns: " + resp.Status() + " " + resp.String())
	}

	return nil
}
