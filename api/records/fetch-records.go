package records

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

type Simplified struct {
	Zone        string  `json:"zone"`
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	VL          string  `json:"vl"`
	TTL         int     `json:"ttl"`
	Comment     string  `json:"comment,omitempty"`
	SVCPriority *int    `json:"svcPriority,omitempty"`
	TargetName  *string `json:"targetName,omitempty"`
	SVCParams   *string `json:"svcParams,omitempty"`
	Weight      *int    `json:"weight,omitempty"`
	Port        *int    `json:"port,omitempty"`
	Target      *string `json:"target,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
}

type Soa struct {
	StartOfAuthority string `json:"startOfAuthority"`
	Email            string `json:"email"`
	Refresh          int    `json:"refresh"`
	Retry            int    `json:"retry"`
	Expire           int    `json:"expire"`
	NegativeCacheTtl int    `json:"negativeCacheTtl"`
}

func FetchRecords(ctx context.Context, zone string) ([]Simplified, *Soa, error) {
	httpc := resty.New().
		SetTimeout(30*time.Second).
		SetBaseURL(env.C.DnsHost).
		SetHeader("X-API-Key", env.C.DnsApiKey).
		SetHeader("Accept", "application/json").
		SetRetryCount(2)

	resp, err := httpc.R().SetContext(ctx).Get(fmt.Sprintf("/api/v1/servers/%s/zones/%s", env.C.DnsServerId, zone))
	if err != nil {
		return []Simplified{}, nil, errors.New("falha ao buscar registros: " + err.Error())
	}

	if resp.IsError() {
		return []Simplified{}, nil, errors.New("falha ao buscar registros:" + resp.String())
	}

	var z Zone

	err = json.Unmarshal(resp.Bytes(), &z)
	if err != nil {
		return []Simplified{}, nil, errors.New("falha ao analisar resposta do PowerDNS: " + err.Error())
	}

	var (
		records []Simplified
		soa     *Soa
	)

	getRecords(&records, z, &soa)

	sort.Slice(records, func(i, j int) bool {
		return records[i].Name < records[j].Name
	})

	return records, soa, nil
}

func appendMXRecord(records *[]Simplified, z Zone, rr rrsetRecord, rec Record, comment string) {
	var (
		priority *int
		value    string
	)

	parts := strings.Fields(rec.Content)
	value = rec.Content

	if len(parts) >= 2 {
		p, err := strconv.Atoi(parts[0])
		if err == nil {
			priority = &p

			host := strings.Join(parts[1:], " ")
			if !strings.HasSuffix(host, ".") {
				host += "."
			}

			value = fmt.Sprintf("%d %s", *priority, host)
		}
	}

	if priority == nil && !strings.HasSuffix(value, ".") {
		value += "."
	}

	*records = append(*records, Simplified{
		Zone:     z.Name,
		Type:     rr.Type,
		Name:     rr.Name,
		VL:       value,
		TTL:      rr.TTL,
		Comment:  comment,
		Priority: priority,
	})
}

func appendSRVRecord(records *[]Simplified, z Zone, rr rrsetRecord, rec Record, comment string) bool {
	parts := strings.Fields(rec.Content)
	if len(parts) < 4 {
		return false
	}

	p, _ := strconv.Atoi(parts[0])
	w, _ := strconv.Atoi(parts[1])
	port, _ := strconv.Atoi(parts[2])

	t := parts[3]
	if !strings.HasSuffix(t, ".") {
		t += "."
	}

	priority := p
	weight := w
	portVal := port
	target := t
	value := rec.Content

	*records = append(*records, Simplified{
		Zone:     z.Name,
		Type:     rr.Type,
		Name:     rr.Name,
		VL:       value,
		TTL:      rr.TTL,
		Comment:  comment,
		Priority: &priority,
		Weight:   &weight,
		Port:     &portVal,
		Target:   &target,
	})

	return true
}

func appendHTTPSRecord(records *[]Simplified, z Zone, rr rrsetRecord, rec Record, comment string) bool {
	parts := strings.Fields(rec.Content)
	if len(parts) < 3 {
		return false
	}

	svcP, _ := strconv.Atoi(parts[0])
	target := parts[1]
	params := strings.Join(parts[2:], " ")

	if !strings.HasSuffix(target, ".") {
		target += "."
	}

	svcPriority := svcP
	targetName := target
	svcParams := params
	value := rec.Content

	*records = append(*records, Simplified{
		Zone:        z.Name,
		Type:        rr.Type,
		Name:        rr.Name,
		VL:          value,
		TTL:         rr.TTL,
		Comment:     comment,
		SVCPriority: &svcPriority,
		TargetName:  &targetName,
		SVCParams:   &svcParams,
	})

	return true
}

func getRecords(records *[]Simplified, z Zone, soa **Soa) {
	for _, rr := range z.RRSets {
		if rr.Type == "SOA" && len(rr.Records) > 0 {
			parts := strings.Fields(rr.Records[0].Content)
			if len(parts) >= 7 {
				refresh, _ := strconv.Atoi(parts[3])
				retry, _ := strconv.Atoi(parts[4])
				expire, _ := strconv.Atoi(parts[5])
				negTTL, _ := strconv.Atoi(parts[6])

				*soa = &Soa{
					StartOfAuthority: parts[0],
					Email:            parts[1],
					Refresh:          refresh,
					Retry:            retry,
					Expire:           expire,
					NegativeCacheTtl: negTTL,
				}
			}

			continue
		}

		for i, rec := range rr.Records {
			comment, _, _ := findCommentForRecordContent(rr.Comments, i, len(rr.Records))

			switch rr.Type {
			case "MX":
				appendMXRecord(records, z, rr, rec, comment)

				continue

			case "SRV":
				if appendSRVRecord(records, z, rr, rec, comment) {
					continue
				}

			case "HTTPS":
				if appendHTTPSRecord(records, z, rr, rec, comment) {
					continue
				}

			default:
				*records = append(*records, Simplified{
					Zone:    z.Name,
					Type:    rr.Type,
					Name:    rr.Name,
					VL:      rec.Content,
					TTL:     rr.TTL,
					Comment: comment,
				})
			}
		}
	}
}
