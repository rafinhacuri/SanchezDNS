package records

import (
	"context"
	"errors"
	"time"

	"resty.dev/v3"
)

func AplicarComentarios(
	ctx context.Context,
	httpc *resty.Client,
	reg Registro,
	name, conteudo string,
	anteriores map[string]string,
) error {
	zone, err := FetchZone(ctx, httpc, reg.Zone)
	if err != nil {
		return err
	}

	rr := FindRRSet(zone, reg.Type, name)
	if rr == nil {
		return errors.New("rrset não encontrado após a alteração")
	}

	agora := time.Now().Unix()
	registros := make([]Record, 0, len(rr.Records))
	comentarios := make([]Comment, 0, len(rr.Records))

	for _, rec := range rr.Records {
		registros = append(registros, Record{Content: rec.Content, Disabled: rec.Disabled})

		texto := ""

		if rec.Content == conteudo {
			texto = reg.Comment
		} else if anterior, ok := anteriores[rec.Content]; ok {
			texto = anterior
		}

		comentarios = append(comentarios, Comment{Content: texto, Account: "", ModifiedAt: agora})
	}

	ttl := rr.TTL

	return Patch(ctx, httpc, reg.Zone, PDNSZonePatchRequest{
		RRSets: []PDNSRRSetChange{
			{
				Name:       name,
				Type:       reg.Type,
				TTL:        &ttl,
				ChangeType: "REPLACE",
				Records:    registros,
				Comments:   comentarios,
			},
		},
	})
}
