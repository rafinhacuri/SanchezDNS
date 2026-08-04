package records

import (
	"fmt"
	"strings"
)

func NomeCompleto(zona, name string) string {
	if strings.HasSuffix(name, ".") {
		return name
	}

	zone := strings.TrimSuffix(zona, ".")
	if strings.HasSuffix(name, zone) {
		return name + "."
	}

	return fmt.Sprintf("%s.%s.", name, zone)
}
