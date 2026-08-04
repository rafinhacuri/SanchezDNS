package records

import (
	"context"

	"resty.dev/v3"
)

func semPTR(ctx context.Context, httpc *resty.Client, candidatos []ReversoAusente) ([]ReversoAusente, error) {
	existentes := map[string]bool{}
	visitadas := map[string]bool{}
	ausentes := []ReversoAusente{}

	for _, candidato := range candidatos {
		if !visitadas[candidato.ZonaReversa] {
			zone, err := FetchZone(ctx, httpc, candidato.ZonaReversa)
			if err != nil {
				return nil, err
			}

			for _, rr := range zone.RRSets {
				if rr.Type == "PTR" {
					existentes[rr.Name] = true
				}
			}

			visitadas[candidato.ZonaReversa] = true
		}

		if existentes[candidato.NomeReverso] {
			continue
		}

		ausentes = append(ausentes, candidato)
	}

	return ausentes, nil
}
