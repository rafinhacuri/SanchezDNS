package zonas

import (
	"context"
	"errors"
	"fmt"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func Insert(ctx context.Context, httpc *resty.Client, domain string) error {
	payload := createZonePayload{
		Name:       domain,
		Kind:       "Native",
		SOAEditAPI: "DEFAULT",
	}

	resp, err := httpc.R().
		SetContext(ctx).
		SetBody(payload).
		Post(fmt.Sprintf("/api/v1/servers/%s/zones", env.C.DnsServerId))
	if err != nil {
		return err
	}

	if resp.IsStatusFailure() {
		return errors.New("pdns: " + resp.Status() + " " + resp.String())
	}

	return nil
}
