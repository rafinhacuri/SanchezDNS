package zonas

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func FetchCryptokeys(ctx context.Context, httpc *resty.Client, domain string) ([]CryptoKey, error) {
	var keys []CryptoKey

	resp, err := httpc.R().
		SetContext(ctx).
		SetResult(&keys).
		Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s/cryptokeys", env.C.DnsServerId, url.PathEscape(domain)))
	if err != nil {
		return nil, err
	}

	if resp.IsStatusFailure() {
		return nil, errors.New("pdns: " + resp.Status() + " " + resp.String())
	}

	return keys, nil
}
