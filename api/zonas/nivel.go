package zonas

import "github.com/rafinhacuri/SanchezDNS/api/users"

func Nivel(level, email, zona string, permissoes map[string]users.Zone) (string, bool) {
	if level == "admin" {
		return "ADMINISTRADOR", true
	}

	zone, ok := permissoes[zona]
	if !ok {
		return "", false
	}

	nivel := users.NivelPermissao(users.Permissao(zone, email))

	return nivel, nivel != ""
}
