//nolint:contextcheck
package zonas

import (
	"context"
	"fmt"
	"time"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
	"github.com/rafinhacuri/SanchezDNS/api/logs"
)

func DeleteZone(ctx context.Context, domain, user string) (string, error) {
	httpc := resty.New().
		SetTimeout(30*time.Second).
		SetBaseURL(env.C.DnsHost).
		SetHeader("X-API-Key", env.C.DnsApiKey).
		SetHeader("Accept", "application/json").
		SetRetryCount(2)

	resp, err := httpc.R().SetContext(ctx).Delete(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, domain))
	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", err
	}

	go logs.InsertLog(domain, user, "delete_zone", "Excluída a zona "+domain)

	return "zona excluída com sucesso", nil
}
