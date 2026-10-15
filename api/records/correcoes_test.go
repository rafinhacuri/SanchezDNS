package records_test

import (
	"errors"
	"testing"

	"github.com/rafinhacuri/SanchezDNS/api/records"
)

func TestNomeCompletoRespeitaFronteiraDeRotulo(t *testing.T) {
	t.Parallel()

	casos := []struct{ zona, nome, esperado string }{
		{"curi.dev.br.", "www", "www.curi.dev.br."},
		{"curi.dev.br.", "curi.dev.br", "curi.dev.br."},
		{"curi.dev.br.", "www.curi.dev.br", "www.curi.dev.br."},
		{"curi.dev.br.", "xcuri.dev.br", "xcuri.dev.br.curi.dev.br."},
		{"curi.dev.br.", "www.outro.br.", "www.outro.br."},
	}

	for _, caso := range casos {
		obtido := records.NomeCompleto(caso.zona, caso.nome)
		if obtido != caso.esperado {
			t.Errorf("NomeCompleto(%q, %q) = %q, esperado %q", caso.zona, caso.nome, obtido, caso.esperado)
		}
	}
}

func TestPrioridadeMXAceitaOsDoisFormatos(t *testing.T) {
	t.Parallel()

	dez := 10

	casos := []struct {
		id       string
		reg      records.Registro
		valido   bool
		esperado string
	}{
		{"campo separado", records.Registro{Type: "MX", VL: "mail.ex.com", Priority: &dez}, true, "10 mail.ex.com."},
		{"prioridade no valor", records.Registro{Type: "MX", VL: "10 mail.ex.com."}, true, "10 mail.ex.com."},
		{
			"os dois, como a tela",
			records.Registro{Type: "MX", VL: "10 mail.ex.com.", Priority: &dez},
			true,
			"10 mail.ex.com.",
		},
		{"sem prioridade", records.Registro{Type: "MX", VL: "mail.ex.com"}, false, ""},
	}

	for _, caso := range casos {
		reg := caso.reg

		valido := records.PrioridadeMX(&reg)
		if valido != caso.valido {
			t.Errorf("%s: PrioridadeMX = %v, esperado %v", caso.id, valido, caso.valido)

			continue
		}

		if valido && records.NormalizarValor(reg) != caso.esperado {
			t.Errorf("%s: valor = %q, esperado %q", caso.id, records.NormalizarValor(reg), caso.esperado)
		}
	}
}

func TestDeleteSemRegistroCorrespondente(t *testing.T) {
	t.Parallel()

	_, httpc := novoPdns(t, diretaCom(rrset("mail.f.example.", "A", "192.0.2.1")))

	casos := []records.Registro{
		{Zone: zonaDireta, Type: "A", Name: "mail", VL: "192.0.2.99"},
		{Zone: zonaDireta, Type: "A", Name: "inexistente", VL: "192.0.2.1"},
	}

	for _, reg := range casos {
		err := records.Delete(t.Context(), httpc, reg)
		if !errors.Is(err, records.ErrRegistroNaoEncontrado) {
			t.Errorf("Delete(%s %s) = %v, esperado ErrRegistroNaoEncontrado", reg.Name, reg.VL, err)
		}
	}
}

func TestUpdateSemRegistroCorrespondente(t *testing.T) {
	t.Parallel()

	_, httpc := novoPdns(t, diretaCom(rrset("mail.f.example.", "A", "192.0.2.1")))

	antigo := records.Registro{Zone: zonaDireta, Type: "A", Name: "mail", VL: "192.0.2.99", TTL: 300}
	novo := records.Registro{Zone: zonaDireta, Type: "A", Name: "mail", VL: "192.0.2.2", TTL: 300}

	err := records.Update(t.Context(), httpc, novo, antigo)
	if !errors.Is(err, records.ErrRegistroNaoEncontrado) {
		t.Errorf("Update = %v, esperado ErrRegistroNaoEncontrado", err)
	}
}
