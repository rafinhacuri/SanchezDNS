package zonas

import "strings"

func NormalizarTipo(tipo string) string {
	tipo = strings.ToLower(strings.TrimSpace(tipo))

	if tipo == "" || tipo == "forward" {
		return "normal"
	}

	return tipo
}
