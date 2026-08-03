package records

import (
	"context"

	"resty.dev/v3"
)

func FetchReversosAusentes(ctx context.Context, httpc *resty.Client, zona string) ([]ReversoAusente, error) {
	zone, err := FetchZone(ctx, httpc, zona)
	if err != nil {
		return nil, err
	}

	zonas, err := FetchZonas(ctx, httpc)
	if err != nil {
		return nil, err
	}

	vistos := map[string]bool{}
	candidatos := []ReversoAusente{}

	for _, registro := range RegistrosIP(zone) {
		nomeReverso := NomeReverso(registro.IP)
		if nomeReverso == "" || vistos[nomeReverso] {
			continue
		}

		zonaReversa := ZonaReversa(zonas, nomeReverso)
		if zonaReversa == "" {
			continue
		}

		vistos[nomeReverso] = true

		candidatos = append(candidatos, ReversoAusente{
			Name:        registro.Name,
			Type:        registro.Type,
			IP:          registro.IP.String(),
			NomeReverso: nomeReverso,
			ZonaReversa: zonaReversa,
		})
	}

	return semPTR(ctx, httpc, candidatos)
}
