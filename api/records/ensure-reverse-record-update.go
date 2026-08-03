package records

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func ensureReverseRecordUpdate(parentCtx context.Context, tipo, vl, zona, name, oldVl string) error {
	if tipo != "A" && tipo != "AAAA" {
		return nil
	}

	ctx, cancel := context.WithTimeout(parentCtx, 6*time.Second)
	defer cancel()

	switch tipo {
	case "A":
		ensureReverseIPv4Update(ctx, vl, zona, name, oldVl)
	case "AAAA":
		ensureReverseIPv6Update(ctx, vl, zona, name, oldVl)
	}

	return nil
}

func findBestReverseZone(zones []pdnsZoneInfo, fullReverseName, suffix string) string {
	bestZone := ""
	fullRevNoDot := strings.TrimSuffix(fullReverseName, ".")

	for _, z := range zones {
		zNameNoDot := strings.TrimSuffix(z.Name, ".")

		if !strings.HasSuffix(zNameNoDot, suffix) {
			continue
		}

		if strings.HasSuffix(fullRevNoDot, zNameNoDot) {
			if len(zNameNoDot) > len(strings.TrimSuffix(bestZone, ".")) {
				bestZone = z.Name
			}
		}
	}

	return bestZone
}

func hasForwardRecordsForIPv4(ctx context.Context, httpc *resty.Client, zones []pdnsZoneInfo, targetIP net.IP) bool {
	if targetIP == nil {
		return false
	}

	for _, z := range zones {
		zNameNoDot := strings.TrimSuffix(z.Name, ".")

		if strings.HasSuffix(zNameNoDot, ".in-addr.arpa") || strings.HasSuffix(zNameNoDot, ".ip6.arpa") {
			continue
		}

		zoneResp, err := httpc.R().SetContext(ctx).Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, z.Name))
		if err != nil || zoneResp.IsStatusFailure() {
			continue
		}

		var zoneData Zone

		err = json.Unmarshal(zoneResp.Bytes(), &zoneData)
		if err != nil {
			continue
		}

		for _, rr := range zoneData.RRSets {
			if rr.Type != "A" {
				continue
			}

			for _, rec := range rr.Records {
				content := strings.TrimSpace(strings.TrimSuffix(rec.Content, "."))

				ip := net.ParseIP(content)
				if ip == nil {
					continue
				}

				if ip.To4() != nil && ip.Equal(targetIP) {
					return true
				}
			}
		}
	}

	return false
}

func hasForwardRecordsForIPv6(ctx context.Context, httpc *resty.Client, zones []pdnsZoneInfo, targetIP net.IP) bool {
	if targetIP == nil {
		return false
	}

	for _, z := range zones {
		zNameNoDot := strings.TrimSuffix(z.Name, ".")

		if strings.HasSuffix(zNameNoDot, ".in-addr.arpa") || strings.HasSuffix(zNameNoDot, ".ip6.arpa") {
			continue
		}

		zoneResp, err := httpc.R().SetContext(ctx).Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, z.Name))
		if err != nil || zoneResp.IsStatusFailure() {
			continue
		}

		var zoneData Zone

		err = json.Unmarshal(zoneResp.Bytes(), &zoneData)
		if err != nil {
			continue
		}

		for _, rr := range zoneData.RRSets {
			if rr.Type != "AAAA" {
				continue
			}

			for _, rec := range rr.Records {
				content := strings.TrimSpace(strings.TrimSuffix(rec.Content, "."))

				ip := net.ParseIP(content)
				if ip == nil {
					continue
				}

				if ip.To4() == nil && ip.To16() != nil && ip.Equal(targetIP) {
					return true
				}
			}
		}
	}

	return false
}

func ensureReverseIPv4Update(ctx context.Context, vl, zona, name, oldVl string) {
	httpc := resty.New().
		SetTimeout(30*time.Second).
		SetBaseURL(env.C.DnsHost).
		SetHeader("X-API-Key", env.C.DnsApiKey).
		SetHeader("Accept", "application/json").
		SetRetryCount(2)

	ipNewStr := strings.TrimSpace(strings.TrimSuffix(vl, "."))

	ipNew := net.ParseIP(ipNewStr)
	if ipNew != nil {
		ipNew = ipNew.To4()
	}

	var ipOld net.IP

	if strings.TrimSpace(oldVl) != "" {
		ipOldStr := strings.TrimSpace(strings.TrimSuffix(oldVl, "."))

		ipOld = net.ParseIP(ipOldStr)
		if ipOld != nil {
			ipOld = ipOld.To4()
		}
	}

	if ipNew == nil && ipOld == nil {
		return
	}

	zoneListResp, err := httpc.R().SetContext(ctx).Get(fmt.Sprintf("/api/v1/servers/%s/zones", env.C.DnsServerId))
	if err != nil || zoneListResp.IsStatusFailure() {
		return
	}

	var zones []pdnsZoneInfo

	err = json.Unmarshal(zoneListResp.Bytes(), &zones)
	if err != nil {
		return
	}

	now := time.Now().Unix()
	fqdn := buildFQDN(zona, name)

	if ipNew != nil {
		octetsNew := strings.Split(ipNew.String(), ".")
		if len(octetsNew) == 4 {
			fullReverseNew := fmt.Sprintf("%s.%s.%s.%s.in-addr.arpa.", octetsNew[3], octetsNew[2], octetsNew[1], octetsNew[0])
			bestZoneNew := findBestReverseZone(zones, fullReverseNew, ".in-addr.arpa")

			if bestZoneNew != "" {
				ttl := 3600

				patchBody := PDNSZonePatchRequest{
					RRSets: []PDNSRRSetChange{
						{
							Name:       fullReverseNew,
							Type:       "PTR",
							TTL:        &ttl,
							ChangeType: "REPLACE",
							Records: []Record{
								{
									Content:  fqdn,
									Disabled: false,
								},
							},
							Comments: []Comment{
								{
									Content:    "",
									Account:    "",
									ModifiedAt: now,
								},
							},
						},
					},
				}

				_, _ = httpc.R().SetContext(ctx).SetBody(patchBody).
					Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, bestZoneNew))
			}
		}
	}

	if ipOld != nil && (ipNew == nil || !ipNew.Equal(ipOld)) {
		octetsOld := strings.Split(ipOld.String(), ".")
		if len(octetsOld) != 4 {
			return
		}

		fullReverseOld := fmt.Sprintf("%s.%s.%s.%s.in-addr.arpa.", octetsOld[3], octetsOld[2], octetsOld[1], octetsOld[0])

		bestZoneOld := findBestReverseZone(zones, fullReverseOld, ".in-addr.arpa")
		if bestZoneOld == "" {
			return
		}

		if hasForwardRecordsForIPv4(ctx, httpc, zones, ipOld) {
			return
		}

		delBody := PDNSZonePatchRequest{
			RRSets: []PDNSRRSetChange{
				{
					Name:       fullReverseOld,
					Type:       "PTR",
					ChangeType: "DELETE",
				},
			},
		}

		_, _ = httpc.R().SetContext(ctx).SetBody(delBody).
			Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, bestZoneOld))
	}
}

