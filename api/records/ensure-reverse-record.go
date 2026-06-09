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

type pdnsZoneInfo struct {
	Name string `json:"name"`
}

func buildFQDN(zone, name string) string {
	z := strings.TrimSuffix(zone, ".")
	n := name

	if !strings.HasSuffix(n, ".") {
		if !strings.HasSuffix(n, z) {
			n = fmt.Sprintf("%s.%s.", n, z)
		} else {
			n += "."
		}
	}

	return n
}

func ensureReverseRecord(parentCtx context.Context, tipo, vl, zona, name string) error {
	if tipo != "A" && tipo != "AAAA" {
		return nil
	}

	ctx, cancel := context.WithTimeout(parentCtx, 6*time.Second)
	defer cancel()

	switch tipo {
	case "A":
		ensureReverseIPv4(ctx, vl, zona, name)
	case "AAAA":
		ensureReverseIPv6(ctx, vl, zona, name)
	}

	return nil
}

func ensureReverseIPv4(ctx context.Context, vl, zona, name string) {
	httpc := resty.New().
		SetTimeout(30*time.Second).
		SetBaseURL(env.C.DnsHost).
		SetHeader("X-API-Key", env.C.DnsApiKey).
		SetHeader("Accept", "application/json").
		SetRetryCount(2)

	ipStr := strings.TrimSpace(vl)
	ipStr = strings.TrimSuffix(ipStr, ".")

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return
	}

	octets := strings.Split(ipv4.String(), ".")
	if len(octets) != 4 {
		return
	}

	fullReverseName := fmt.Sprintf("%s.%s.%s.%s.in-addr.arpa.", octets[3], octets[2], octets[1], octets[0])

	zoneListResp, err := httpc.R().SetContext(ctx).Get(fmt.Sprintf("/api/v1/servers/%s/zones", env.C.DnsServerId))
	if err != nil || zoneListResp.IsError() {
		return
	}

	var zones []pdnsZoneInfo

	err = json.Unmarshal(zoneListResp.Bytes(), &zones)
	if err != nil {
		return
	}

	bestZone := ""

	fullRevNoDot := strings.TrimSuffix(fullReverseName, ".")

	for _, z := range zones {
		zNameNoDot := strings.TrimSuffix(z.Name, ".")

		if !strings.HasSuffix(zNameNoDot, ".in-addr.arpa") {
			continue
		}

		if strings.HasSuffix(fullRevNoDot, zNameNoDot) {
			if len(zNameNoDot) > len(strings.TrimSuffix(bestZone, ".")) {
				bestZone = z.Name
			}
		}
	}

	if bestZone == "" {
		return
	}

	zoneResp, err := httpc.R().SetContext(ctx).Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, bestZone))
	if err != nil || zoneResp.IsError() {
		return
	}

	var zoneData Zone

	err = json.Unmarshal(zoneResp.Bytes(), &zoneData)
	if err != nil {
		return
	}

	for _, rr := range zoneData.RRSets {
		if rr.Type == "PTR" && rr.Name == fullReverseName {
			return
		}
	}

	fqdn := buildFQDN(zona, name)
	ttl := 3600
	now := time.Now().Unix()

	patchBody := PDNSZonePatchRequest{
		RRSets: []PDNSRRSetChange{
			{
				Name:       fullReverseName,
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
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, bestZone))
}

func ensureReverseIPv6(ctx context.Context, vl, zona, name string) {
	httpc := resty.New().
		SetTimeout(30*time.Second).
		SetBaseURL(env.C.DnsHost).
		SetHeader("X-API-Key", env.C.DnsApiKey).
		SetHeader("Accept", "application/json").
		SetRetryCount(2)

	raw := strings.TrimSpace(vl)
	ipStr := strings.TrimSuffix(raw, ".")

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return
	}

	if ip.To4() != nil {
		return
	}

	ipv6 := ip.To16()
	if ipv6 == nil {
		return
	}

	var hexBuilder strings.Builder

	for i := 0; i < 16; i += 2 {
		group := uint16(ipv6[i])<<8 | uint16(ipv6[i+1])
		fmt.Fprintf(&hexBuilder, "%04x", group)
	}

	hexStr := hexBuilder.String()

	var nibbles []string
	for i := len(hexStr) - 1; i >= 0; i-- {
		nibbles = append(nibbles, string(hexStr[i]))
	}

	fullReverseName := strings.Join(nibbles, ".") + ".ip6.arpa."

	zoneListResp, err := httpc.R().SetContext(ctx).Get(fmt.Sprintf("/api/v1/servers/%s/zones", env.C.DnsServerId))
	if err != nil || zoneListResp.IsError() {
		return
	}

	var zones []pdnsZoneInfo

	err = json.Unmarshal(zoneListResp.Bytes(), &zones)
	if err != nil {
		return
	}

	bestZone := ""

	fullRevNoDot := strings.TrimSuffix(fullReverseName, ".")

	for _, z := range zones {
		zNameNoDot := strings.TrimSuffix(z.Name, ".")

		if !strings.HasSuffix(zNameNoDot, ".ip6.arpa") {
			continue
		}

		if strings.HasSuffix(fullRevNoDot, zNameNoDot) {
			if len(zNameNoDot) > len(strings.TrimSuffix(bestZone, ".")) {
				bestZone = z.Name
			}
		}
	}

	if bestZone == "" {
		return
	}

	zoneResp, err := httpc.R().SetContext(ctx).Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, bestZone))
	if err != nil || zoneResp.IsError() {
		return
	}

	var zoneData Zone

	err = json.Unmarshal(zoneResp.Bytes(), &zoneData)
	if err != nil {
		return
	}

	for _, rr := range zoneData.RRSets {
		if rr.Type == "PTR" && rr.Name == fullReverseName {
			return
		}
	}

	fqdn := buildFQDN(zona, name)
	ttl := 3600
	now := time.Now().Unix()

	patchBody := PDNSZonePatchRequest{
		RRSets: []PDNSRRSetChange{
			{
				Name:       fullReverseName,
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
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, bestZone))
}
