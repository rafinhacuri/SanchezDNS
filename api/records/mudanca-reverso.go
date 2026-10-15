package records

import "fmt"

const (
	AcaoReversoCriado     = "auto_insert_reverse"
	AcaoReversoReapontado = "auto_update_reverse"
	AcaoReversoRemovido   = "auto_delete_reverse"
)

type MudancaReverso struct {
	Acao         string
	ZonaReversa  string
	NomeReverso  string
	IP           string
	Alvo         string
	AlvoAnterior string
}

func (m MudancaReverso) Descricao() string {
	switch m.Acao {
	case AcaoReversoCriado:
		return fmt.Sprintf("Criado PTR %s (%s) apontando para %s", m.NomeReverso, m.IP, m.Alvo)
	case AcaoReversoReapontado:
		return fmt.Sprintf("Reapontado PTR %s (%s) de %s para %s", m.NomeReverso, m.IP, m.AlvoAnterior, m.Alvo)
	case AcaoReversoRemovido:
		return fmt.Sprintf("Removido PTR %s (%s) que apontava para %s", m.NomeReverso, m.IP, m.AlvoAnterior)
	}

	return ""
}
