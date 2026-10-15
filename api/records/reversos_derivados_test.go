package records_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"resty.dev/v3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
	"github.com/rafinhacuri/SanchezDNS/api/records"
)

const (
	zonaDireta  = "f.example."
	zonaReversa = "2.0.192.in-addr.arpa."
	ptrOutra    = "servidor.outra-equipe.example."
)

type pdnsFalso struct {
	mu      sync.Mutex
	zonas   map[string]*records.Zone
	falhar  map[string]bool
	recusar bool
}

func TestMain(m *testing.M) {
	env.C = &env.EnvConfig{DnsServerId: "localhost"}

	os.Exit(m.Run())
}

func novoPdns(t *testing.T, zonas ...records.Zone) (*pdnsFalso, *resty.Client) {
	t.Helper()

	p := &pdnsFalso{zonas: map[string]*records.Zone{}, falhar: map[string]bool{}}
	for i := range zonas {
		p.zonas[zonas[i].Name] = &zonas[i]
	}

	srv := httptest.NewServer(p)
	t.Cleanup(srv.Close)

	return p, resty.New().SetBaseURL(srv.URL)
}

func (p *pdnsFalso) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mu.Lock()
	defer p.mu.Unlock()

	const prefixo = "/api/v1/servers/localhost/zones"

	nome := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, prefixo), "/")

	switch {
	case nome == "" && r.Method == http.MethodGet:
		lista := make([]records.ZoneInfo, 0, len(p.zonas))
		for n := range p.zonas {
			lista = append(lista, records.ZoneInfo{Name: n})
		}

		responder(w, lista)
	case p.falhar[nome]:
		http.Error(w, `{"error": "falha simulada"}`, http.StatusInternalServerError)
	case r.Method == http.MethodGet:
		responder(w, p.zonas[nome])
	case r.Method == http.MethodPatch:
		if p.recusar {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`{"error": "Name is out of zone"}`))

			return
		}

		var body records.PDNSZonePatchRequest

		_ = json.NewDecoder(r.Body).Decode(&body)
		p.aplicar(p.zonas[nome], body)
		w.WriteHeader(http.StatusNoContent)
	}
}

func responder(w http.ResponseWriter, v any) {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *pdnsFalso) aplicar(zone *records.Zone, body records.PDNSZonePatchRequest) {
	for _, ch := range body.RRSets {
		restantes := zone.RRSets[:0]
		for _, rr := range zone.RRSets {
			if rr.Name != ch.Name || rr.Type != ch.Type {
				restantes = append(restantes, rr)
			}
		}

		zone.RRSets = restantes
		if ch.ChangeType == "REPLACE" {
			zone.RRSets = append(zone.RRSets, records.RRSet{Name: ch.Name, Type: ch.Type, Records: ch.Records})
		}
	}
}

func (p *pdnsFalso) ptr(nome string) string {
	p.mu.Lock()
	defer p.mu.Unlock()

	rr := records.FindRRSet(*p.zonas[zonaReversa], "PTR", nome)
	if rr == nil {
		return ""
	}

	return rr.Records[0].Content
}

func rrset(nome, tipo, conteudo string) records.RRSet {
	return records.RRSet{Name: nome, Type: tipo, Records: []records.Record{{Content: conteudo}}}
}

func diretaCom(rrs ...records.RRSet) records.Zone {
	return records.Zone{Name: zonaDireta, RRSets: rrs}
}

func reversaCom(rrs ...records.RRSet) records.Zone {
	return records.Zone{Name: zonaReversa, RRSets: rrs}
}

func TestInsertReversoCriaPTR(t *testing.T) {
	t.Parallel()

	p, httpc := novoPdns(t, diretaCom(rrset("a3.f.example.", "A", "192.0.2.10")), reversaCom())

	mudancas, err := records.InsertReverso(t.Context(), httpc, "A", "192.0.2.10", zonaDireta, "a3")
	if err != nil {
		t.Fatal(err)
	}

	if got := p.ptr("10.2.0.192.in-addr.arpa."); got != "a3.f.example." {
		t.Errorf("PTR = %q, esperado a3.f.example.", got)
	}

	if len(mudancas) != 1 || mudancas[0].Acao != records.AcaoReversoCriado {
		t.Errorf("mudanças = %+v, esperado uma criação", mudancas)
	}
}

