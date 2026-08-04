package records

import (
	"strconv"
	"strings"
)

func ParseSoa(content string) *Soa {
	parts := strings.Fields(content)
	if len(parts) < 7 {
		return nil
	}

	refresh, _ := strconv.Atoi(parts[3])
	retry, _ := strconv.Atoi(parts[4])
	expire, _ := strconv.Atoi(parts[5])
	negativeCacheTtl, _ := strconv.Atoi(parts[6])

	return &Soa{
		StartOfAuthority: parts[0],
		Email:            parts[1],
		Refresh:          refresh,
		Retry:            retry,
		Expire:           expire,
		NegativeCacheTtl: negativeCacheTtl,
	}
}
