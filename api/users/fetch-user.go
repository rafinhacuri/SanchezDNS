package users

func FetchUser(zone Zone, permissao string, index int) (string, bool) {
	lista := zone.Leitura
	if permissao == "escrita" {
		lista = zone.Escrita
	}

	if index < 0 || index >= len(lista) {
		return "", false
	}

	return lista[index], true
}
