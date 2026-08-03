package records

import (
	"context"
	"errors"

	"resty.dev/v3"
)

func Update(ctx context.Context, httpc *resty.Client, novo, antigo Registro) error {
	valorNovo := NormalizarValor(novo)
	valorAntigo := NormalizarValor(antigo)
	name := NomeCompleto(novo.Zone, novo.Name)

	zone, err := FetchZone(ctx, httpc, novo.Zone)
	if err != nil {
		return err
	}

	existente := FindRRSet(zone, novo.Type, name)
	if existente == nil {
		return errors.New("rrset não encontrado para edição")
	}

	anteriores := Comentarios(existente)

	ttl := existente.TTL
	if novo.TTL > 0 {
		ttl = novo.TTL
	}

	err = Patch(ctx, httpc, novo.Zone, PDNSZonePatchRequest{
		RRSets: []PDNSRRSetChange{
			{
				Name:       name,
				Type:       novo.Type,
				TTL:        &ttl,
				ChangeType: "REPLACE",
				Records:    registrosSubstituindo(existente, valorAntigo, valorNovo),
			},
		},
	})
	if err != nil {
		return err
	}

	return AplicarComentarios(ctx, httpc, novo, name, valorNovo, anteriores)
}
