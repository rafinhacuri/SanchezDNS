package records

import (
	"context"

	"resty.dev/v3"
)

func InsertReversos(ctx context.Context, httpc *resty.Client, ausentes []ReversoAusente) error {
	for _, ausente := range ausentes {
		err := InsertPTR(ctx, httpc, ausente.ZonaReversa, ausente.NomeReverso, ausente.Name)
		if err != nil {
			return err
		}
	}

	return nil
}
