package records

import (
	"context"
	"time"

	"resty.dev/v3"
)

func InsertReverso(parentCtx context.Context, httpc *resty.Client, tipo, vl, zona, name string) error {
	if tipo != "A" && tipo != "AAAA" {
		return nil
	}

	ctx, cancel := context.WithTimeout(parentCtx, 6*time.Second)
	defer cancel()

	ip := ParseIP(tipo, vl)
	if ip == nil {
		return nil
	}

	nomeReverso := NomeReverso(ip)
	if nomeReverso == "" {
		return nil
	}

	zonas, err := FetchZonas(ctx, httpc)
	if err != nil {
		return err
	}

	zonaReversa := ZonaReversa(zonas, nomeReverso)
	if zonaReversa == "" {
		return nil
	}

	zone, err := FetchZone(ctx, httpc, zonaReversa)
	if err != nil {
		return err
	}

	if FindRRSet(zone, "PTR", nomeReverso) != nil {
		return nil
	}

	return InsertPTR(ctx, httpc, zonaReversa, nomeReverso, NomeCompleto(zona, name))
}
