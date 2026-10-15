package logs

import "github.com/rafinhacuri/SanchezDNS/api/records"

func ReversosDerivados(username, origem string, mudancas []records.MudancaReverso) {
	for _, m := range mudancas {
		Insert(m.ZonaReversa, username, m.Acao, m.Descricao()+", derivado de "+origem)
	}
}
