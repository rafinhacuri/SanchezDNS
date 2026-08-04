package records

import (
	"strconv"
	"strings"
)

func simplificarSRV(zona string, rr RRSet, rec Record, comentario string) (Simplified, bool) {
	parts := strings.Fields(rec.Content)
	if len(parts) < 4 {
		return Simplified{}, false
	}

	priority, _ := strconv.Atoi(parts[0])
	weight, _ := strconv.Atoi(parts[1])
	port, _ := strconv.Atoi(parts[2])

	target := parts[3]
	if !strings.HasSuffix(target, ".") {
		target += "."
	}

	return Simplified{
		Zone:     zona,
		Type:     rr.Type,
		Name:     rr.Name,
		VL:       rec.Content,
		TTL:      rr.TTL,
		Comment:  comentario,
		Priority: &priority,
		Weight:   &weight,
		Port:     &port,
		Target:   &target,
	}, true
}
