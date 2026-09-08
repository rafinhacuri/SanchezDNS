package zonas

import "strings"

const (
	DsOk           = "ok"
	DsAusente      = "ausente"
	DsDivergente   = "divergente"
	DsIndisponivel = "indisponivel"
)

func normalizarDs(ds string) string {
	return strings.ToLower(strings.Join(strings.Fields(ds), " "))
}

func StatusDs(keys []DnssecKey, dsPai []string) string {
	if len(dsPai) == 0 {
		return DsAusente
	}

	publicados := make(map[string]bool, len(keys))

	for _, key := range keys {
		for _, ds := range key.Ds {
			publicados[normalizarDs(ds)] = true
		}
	}

	for _, ds := range dsPai {
		if publicados[normalizarDs(ds)] {
			return DsOk
		}
	}

	return DsDivergente
}
