package statistics

func Valores(stats []Stat) map[string]string {
	valores := make(map[string]string, len(stats))
	for _, stat := range stats {
		valores[stat.Name] = valor(stat.Value)
	}

	return valores
}
