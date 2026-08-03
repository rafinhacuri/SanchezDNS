package records

func simplificarRecord(zona string, rr RRSet, rec Record, comentario string) (Simplified, bool) {
	switch rr.Type {
	case "MX":
		return simplificarMX(zona, rr, rec, comentario), true
	case "SRV":
		return simplificarSRV(zona, rr, rec, comentario)
	case "HTTPS":
		return simplificarHTTPS(zona, rr, rec, comentario)
	}

	return Simplified{
		Zone:    zona,
		Type:    rr.Type,
		Name:    rr.Name,
		VL:      rec.Content,
		TTL:     rr.TTL,
		Comment: comentario,
	}, true
}
