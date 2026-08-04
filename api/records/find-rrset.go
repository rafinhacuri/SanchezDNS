package records

func FindRRSet(zone Zone, tipo, name string) *RRSet {
	for i := range zone.RRSets {
		rr := &zone.RRSets[i]
		if rr.Type == tipo && rr.Name == name {
			return rr
		}
	}

	return nil
}
