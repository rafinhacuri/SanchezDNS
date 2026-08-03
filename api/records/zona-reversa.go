package records

import "strings"

func ZonaReversa(zonas []ZoneInfo, nomeReverso string) string {
	sufixo := ".in-addr.arpa"
	if strings.HasSuffix(nomeReverso, ".ip6.arpa.") {
		sufixo = ".ip6.arpa"
	}

	melhor := ""
	reverso := strings.TrimSuffix(nomeReverso, ".")

	for _, zona := range zonas {
		nome := strings.TrimSuffix(zona.Name, ".")

		if !strings.HasSuffix(nome, sufixo) || !strings.HasSuffix(reverso, nome) {
			continue
		}

		if len(nome) > len(strings.TrimSuffix(melhor, ".")) {
			melhor = zona.Name
		}
	}

	return melhor
}
