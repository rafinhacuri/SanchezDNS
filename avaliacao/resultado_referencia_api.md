# Condição de referência — API do PowerDNS

Execução: 2026-10-15 16:04 -0300 · PowerDNS 5.0.1 · zona descartável `ref-sanchezdns.example.` · critérios em `avaliacao/plano.md`, seção 4.

| Caso | Entrada | Enviado (tipo nome ttl conteúdo) | HTTP | Mensagem do PowerDNS | Gravado | Publicado (dig) | Classificação | Corresponde ao pretendido |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| N1 | CNAME sem ponto final | CNAME n1.ref-sanchezdns.example. 300 www.destino.com.br | 422 | Record n1.ref-sanchezdns.example./CNAME 'www.destino.com.br': Not in expected format (parsed as 'www.destino.com.br.') | (nada) | - | recusada pelo servidor | - |
| N2 | NS sem ponto final | NS n2.ref-sanchezdns.example. 300 ns1.exemplo.com.br | 422 | Record n2.ref-sanchezdns.example./NS 'ns1.exemplo.com.br': Not in expected format (parsed as 'ns1.exemplo.com.br.') | (nada) | - | recusada pelo servidor | - |
| N3a | MX sem ponto final | MX n3a.ref-sanchezdns.example. 300 10 mail.exemplo.com.br | 422 | Record n3a.ref-sanchezdns.example./MX '10 mail.exemplo.com.br': Not in expected format (parsed as '10 mail.exemplo.com.br.') | (nada) | - | recusada pelo servidor | - |
| N3b | MX sem prioridade | MX n3b.ref-sanchezdns.example. 300 mail.exemplo.com.br. | 422 | Record n3b.ref-sanchezdns.example./MX 'mail.exemplo.com.br.': Parsing record content (try 'pdnsutil check-zone'): expected digits at position 0 in 'mail.exemplo.com.br.' | (nada) | - | recusada pelo servidor | - |
| N4 | TXT sem aspas | TXT n4.ref-sanchezdns.example. 300 v=spf1 -all | 422 | Record n4.ref-sanchezdns.example./TXT 'v=spf1 -all': Parsing record content (try 'pdnsutil check-zone'): Data field in DNS should start with quote (") at position 0 of 'v=spf1 -all' | (nada) | - | recusada pelo servidor | - |
| N5 | SRV sem ponto final no alvo | SRV n5.ref-sanchezdns.example. 300 10 5 5060 sip.exemplo.com.br | 422 | Record n5.ref-sanchezdns.example./SRV '10 5 5060 sip.exemplo.com.br': Not in expected format (parsed as '10 5 5060 sip.exemplo.com.br.') | (nada) | - | recusada pelo servidor | - |
| N6 | HTTPS sem destino | HTTPS n6.ref-sanchezdns.example. 300 1 | 422 | Record n6.ref-sanchezdns.example./HTTPS '1': Parsing record content (try 'pdnsutil check-zone'): missing field at the end of record content '1' | (nada) | - | recusada pelo servidor | - |
| N7 | CAA sem flag e tag | CAA n7.ref-sanchezdns.example. 300 letsencrypt.org | 422 | Record n7.ref-sanchezdns.example./CAA 'letsencrypt.org': Parsing record content (try 'pdnsutil check-zone'): expected digits at position 0 in 'letsencrypt.org' | (nada) | - | recusada pelo servidor | - |
| N8 | nome relativo | A www 300 192.0.2.8 | 422 | DNS Name 'www' is not canonical | (nada) | - | recusada pelo servidor | - |
| N9 | nome completo sem ponto final | A www.ref-sanchezdns.example 300 192.0.2.9 | 422 | DNS Name 'www.ref-sanchezdns.example' is not canonical | (nada) | - | recusada pelo servidor | - |
| N10 | nome já qualificado (controle) | A n10.ref-sanchezdns.example. 300 192.0.2.10 | 204 | - | n10.ref-sanchezdns.example. 300 192.0.2.10 | 192.0.2.10 | aceita como enviada | sim |
| N11 | endereço IPv4 inválido | A n11.ref-sanchezdns.example. 300 999.1.1.1 | 422 | Record n11.ref-sanchezdns.example./A '999.1.1.1': Parsing record content (try 'pdnsutil check-zone'): unable to parse IP address | (nada) | - | recusada pelo servidor | - |
| N12 | endereço IPv4 em AAAA | AAAA n12.ref-sanchezdns.example. 300 192.0.2.1 | 422 | Record n12.ref-sanchezdns.example./AAAA '192.0.2.1': Invalid IPv6 address | (nada) | - | recusada pelo servidor | - |
| N13 | TTL de 30 segundos | A n13.ref-sanchezdns.example. 30 192.0.2.13 | 204 | - | n13.ref-sanchezdns.example. 30 192.0.2.13 | 192.0.2.13 | aceita como enviada | não (deveria ser recusada) |
