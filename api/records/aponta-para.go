package records

import "strings"

func apontaPara(rr *RRSet, fqdn string) bool {
	if rr == nil || len(rr.Records) != 1 {
		return false
	}

	alvo := strings.TrimSuffix(strings.TrimSpace(rr.Records[0].Content), ".")

	return strings.EqualFold(alvo, strings.TrimSuffix(fqdn, "."))
}