func TestUpdateReversoNaoSobrescrevePTRDeOutraEquipe(t *testing.T) {
	t.Parallel()

	p, httpc := novoPdns(t,
		diretaCom(rrset("a8.f.example.", "A", "192.0.2.50")),
		reversaCom(
			rrset("50.2.0.192.in-addr.arpa.", "PTR", ptrOutra),
			rrset("20.2.0.192.in-addr.arpa.", "PTR", "a8.f.example."),
		))

	mudancas, err := records.UpdateReverso(t.Context(), httpc, "A", "192.0.2.50", "192.0.2.20", zonaDireta, "a8")
	if err != nil {
		t.Fatal(err)
	}

	if got := p.ptr("50.2.0.192.in-addr.arpa."); got != ptrOutra {
		t.Errorf("PTR de outra equipe = %q, esperado %q", got, ptrOutra)
	}

	if got := p.ptr("20.2.0.192.in-addr.arpa."); got != "" {
		t.Errorf("PTR do endereço antigo = %q, esperado removido", got)
	}

	if len(mudancas) != 1 || mudancas[0].Acao != records.AcaoReversoRemovido {
		t.Errorf("mudanças = %+v, esperado só a remoção do PTR antigo", mudancas)
	}
}

func TestDeleteReversoReapontaPTR(t *testing.T) {
	t.Parallel()

	p, httpc := novoPdns(t,
		diretaCom(rrset("b.f.example.", "A", "192.0.2.10")),
		reversaCom(rrset("10.2.0.192.in-addr.arpa.", "PTR", "a3.f.example.")))

	mudancas, err := records.DeleteReverso(t.Context(), httpc, "A", "192.0.2.10", zonaDireta, "a3")
	if err != nil {
		t.Fatal(err)
	}

	if got := p.ptr("10.2.0.192.in-addr.arpa."); got != "b.f.example." {
		t.Errorf("PTR = %q, esperado reapontado para b.f.example.", got)
	}

	if len(mudancas) != 1 || mudancas[0].Acao != records.AcaoReversoReapontado {
		t.Errorf("mudanças = %+v, esperado um reapontamento", mudancas)
	}
}

func TestDeleteReversoRemoveUltimaReferencia(t *testing.T) {
	t.Parallel()

	p, httpc := novoPdns(t,
		diretaCom(),
		reversaCom(rrset("10.2.0.192.in-addr.arpa.", "PTR", "a3.f.example.")))

	_, err := records.DeleteReverso(t.Context(), httpc, "A", "192.0.2.10", zonaDireta, "a3")
	if err != nil {
		t.Fatal(err)
	}

	if got := p.ptr("10.2.0.192.in-addr.arpa."); got != "" {
		t.Errorf("PTR = %q, esperado removido", got)
	}
}

func TestDeleteReversoPreservaPTRDeOutraEquipe(t *testing.T) {
	t.Parallel()

	p, httpc := novoPdns(t,
		diretaCom(),
		reversaCom(rrset("50.2.0.192.in-addr.arpa.", "PTR", ptrOutra)))

	mudancas, err := records.DeleteReverso(t.Context(), httpc, "A", "192.0.2.50", zonaDireta, "a4")
	if err != nil {
		t.Fatal(err)
	}

	if got := p.ptr("50.2.0.192.in-addr.arpa."); got != ptrOutra {
		t.Errorf("PTR de outra equipe = %q, esperado %q", got, ptrOutra)
	}

	if len(mudancas) != 0 {
		t.Errorf("mudanças = %+v, esperado nenhuma", mudancas)
	}
}

func TestDeleteReversoNaoRemoveComZonaIlegivel(t *testing.T) {
	t.Parallel()

	p, httpc := novoPdns(t,
		diretaCom(),
		records.Zone{Name: "g.example."},
		reversaCom(rrset("10.2.0.192.in-addr.arpa.", "PTR", "a3.f.example.")))
	p.falhar["g.example."] = true

	_, err := records.DeleteReverso(t.Context(), httpc, "A", "192.0.2.10", zonaDireta, "a3")
	if err == nil {
		t.Error("esperado erro com zona direta ilegível")
	}

	if got := p.ptr("10.2.0.192.in-addr.arpa."); got != "a3.f.example." {
		t.Errorf("PTR = %q, esperado intacto", got)
	}
}

func TestRecusaPdns(t *testing.T) {
	t.Parallel()

	p, httpc := novoPdns(t, diretaCom())
	p.recusar = true

	err := records.Patch(t.Context(), httpc, zonaDireta, records.PDNSZonePatchRequest{})

	motivo, ok := records.RecusaPdns(err)
	if !ok || motivo != "Name is out of zone" {
		t.Errorf("RecusaPdns = (%q, %v), esperado (\"Name is out of zone\", true)", motivo, ok)
	}
}
