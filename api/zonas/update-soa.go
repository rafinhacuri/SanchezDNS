//nolint:contextcheck
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

func UpdateSoa(
	ctx context.Context,
	zoneID,
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

	soaName := strings.TrimSuffix(startOfAuthority, ".")
	soaEmail := strings.TrimSuffix(email, ".")

	ttl := 3600

	body := records.PDNSZonePatchRequest{
		RRSets: []records.PDNSRRSetChange{
			{
				Name:       zoneID,
				Type:       "SOA",
				TTL:        &ttl,
				ChangeType: "REPLACE",
				Records: []records.Record{
					{
						Content:  fmt.Sprintf("%s. %s. 1 %d %d %d %d", soaName, soaEmail, refresh, retry, expire, negativeCacheTtl),
						Disabled: false,
					},
				},
			},
		},
	}

	resp, err := httpc.R().
		SetContext(ctx).
		SetBody(body).
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, zoneID))
	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", errors.New("falha ao atualizar SOA: " + resp.String())
	}

	go logs.InsertLog(zoneID, user, "update_soa", "Atualizado registro SOA para a zona "+zoneID)

	return "registro SOA atualizado com sucesso", nil
}
