package records

import (
	"fmt"
	"strconv"
	"strings"
)

func simplificarMX(zona string, rr RRSet, rec Record, comentario string) Simplified {
	var priority *int

	valor := rec.Content
	parts := strings.Fields(rec.Content)

	if len(parts) >= 2 {
		p, err := strconv.Atoi(parts[0])
		if err == nil {
			priority = &p

			host := strings.Join(parts[1:], " ")
			if !strings.HasSuffix(host, ".") {
				host += "."
			}

			valor = fmt.Sprintf("%d %s", p, host)
		}
	}

	if priority == nil && !strings.HasSuffix(valor, ".") {
		valor += "."
	}

	return Simplified{
		Zone:     zona,
		Type:     rr.Type,
		Name:     rr.Name,
		VL:       valor,
		TTL:      rr.TTL,
		Comment:  comentario,
		Priority: priority,
	}
}
