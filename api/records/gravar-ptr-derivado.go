package records

import (
	"context"
	"net"

	"resty.dev/v3"
)

func gravarPTRDerivado(
	ctx context.Context,
	httpc *resty.Client,
	zonas []ZoneInfo,
	ip net.IP,
	fqdn string,
) ([]MudancaReverso, error) {
	nomeReverso := NomeReverso(ip)
	if nomeReverso == "" {
		return nil, nil
	}

	zonaReversa := ZonaReversa(zonas, nomeReverso)
	if zonaReversa == "" {
		return nil, nil
	}

	zone, err := FetchZone(ctx, httpc, zonaReversa)
	if err != nil {
		return nil, err
	}

	if FindRRSet(zone, "PTR", nomeReverso) != nil {
		return nil, nil
	}

	err = InsertPTR(ctx, httpc, zonaReversa, nomeReverso, fqdn)
	if err != nil {
		return nil, err
	}

	return []MudancaReverso{{
		Acao:        AcaoReversoCriado,
		ZonaReversa: zonaReversa,
		NomeReverso: nomeReverso,
		IP:          ip.String(),
		Alvo:        fqdn,
	}}, nil
}
