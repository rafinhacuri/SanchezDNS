package records_test

import (
	"testing"

	"github.com/rafinhacuri/SanchezDNS/api/records"
)

func TestNomeReversoV4EV6(t *testing.T) {
	t.Parallel()

	casos := []struct {
		tipo, ip, esperado string
	}{
		{"A", "152.84.120.224", "224.120.84.152.in-addr.arpa."},
		{"A", "10.0.1.5.", "5.1.0.10.in-addr.arpa."},
		{
			"AAAA",
			"2001:db8:85a3::8a2e:370:7334",
			"4.3.3.7.0.7.3.0.e.2.a.8.0.0.0.0.0.0.0.0.3.a.5.8.8.b.d.0.1.0.0.2.ip6.arpa.",
		},
		{"AAAA", "152.84.120.224", ""},
		{"A", "nao-e-ip", ""},
	}

	for _, caso := range casos {
		ip := records.ParseIP(caso.tipo, caso.ip)

		obtido := ""
		if ip != nil {
			obtido = records.NomeReverso(ip)
		}

		if obtido != caso.esperado {
			t.Errorf("ParseIP+NomeReverso(%s, %q) = %q, esperado %q", caso.tipo, caso.ip, obtido, caso.esperado)
		}
	}
}

func TestIPDoReversoInverteNomeReverso(t *testing.T) {
	t.Parallel()

	casos := []string{"152.84.120.224", "10.0.1.5", "2001:db8:85a3::8a2e:370:7334", "::1"}

	for _, original := range casos {
		tipo := "AAAA"
		if records.ParseIP("A", original) != nil {
			tipo = "A"
		}

		ip := records.ParseIP(tipo, original)

		volta := records.IPDoReverso(records.NomeReverso(ip))
		if volta == nil || !volta.Equal(ip) {
			t.Errorf("IPDoReverso(NomeReverso(%q)) = %v, esperado %v", original, volta, ip)
		}
	}

	if records.IPDoReverso("nao.e.reverso.") != nil {
		t.Error("nome inválido deveria retornar nil")
	}

	if records.IPDoReverso("1.2.in-addr.arpa.") != nil {
		t.Error("reverso incompleto deveria retornar nil")
	}
}

func TestRegistrosIPFiltraTipos(t *testing.T) {
	t.Parallel()

	zone := records.Zone{
		RRSets: []records.RRSet{
			{Type: "A", Name: "host.sanchez.br.", Records: []records.Record{{Content: "10.0.1.5"}}},
			{Type: "AAAA", Name: "v6.sanchez.br.", Records: []records.Record{{Content: "2001:db8::1"}}},
			{Type: "CNAME", Name: "www.sanchez.br.", Records: []records.Record{{Content: "host.sanchez.br."}}},
			{Type: "MX", Name: "sanchez.br.", Records: []records.Record{{Content: "10 mail.sanchez.br."}}},
		},
	}

	lista := records.RegistrosIP(zone)
	if len(lista) != 2 {
		t.Fatalf("esperado 2 registros (A e AAAA), obtido %d", len(lista))
	}

	if lista[0].Type != "A" || lista[1].Type != "AAAA" {
		t.Errorf("tipos inesperados: %s, %s", lista[0].Type, lista[1].Type)
	}
}

func TestZonaReversaV4EV6(t *testing.T) {
	t.Parallel()

	zonas := []records.ZoneInfo{
		{Name: "sanchez.br."},
		{Name: "84.152.in-addr.arpa."},
		{Name: "120.84.152.in-addr.arpa."},
		{Name: "0.0.0.0.3.a.5.8.8.b.d.0.1.0.0.2.ip6.arpa."},
	}

	v4 := records.ZonaReversa(zonas, "224.120.84.152.in-addr.arpa.")
	if v4 != "120.84.152.in-addr.arpa." {
		t.Errorf("v4 = %q, esperado a zona mais específica", v4)
	}

	v6 := records.ZonaReversa(zonas, "4.3.3.7.0.7.3.0.e.2.a.8.0.0.0.0.0.0.0.0.3.a.5.8.8.b.d.0.1.0.0.2.ip6.arpa.")
	if v6 != "0.0.0.0.3.a.5.8.8.b.d.0.1.0.0.2.ip6.arpa." {
		t.Errorf("v6 = %q", v6)
	}

	if records.ZonaReversa(zonas, "1.2.3.4.in-addr.arpa.") != "" {
		t.Error("sem zona reversa cadastrada deveria retornar vazio")
	}
}

func TestNomeCompleto(t *testing.T) {
	t.Parallel()

	casos := []struct {
		zona, name, esperado string
	}{
		{"sanchez.br.", "www", "www.sanchez.br."},
		{"sanchez.br.", "www.sanchez.br", "www.sanchez.br."},
		{"sanchez.br.", "www.sanchez.br.", "www.sanchez.br."},
		{"sanchez.br", "mail", "mail.sanchez.br."},
	}

	for _, caso := range casos {
		obtido := records.NomeCompleto(caso.zona, caso.name)
		if obtido != caso.esperado {
			t.Errorf("NomeCompleto(%q, %q) = %q, esperado %q", caso.zona, caso.name, obtido, caso.esperado)
		}
	}
}
