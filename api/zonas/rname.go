package zonas

import "strings"

func Rname(email string) string {
	email = strings.TrimSuffix(strings.TrimSpace(email), ".")

	local, dominio, ok := strings.Cut(email, "@")
	if !ok {
		return email
	}

	return strings.ReplaceAll(local, ".", `\.`) + "." + dominio
}
