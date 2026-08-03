package users

func NivelPermissao(permissao string) string {
	switch permissao {
	case "escrita":
		return "ESCRITA"
	case "leitura":
		return "LEITURA"
	}

	return ""
}
