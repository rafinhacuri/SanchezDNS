package records

func registrosRestantes(zone Zone, tipo, name, valor string) ([]Record, []Comment) {
	var (
		registros   []Record
		comentarios []Comment
	)

	for _, rr := range zone.RRSets {
		if rr.Type != tipo || rr.Name != name {
			continue
		}

		for i, rec := range rr.Records {
			if rec.Content == valor {
				continue
			}

			registros = append(registros, Record{Content: rec.Content, Disabled: rec.Disabled})

			texto, conta, modificado := Comentario(rr.Comments, i, len(rr.Records))
			comentarios = append(comentarios, Comment{Content: texto, Account: conta, ModifiedAt: modificado})
		}
	}

	return registros, comentarios
}
