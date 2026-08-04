package zonas

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func Delete(ctx context.Context, httpc *resty.Client, domain string) error {
	resp, err := httpc.R().
		SetContext(ctx).
		Delete(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, url.PathEscape(domain)))
	if err != nil {
		return err
	}

	if resp.IsStatusFailure() {
		return errors.New("pdns: " + resp.Status() + " " + resp.String())
	}

	return nil
}
