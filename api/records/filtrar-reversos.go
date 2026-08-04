package records

func FiltrarReversos(ausentes []ReversoAusente, nomes []string) []ReversoAusente {
	escolhidos := make(map[string]bool, len(nomes))
	for _, nome := range nomes {
		escolhidos[nome] = true
	}

	filtrados := make([]ReversoAusente, 0, len(nomes))

	for _, ausente := range ausentes {
		if escolhidos[ausente.NomeReverso] {
			filtrados = append(filtrados, ausente)
		}
	}

	return filtrados
}
