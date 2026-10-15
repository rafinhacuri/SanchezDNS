package records

import (
	"strconv"
	"strings"
)

func PrioridadeMX(reg *Registro) bool {
	partes := strings.Fields(reg.VL)

	if len(partes) >= 2 {
		prioridade, err := strconv.Atoi(partes[0])
		if err == nil {
			reg.Priority = &prioridade
			reg.VL = strings.Join(partes[1:], " ")

			return true
		}
	}

	return len(partes) == 1 && reg.Priority != nil
}
