package records

import (
	"context"
	"net"

	"resty.dev/v3"
)

func liberarPTRDerivado(
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

	rr := FindRRSet(zone, "PTR", nomeReverso)
	if !apontaPara(rr, fqdn) {
		return nil, nil
	}

	anterior := rr.Records[0].Content

	nomes, err := NomesComIP(ctx, httpc, zonas, ip)
	if err != nil {
		return nil, err
	}

	if len(nomes) > 0 {
		err = InsertPTR(ctx, httpc, zonaReversa, nomeReverso, nomes[0])
		if err != nil {
			return nil, err
		}

		return []MudancaReverso{{
			Acao:         AcaoReversoReapontado,
			ZonaReversa:  zonaReversa,
			NomeReverso:  nomeReverso,
			IP:           ip.String(),
			Alvo:         nomes[0],
			AlvoAnterior: anterior,
		}}, nil
	}

	err = DeletePTR(ctx, httpc, zonaReversa, nomeReverso)
	if err != nil {
		return nil, err
	}

	return []MudancaReverso{{
		Acao:         AcaoReversoRemovido,
		ZonaReversa:  zonaReversa,
		NomeReverso:  nomeReverso,
		IP:           ip.String(),
		AlvoAnterior: anterior,
	}}, nil
}
