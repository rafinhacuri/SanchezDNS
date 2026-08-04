package records

import (
	"context"
	"time"

	"resty.dev/v3"
)

func UpdateReverso(parentCtx context.Context, httpc *resty.Client, tipo, vl, vlAntigo, zona, name string) error {
	if tipo != "A" && tipo != "AAAA" {
		return nil
	}

	ctx, cancel := context.WithTimeout(parentCtx, 6*time.Second)
	defer cancel()

	novo := ParseIP(tipo, vl)
	antigo := ParseIP(tipo, vlAntigo)

	if novo == nil && antigo == nil {
		return nil
	}

	zonas, err := FetchZonas(ctx, httpc)
	if err != nil {
		return err
	}

	if novo != nil {
		nomeReverso := NomeReverso(novo)

		zonaReversa := ZonaReversa(zonas, nomeReverso)
		if zonaReversa != "" {
			err = InsertPTR(ctx, httpc, zonaReversa, nomeReverso, NomeCompleto(zona, name))
			if err != nil {
				return err
			}
		}
	}

	if antigo == nil || (novo != nil && novo.Equal(antigo)) {
		return nil
	}

	nomeAntigo := NomeReverso(antigo)

	zonaAntiga := ZonaReversa(zonas, nomeAntigo)
	if zonaAntiga == "" {
		return nil
	}

	if TemForward(ctx, httpc, zonas, antigo) {
		return nil
	}

	return DeletePTR(ctx, httpc, zonaAntiga, nomeAntigo)
}