func parseIPv6(v string) net.IP {
	raw := strings.TrimSpace(strings.TrimSuffix(v, "."))
	if raw == "" {
		return nil
	}

	ip := net.ParseIP(raw)
	if ip == nil || ip.To4() != nil {
		return nil
	}

	return ip.To16()
}

func fetchZones(ctx context.Context, httpc *resty.Client) ([]pdnsZoneInfo, bool) {
	resp, err := httpc.R().
		SetContext(ctx).
		Get(fmt.Sprintf("/api/v1/servers/%s/zones", env.C.DnsServerId))

	if err != nil || resp.IsStatusFailure() {
		return nil, false
	}

	var zones []pdnsZoneInfo

	err = json.Unmarshal(resp.Bytes(), &zones)
	if err != nil {
		return nil, false
	}

	return zones, true
}

func buildIPv6Reverse(ip net.IP) string {
	var hexBuilder strings.Builder

	for i := 0; i < 16; i += 2 {
		group := uint16(ip[i])<<8 | uint16(ip[i+1])
		fmt.Fprintf(&hexBuilder, "%04x", group)
	}

	hexStr := hexBuilder.String()

	var nibbles []string
	for i := len(hexStr) - 1; i >= 0; i-- {
		nibbles = append(nibbles, string(hexStr[i]))
	}

	return strings.Join(nibbles, ".") + ".ip6.arpa."
}

func replacePTR(ctx context.Context, httpc *resty.Client, zone, name, fqdn string, now int64) {
	ttl := 3600

	body := PDNSZonePatchRequest{
		RRSets: []PDNSRRSetChange{
			{
				Name:       name,
				Type:       "PTR",
				TTL:        &ttl,
				ChangeType: "REPLACE",
				Records: []Record{
					{
						Content:  fqdn,
						Disabled: false,
					},
				},
				Comments: []Comment{
					{
						Content:    "",
						Account:    "",
						ModifiedAt: now,
					},
				},
			},
		},
	}

	_, _ = httpc.R().
		SetContext(ctx).
		SetBody(body).
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, zone))
}

func deletePTR(ctx context.Context, httpc *resty.Client, zone, name string) {
	body := PDNSZonePatchRequest{
		RRSets: []PDNSRRSetChange{
			{
				Name:       name,
				Type:       "PTR",
				ChangeType: "DELETE",
			},
		},
	}

	_, _ = httpc.R().
		SetContext(ctx).
		SetBody(body).
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, zone))
}

func ensureReverseIPv6Update(ctx context.Context, vl, zona, name, oldVl string) {
	httpc := resty.New().
		SetTimeout(30*time.Second).
		SetBaseURL(env.C.DnsHost).
		SetHeader("X-API-Key", env.C.DnsApiKey).
		SetHeader("Accept", "application/json").
		SetRetryCount(2)

	ipNew := parseIPv6(vl)
	ipOld := parseIPv6(oldVl)

	if ipNew == nil && ipOld == nil {
		return
	}

	zones, ok := fetchZones(ctx, httpc)
	if !ok {
		return
	}

	now := time.Now().Unix()
	fqdn := buildFQDN(zona, name)

	if ipNew != nil {
		fullReverseNew := buildIPv6Reverse(ipNew)
		bestZoneNew := findBestReverseZone(zones, fullReverseNew, ".ip6.arpa")

		if bestZoneNew != "" {
			replacePTR(ctx, httpc, bestZoneNew, fullReverseNew, fqdn, now)
		}
	}

	if ipOld != nil && (ipNew == nil || !ipNew.Equal(ipOld)) {
		fullReverseOld := buildIPv6Reverse(ipOld)
		bestZoneOld := findBestReverseZone(zones, fullReverseOld, ".ip6.arpa")

		if bestZoneOld == "" {
			return
		}

		if hasForwardRecordsForIPv6(ctx, httpc, zones, ipOld) {
			return
		}

		deletePTR(ctx, httpc, bestZoneOld, fullReverseOld)
	}
}
