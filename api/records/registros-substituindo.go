package records

func registrosSubstituindo(rr *RRSet, antigo, novo string) []Record {
	registros := make([]Record, 0, len(rr.Records))

	for _, rec := range rr.Records {
		conteudo := rec.Content
		if conteudo == antigo {
			conteudo = novo
		}

		registros = append(registros, Record{Content: conteudo, Disabled: rec.Disabled})
	}

	return registros
}
