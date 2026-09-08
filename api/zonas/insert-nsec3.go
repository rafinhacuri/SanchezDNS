package zonas

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

const Nsec3Param = "1 0 0 -"

func InsertNsec3(ctx context.Context, httpc *resty.Client, domain string) error {
	payload := nsec3ParamPayload{Nsec3Param: Nsec3Param}

	resp, err := httpc.R().
		SetContext(ctx).
		SetBody(payload).
		Put(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, url.PathEscape(domain)))
	if err != nil {
		return err
	}

	if resp.IsStatusFailure() {
		return errors.New("pdns: " + resp.Status() + " " + resp.String())
	}

	return nil
}
