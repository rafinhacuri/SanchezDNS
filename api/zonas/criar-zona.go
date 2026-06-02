//nolint:tagliatelle, contextcheck
package zonas

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
	"github.com/rafinhacuri/SanchezDNS/api/logs"
	"github.com/rafinhacuri/SanchezDNS/api/records"
)

type pdnsCreateZoneRequest struct {
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	SOAEditAPI string `json:"soa_edit_api"`
}
type pdnsCryptoKeyRequest struct {
	Active  bool   `json:"active"`
	KeyType string `json:"keytype"`
}

func CreateZone(
	ctx context.Context,
	domain,
	zoneType,
	startOfAuthority,
	email string,
	refresh,
	retry,
	expire,
	negativeCacheTtl int,
	user string,
) (string, error) {
	httpc := resty.New().
		SetTimeout(30*time.Second).
		SetBaseURL(env.C.DnsHost).
		SetHeader("X-API-Key", env.C.DnsApiKey).
		SetHeader("Accept", "application/json").
		SetRetryCount(2)

	domain = strings.TrimSuffix(domain, ".")
	domainWithDot := domain + "."

	soaMname := strings.TrimSuffix(startOfAuthority, ".")
	soaRname := strings.TrimSuffix(email, ".")

	zonePayload := pdnsCreateZoneRequest{
		Name:       domainWithDot,
		Kind:       "Native",
		SOAEditAPI: "DEFAULT",
	}

	resp, err := httpc.R().
		SetContext(ctx).
		SetBody(zonePayload).
		Post(fmt.Sprintf("/api/v1/servers/%s/zones", env.C.DnsServerId))
	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", errors.New("erro do PowerDNS ao criar zona: status=" + resp.Status() + ", body=" + resp.String())
	}

	dnssecPayload := pdnsCryptoKeyRequest{
		Active:  true,
		KeyType: "ksk",
	}

	respDNSSEC, err := httpc.R().
		SetContext(ctx).
		SetBody(dnssecPayload).
		Post(fmt.Sprintf("/api/v1/servers/%s/zones/%s/cryptokeys", env.C.DnsServerId, domainWithDot))
	if err != nil {
		return "", err
	} else if respDNSSEC.IsError() {
		return "", errors.New(
			"erro do PowerDNS ao ativar DNSSEC: status=" +
				respDNSSEC.Status() +
				", body=" +
				respDNSSEC.String(),
		)
	}

	ttl := 3600
	patchBody := records.PDNSZonePatchRequest{
		RRSets: []records.PDNSRRSetChange{
			{
				Name:       domainWithDot,
				Type:       "SOA",
				TTL:        &ttl,
				ChangeType: "REPLACE",
				Records: []records.Record{
					{
						Content:  fmt.Sprintf("%s. %s. 1 %d %d %d %d", soaMname, soaRname, refresh, retry, expire, negativeCacheTtl),
						Disabled: false,
					},
				},
			},
		},
	}

	respPatch, err := httpc.R().
		SetContext(ctx).
		SetBody(patchBody).
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, domainWithDot))
	if err != nil {
		return "", err
	} else if respPatch.IsError() {
		return "", errors.New("erro do PowerDNS no PATCH: status=" + respPatch.Status() + ", body=" + respPatch.String())
	}

	domain = strings.TrimSuffix(domain, ".")
	go logs.InsertLog(domain, user, "create_zone", fmt.Sprintf("Criada zona %s do tipo %s", domain, zoneType))

	return "zona criada com sucesso", nil
}
