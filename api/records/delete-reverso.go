package records

import (
	"context"
	"time"

	"resty.dev/v3"
)

func DeleteReverso(parentCtx context.Context, httpc *resty.Client, tipo, vl string) error {
	if tipo != "A" && tipo != "AAAA" {
		return nil
	}

	ctx, cancel := context.WithTimeout(parentCtx, 6*time.Second)
	defer cancel()

	ip := ParseIP(tipo, vl)
	if ip == nil {
		return nil
	}

	zonas, err := FetchZonas(ctx, httpc)
	if err != nil {
		return err
	}

	if TemForward(ctx, httpc, zonas, ip) {
		return nil
	}

	nomeReverso := NomeReverso(ip)

	zonaReversa := ZonaReversa(zonas, nomeReverso)
	if zonaReversa == "" {
		return nil
	}

	return DeletePTR(ctx, httpc, zonaReversa, nomeReverso)
}
