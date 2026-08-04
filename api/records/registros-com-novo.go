package records

func registrosComNovo(rr *RRSet, valor string, ttl int) ([]Record, int) {
	if rr == nil {
		return []Record{{Content: valor, Disabled: false}}, ttl
	}

	registros := make([]Record, 0, len(rr.Records)+1)
	for _, rec := range rr.Records {
		registros = append(registros, Record{Content: rec.Content, Disabled: rec.Disabled})
	}

	if ttl <= 0 {
		ttl = rr.TTL
	}

	return append(registros, Record{Content: valor, Disabled: false}), ttl
}
