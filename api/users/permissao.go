package users

import (
	"slices"
	"strings"
)

func Permissao(zone Zone, email string) string {
	igual := func(u string) bool {
		return strings.EqualFold(strings.TrimSpace(u), strings.TrimSpace(email))
	}

	if slices.ContainsFunc(zone.Escrita, igual) {
		return "escrita"
	}

	if slices.ContainsFunc(zone.Leitura, igual) {
		return "leitura"
	}

	return ""
}
