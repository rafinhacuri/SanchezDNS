package records

import (
	"slices"
	"strings"
)

func Simplificar(zone Zone) ([]Simplified, *Soa) {
	var soa *Soa

	lista := []Simplified{}

	for _, rr := range zone.RRSets {
		if rr.Type == "SOA" {
			if len(rr.Records) > 0 && soa == nil {
				soa = ParseSoa(rr.Records[0].Content)
			}

			continue
		}

		for i, rec := range rr.Records {
			comentario, _, _ := Comentario(rr.Comments, i, len(rr.Records))

			simplificado, ok := simplificarRecord(zone.Name, rr, rec, comentario)
			if !ok {
				continue
			}

			lista = append(lista, simplificado)
		}
	}

	slices.SortFunc(lista, func(a, b Simplified) int {
		return strings.Compare(a.Name, b.Name)
	})

	return lista, soa
}
