# Condição de referência — pdnsutil

Execução: 2026-10-15 16:04 -0300 · pdnsutil 5.0.1 · zona descartável `ref-pdnsutil.example.` · critérios em `avaliacao/plano.md`, seção 4.

| Caso | Entrada | Comando | Saída (rc) | Mensagem | Gravado (nome conteúdo) | Classificação | Corresponde ao pretendido |
| --- | --- | --- | --- | --- | --- | --- | --- |
| N1 | CNAME sem ponto final | add-record ref-pdnsutil.example n1.ref-pdnsutil.example CNAME 300 "www.destino.com.br" | 0 | New rrset: n1.ref-pdnsutil.example. 300 IN CNAME www.destino.com.br | n1.ref-pdnsutil.example. www.destino.com.br. | aceita e reescrita | sim |
| N2 | NS sem ponto final | add-record ref-pdnsutil.example n2.ref-pdnsutil.example NS 300 "ns1.exemplo.com.br" | 0 | New rrset: n2.ref-pdnsutil.example. 300 IN NS ns1.exemplo.com.br | n2.ref-pdnsutil.example. ns1.exemplo.com.br. | aceita e reescrita | sim |
| N3a | MX sem ponto final | add-record ref-pdnsutil.example n3a.ref-pdnsutil.example MX 300 "10 mail.exemplo.com.br" | 0 | New rrset: n3a.ref-pdnsutil.example. 300 IN MX 10 mail.exemplo.com.br | n3a.ref-pdnsutil.example. 10 mail.exemplo.com.br. | aceita e reescrita | sim |
| N3b | MX sem prioridade | add-record ref-pdnsutil.example n3b.ref-pdnsutil.example MX 300 "mail.exemplo.com.br." | 1 | Error: Parsing record content (try 'pdnsutil check-zone'): expected digits at position 0 in 'mail.exemplo.com.br.' | (nada) | recusada pelo servidor | - |
| N4 | TXT sem aspas | add-record ref-pdnsutil.example n4.ref-pdnsutil.example TXT 300 "v=spf1 -all" | 1 | Error: Parsing record content (try 'pdnsutil check-zone'): Data field in DNS should start with quote (") at position 0 of 'v=spf1 -all' | (nada) | recusada pelo servidor | - |
| N5 | SRV sem ponto final no alvo | add-record ref-pdnsutil.example n5.ref-pdnsutil.example SRV 300 "10 5 5060 sip.exemplo.com.br" | 0 | New rrset: n5.ref-pdnsutil.example. 300 IN SRV 10 5 5060 sip.exemplo.com.br | n5.ref-pdnsutil.example. 10 5 5060 sip.exemplo.com.br. | aceita e reescrita | sim |
| N6 | HTTPS sem destino | add-record ref-pdnsutil.example n6.ref-pdnsutil.example HTTPS 300 "1" | 1 | Error: Parsing record content (try 'pdnsutil check-zone'): missing field at the end of record content '1' | (nada) | recusada pelo servidor | - |
| N7 | CAA sem flag e tag | add-record ref-pdnsutil.example n7.ref-pdnsutil.example CAA 300 "letsencrypt.org" | 1 | Error: Parsing record content (try 'pdnsutil check-zone'): expected digits at position 0 in 'letsencrypt.org' | (nada) | recusada pelo servidor | - |
| N8 | nome relativo | add-record ref-pdnsutil.example www8 A 300 "192.0.2.8" | 1 | Error: Name "www8." to add is not part of zone "ref-pdnsutil.example.". | (nada) | recusada pelo servidor | - |
| N9 | nome completo sem ponto final | add-record ref-pdnsutil.example www9.ref-pdnsutil.example A 300 "192.0.2.9" | 0 | New rrset: www9.ref-pdnsutil.example. 300 IN A 192.0.2.9 | www9.ref-pdnsutil.example. 192.0.2.9 | aceita como enviada | sim |
| N10 | nome já qualificado (controle) | add-record ref-pdnsutil.example www10.ref-pdnsutil.example. A 300 "192.0.2.10" | 0 | New rrset: www10.ref-pdnsutil.example. 300 IN A 192.0.2.10 | www10.ref-pdnsutil.example. 192.0.2.10 | aceita como enviada | sim |
| N11 | endereço IPv4 inválido | add-record ref-pdnsutil.example n11.ref-pdnsutil.example A 300 "999.1.1.1" | 1 | Error: Parsing record content (try 'pdnsutil check-zone'): unable to parse IP address | (nada) | recusada pelo servidor | - |
| N12 | endereço IPv4 em AAAA | add-record ref-pdnsutil.example n12.ref-pdnsutil.example AAAA 300 "192.0.2.1" | 1 | Error: Parsing record content (try 'pdnsutil check-zone'): while parsing IPv6 address: '192.0.2.1' is invalid | (nada) | recusada pelo servidor | - |
| N13 | TTL de 30 segundos | add-record ref-pdnsutil.example n13.ref-pdnsutil.example A 30 "192.0.2.13" | 0 | New rrset: n13.ref-pdnsutil.example. 30 IN A 192.0.2.13 | n13.ref-pdnsutil.example. 192.0.2.13 | aceita como enviada | não (deveria ser recusada) |
