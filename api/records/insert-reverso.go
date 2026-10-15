package records

import (
	"context"
	"time"

	"resty.dev/v3"
)

func InsertReverso(
	parentCtx context.Context,
	httpc *resty.Client,
	tipo, vl, zona, name string,
) ([]MudancaReverso, error) {
	if tipo != "A" && tipo != "AAAA" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(parentCtx, 6*time.Second)
	defer cancel()

	ip := ParseIP(tipo, vl)
	if ip == nil {
		return nil, nil
	}

	zonas, err := FetchZonas(ctx, httpc)
	if err != nil {
		return nil, err
	}

	return gravarPTRDerivado(ctx, httpc, zonas, ip, NomeCompleto(zona, name))
}
