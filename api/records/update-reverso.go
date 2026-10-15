package records

import (
	"context"
	"net"
	"time"

	"resty.dev/v3"
)

func UpdateReverso(
	parentCtx context.Context,
	httpc *resty.Client,
	tipo, vl, vlAntigo, zona, name string,
) ([]MudancaReverso, error) {
	if tipo != "A" && tipo != "AAAA" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(parentCtx, 6*time.Second)
	defer cancel()

	novo := ParseIP(tipo, vl)
	antigo := ParseIP(tipo, vlAntigo)

	if novo == nil && antigo == nil {
		return nil, nil
	}

	zonas, err := FetchZonas(ctx, httpc)
	if err != nil {
		return nil, err
	}

	fqdn := NomeCompleto(zona, name)

	gravadas, err := gravarDoNovo(ctx, httpc, zonas, novo, fqdn)
	if err != nil {
		return gravadas, err
	}

	if antigo == nil || (novo != nil && novo.Equal(antigo)) {
		return gravadas, nil
	}

	liberadas, err := liberarPTRDerivado(ctx, httpc, zonas, antigo, fqdn)

	return append(gravadas, liberadas...), err
}

func gravarDoNovo(
	ctx context.Context,
	httpc *resty.Client,
	zonas []ZoneInfo,
	novo net.IP,
	fqdn string,
) ([]MudancaReverso, error) {
	if novo == nil {
		return nil, nil
	}

	return gravarPTRDerivado(ctx, httpc, zonas, novo, fqdn)
}
