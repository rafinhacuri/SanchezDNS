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

type pdnsRRSetDelete struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	ChangeType string `json:"changetype"`
}

type pdnsPatchRequest struct {
	RRSets []pdnsRRSetDelete `json:"rrsets"`
}

func deleteReverseRecord(parentCtx context.Context, tipo, vl string) error {
	if tipo != "A" && tipo != "AAAA" {
		return nil
	}

	ctx, cancel := context.WithTimeout(parentCtx, 6*time.Second)
	defer cancel()

	switch tipo {
	case "A":
		deleteIPv4(ctx, vl)
	case "AAAA":
		deleteIPv6(ctx, vl)
	}

	return nil
}

func deleteIPv4(ctx context.Context, vl string) {
	httpc := resty.New().
		SetTimeout(30*time.Second).
		SetBaseURL(env.C.DnsHost).
		SetHeader("X-API-Key", env.C.DnsApiKey).
		SetHeader("Accept", "application/json").
		SetRetryCount(2)

	ipStr := strings.TrimSpace(strings.TrimSuffix(vl, "."))

	ip := net.ParseIP(ipStr)
	if ip != nil {
		ip = ip.To4()
	}

	if ip == nil {
		return
	}

	zoneListResp, err := httpc.R().SetContext(ctx).Get(fmt.Sprintf("/api/v1/servers/%s/zones", env.C.DnsServerId))
	if err != nil || zoneListResp.IsError() {
		return
	}

	var zones []pdnsZoneInfo

	err = json.Unmarshal(zoneListResp.Bytes(), &zones)
	if err != nil {
		return
	}

	if hasForwardRecordsForIPv4(ctx, httpc, zones, ip) {
		return
	}

	octets := strings.Split(ip.String(), ".")
	if len(octets) != 4 {
		return
	}

	fullReverse := fmt.Sprintf("%s.%s.%s.%s.in-addr.arpa.", octets[3], octets[2], octets[1], octets[0])

	bestZone := findBestReverseZone(zones, fullReverse, ".in-addr.arpa")
	if bestZone == "" {
		return
	}

	delBody := pdnsPatchRequest{
		RRSets: []pdnsRRSetDelete{
			{
				Name:       fullReverse,
				Type:       "PTR",
				ChangeType: "DELETE",
			},
		},
	}

	_, _ = httpc.R().SetContext(ctx).SetBody(delBody).
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, bestZone))
}

func deleteIPv6(ctx context.Context, vl string) {
	httpc := resty.New().
		SetTimeout(30*time.Second).
		SetBaseURL(env.C.DnsHost).
		SetHeader("X-API-Key", env.C.DnsApiKey).
		SetHeader("Accept", "application/json").
		SetRetryCount(2)

	raw := strings.TrimSpace(strings.TrimSuffix(vl, "."))

	ip := net.ParseIP(raw)
	if ip != nil && ip.To4() != nil {
		ip = nil
	}

	if ip != nil {
		ip = ip.To16()
	}

	if ip == nil {
		return
	}

	zoneListResp, err := httpc.R().SetContext(ctx).Get(fmt.Sprintf("/api/v1/servers/%s/zones", env.C.DnsServerId))
	if err != nil || zoneListResp.IsError() {
		return
	}

	var zones []pdnsZoneInfo

	err = json.Unmarshal(zoneListResp.Bytes(), &zones)
	if err != nil {
		return
	}

	if hasForwardRecordsForIPv6(ctx, httpc, zones, ip) {
		return
	}

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

	fullReverse := strings.Join(nibbles, ".") + ".ip6.arpa."

	bestZone := findBestReverseZone(zones, fullReverse, ".ip6.arpa")
	if bestZone == "" {
		return
	}

	delBody := pdnsPatchRequest{
		RRSets: []pdnsRRSetDelete{
			{
				Name:       fullReverse,
				Type:       "PTR",
				ChangeType: "DELETE",
			},
		},
	}

	_, _ = httpc.R().SetContext(ctx).SetBody(delBody).
		Patch(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, bestZone))
}
