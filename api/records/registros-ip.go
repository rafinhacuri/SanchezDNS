package records

func RegistrosIP(zone Zone) []RegistroIP {
	lista := []RegistroIP{}

	for _, rr := range zone.RRSets {
		if rr.Type != "A" && rr.Type != "AAAA" {
			continue
		}

		for _, rec := range rr.Records {
			ip := ParseIP(rr.Type, rec.Content)
			if ip == nil {
				continue
			}

			lista = append(lista, RegistroIP{Name: rr.Name, Type: rr.Type, IP: ip})
		}
	}

	return lista
}
