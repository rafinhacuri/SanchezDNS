package records

import (
	"context"

	"resty.dev/v3"
)

func Insert(ctx context.Context, httpc *resty.Client, reg Registro) error {
	valor := NormalizarValor(reg)
	name := NomeCompleto(reg.Zone, reg.Name)

	zone, err := FetchZone(ctx, httpc, reg.Zone)
	if err != nil {
		return err
	}

	existente := FindRRSet(zone, reg.Type, name)
	anteriores := Comentarios(existente)

	registros, ttl := registrosComNovo(existente, valor, reg.TTL)

	err = Patch(ctx, httpc, reg.Zone, PDNSZonePatchRequest{
		RRSets: []PDNSRRSetChange{
			{
				Name:       name,
				Type:       reg.Type,
				TTL:        &ttl,
				ChangeType: "REPLACE",
				Records:    registros,
			},
		},
	})
	if err != nil {
		return err
	}

	return AplicarComentarios(ctx, httpc, reg, name, valor, anteriores)
}
