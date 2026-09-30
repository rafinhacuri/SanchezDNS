package records_test

import (
	"net"
	"testing"

	"github.com/rafinhacuri/SanchezDNS/api/records"
)

// Casos negativos do episódio E2 da monografia (Quadro 12).

func TestNormalizacaoRDATA(t *testing.T) {
	t.Parallel()

	casos := []struct {
		id       string
		reg      records.Registro
		esperado string
	}{
		{"N1", records.Registro{Type: "CNAME", VL: "www.destino.com.br"}, "www.destino.com.br."},
		{"N2", records.Registro{Type: "NS", VL: "ns1.exemplo.com.br"}, "ns1.exemplo.com.br."},
		{"N3", records.Registro{Type: "MX", VL: "mail.exemplo.com.br", Priority: new(10)}, "10 mail.exemplo.com.br."},
		{"N4", records.Registro{Type: "TXT", VL: "v=spf1 -all"}, `"v=spf1 -all"`},
		{
			"N5",
			records.Registro{Type: "SRV", Priority: new(10), Weight: new(5), Port: new(5060), Target: "sip.exemplo.com.br"},
			"10 5 5060 sip.exemplo.com.br.",
		},
		{"N6", records.Registro{Type: "HTTPS", SvcPriority: new(1)}, "1 . alpn=h2"},
		{"N7", records.Registro{Type: "CAA", VL: "letsencrypt.org"}, `0 issue "letsencrypt.org"`},
	}

	for _, caso := range casos {
		obtido := records.NormalizarValor(caso.reg)
		if obtido != caso.esperado {
			t.Errorf("%s: NormalizarValor = %q, esperado %q", caso.id, obtido, caso.esperado)
		}
	}
}

func TestQualificacaoDoNome(t *testing.T) {
	t.Parallel()

	casos := []struct{ id, name string }{
		{"N8", "www"},
		{"N9", "www.exemplo.com.br"},
		{"N10", "www.exemplo.com.br."},
	}

	for _, caso := range casos {
		obtido := records.NomeCompleto("exemplo.com.br.", caso.name)
		if obtido != "www.exemplo.com.br." {
			t.Errorf("%s: NomeCompleto(%q) = %q", caso.id, caso.name, obtido)
		}
	}
}

func TestEnderecoInvalido(t *testing.T) {
	t.Parallel()

	if records.ParseIP("A", "999.1.1.1") != nil {
		t.Error("N11: endereço A inválido deveria ser rejeitado")
	}

	if records.ParseIP("AAAA", "192.0.2.10") != nil {
		t.Error("N12: AAAA com endereço IPv4 deveria ser rejeitado")
	}
}

func TestTTLMinimo(t *testing.T) {
	t.Parallel()

	if records.TTLValido(30) {
		t.Error("N13: TTL de 30 segundos deveria ser bloqueado")
	}

	if !records.TTLValido(60) || !records.TTLValido(3600) || !records.TTLValido(0) {
		t.Error("TTL de 60, 3600 ou não informado (0) deveria ser aceito")
	}
}

func TestDerivacaoDoNomeReverso(t *testing.T) {
	t.Parallel()

	if obtido := records.NomeReverso(net.ParseIP("192.0.2.10")); obtido != "10.2.0.192.in-addr.arpa." {
		t.Errorf("P1: %q", obtido)
	}

	esperado := "1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa."
	if obtido := records.NomeReverso(net.ParseIP("2001:db8::1")); obtido != esperado {
		t.Errorf("P2: %q", obtido)
	}
}

func TestEscolhaDaZonaReversa(t *testing.T) {
	t.Parallel()

	zonas := []records.ZoneInfo{
		{Name: "0.192.in-addr.arpa."},
		{Name: "2.0.192.in-addr.arpa."},
		{Name: "exemplo.com.br."},
	}

	casos := []struct{ id, nome, esperado string }{
		{"P3", "10.2.0.192.in-addr.arpa.", "2.0.192.in-addr.arpa."},
		{"P4", "5.12.0.192.in-addr.arpa.", "0.192.in-addr.arpa."},
		{"P5", "7.100.51.198.in-addr.arpa.", ""},
	}

	for _, caso := range casos {
		obtido := records.ZonaReversa(zonas, caso.nome)
		if obtido != caso.esperado {
			t.Errorf("%s: ZonaReversa(%q) = %q, esperado %q", caso.id, caso.nome, obtido, caso.esperado)
		}
	}
}

func TestDelegacaoClassless(t *testing.T) {
	t.Parallel()
	t.Skip("P6: delegação classless (RFC 2317) não é suportada; limitação registrada na monografia")

	zonas := []records.ZoneInfo{{Name: "2.0.192.in-addr.arpa."}, {Name: "128/25.2.0.192.in-addr.arpa."}}
	if obtido := records.ZonaReversa(zonas, "130.2.0.192.in-addr.arpa."); obtido != "128/25.2.0.192.in-addr.arpa." {
		t.Errorf("P6: %q", obtido)
	}
}
