package records

import (
	"context"

	"resty.dev/v3"
)

func DeletePTR(ctx context.Context, httpc *resty.Client, zona, nomeReverso string) error {
	body := PDNSZonePatchRequest{
		RRSets: []PDNSRRSetChange{
			{
				Name:       nomeReverso,
				Type:       "PTR",
				ChangeType: "DELETE",
			},
		},
	}

	return Patch(ctx, httpc, zona, body)
}
