package zonas

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func FetchDetalhe(ctx context.Context, httpc *resty.Client, domain string) (ZoneDetalhe, error) {
	var detail zoneDetail

	resp, err := httpc.R().
		SetContext(ctx).
		SetResult(&detail).
		Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, url.PathEscape(domain)))
	if err != nil {
		return ZoneDetalhe{}, err
	}

	if resp.IsStatusFailure() {
		return ZoneDetalhe{}, errors.New("pdns: " + resp.Status() + " " + resp.String())
	}

	return ZoneDetalhe{
		Dnssec: detail.Dnssec,
		Nsec3:  strings.TrimSpace(detail.Nsec3Param) != "",
	}, nil
}
