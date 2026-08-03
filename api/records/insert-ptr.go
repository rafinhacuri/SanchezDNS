package records

import (
	"context"
	"time"

	"resty.dev/v3"
)

func InsertPTR(ctx context.Context, httpc *resty.Client, zona, nomeReverso, fqdn string) error {
	ttl := 3600

	body := PDNSZonePatchRequest{
		RRSets: []PDNSRRSetChange{
			{
				Name:       nomeReverso,
				Type:       "PTR",
				TTL:        &ttl,
				ChangeType: "REPLACE",
				Records:    []Record{{Content: fqdn, Disabled: false}},
				Comments:   []Comment{{Content: "", Account: "", ModifiedAt: time.Now().Unix()}},
			},
		},
	}

	return Patch(ctx, httpc, zona, body)
}
