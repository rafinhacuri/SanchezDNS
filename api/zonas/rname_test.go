package zonas_test

import (
	"testing"

	"github.com/rafinhacuri/SanchezDNS/api/zonas"
)

func TestRname(t *testing.T) {
	t.Parallel()

	casos := []struct{ email, esperado string }{
		{"rafael@curi.dev.br", "rafael.curi.dev.br"},
		{"nome.sobrenome@exemplo.com.br", `nome\.sobrenome.exemplo.com.br`},
		{"hostmaster.exemplo.com.br", "hostmaster.exemplo.com.br"},
		{"hostmaster.exemplo.com.br.", "hostmaster.exemplo.com.br"},
	}

	for _, caso := range casos {
		obtido := zonas.Rname(caso.email)
		if obtido != caso.esperado {
			t.Errorf("Rname(%q) = %q, esperado %q", caso.email, obtido, caso.esperado)
		}
	}
}
