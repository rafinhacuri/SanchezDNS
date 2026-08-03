package records

import (
	"context"

	"resty.dev/v3"
)

func Delete(ctx context.Context, httpc *resty.Client, reg Registro) error {
	valor := NormalizarValor(reg)
	name := NomeCompleto(reg.Zone, reg.Name)

	zone, err := FetchZone(ctx, httpc, reg.Zone)
	if err != nil {
		return err
	}

	restantes, comentarios := registrosRestantes(zone, reg.Type, name, valor)

	if len(restantes) == 0 {
		return Patch(ctx, httpc, reg.Zone, PDNSZonePatchRequest{
			RRSets: []PDNSRRSetChange{
				{
					Name:       name,
					Type:       reg.Type,
					ChangeType: "DELETE",
				},
			},
		})
	}

	ttl := reg.TTL

	return Patch(ctx, httpc, reg.Zone, PDNSZonePatchRequest{
		RRSets: []PDNSRRSetChange{
			{
				Name:       name,
				Type:       reg.Type,
				TTL:        &ttl,
				ChangeType: "REPLACE",
				Records:    restantes,
				Comments:   comentarios,
			},
		},
	})
}
