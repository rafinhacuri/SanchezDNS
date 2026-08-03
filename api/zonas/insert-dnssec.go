package zonas

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func InsertDnssec(ctx context.Context, httpc *resty.Client, domain string) error {
	payload := cryptoKeyPayload{
		Active:  true,
		KeyType: "ksk",
	}

	resp, err := httpc.R().
		SetContext(ctx).
		SetBody(payload).
		Post(fmt.Sprintf("/api/v1/servers/%s/zones/%s/cryptokeys", env.C.DnsServerId, url.PathEscape(domain)))
	if err != nil {
		return err
	}

	if resp.IsStatusFailure() {
		return errors.New("pdns: " + resp.Status() + " " + resp.String())
	}

	return nil
}
