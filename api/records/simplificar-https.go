package records

import (
	"strconv"
	"strings"
)

func simplificarHTTPS(zona string, rr RRSet, rec Record, comentario string) (Simplified, bool) {
	parts := strings.Fields(rec.Content)
	if len(parts) < 3 {
		return Simplified{}, false
	}

	svcPriority, _ := strconv.Atoi(parts[0])

	targetName := parts[1]
	if !strings.HasSuffix(targetName, ".") {
		targetName += "."
	}

	svcParams := strings.Join(parts[2:], " ")

	return Simplified{
		Zone:        zona,
		Type:        rr.Type,
		Name:        rr.Name,
		VL:          rec.Content,
		TTL:         rr.TTL,
		Comment:     comentario,
		SVCPriority: &svcPriority,
		TargetName:  &targetName,
		SVCParams:   &svcParams,
	}, true
}
