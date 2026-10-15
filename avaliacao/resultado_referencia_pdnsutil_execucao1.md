# Condição de referência — pdnsutil

Execução: 2026-10-15 15:58 -0300 · pdnsutil 5.0.1 · zona descartável `ref-pdnsutil.example.` · critérios em `avaliacao/plano.md`, seção 4.

| Caso | Entrada | Comando | Saída (rc) | Mensagem | Gravado (nome ttl conteúdo) | Classificação |
| --- | --- | --- | --- | --- | --- | --- |
| N1 | CNAME sem ponto final | add-record ref-pdnsutil.example n1 CNAME 300 "www.destino.com.br" | 1 | Error: Name "n1." to add is not part of zone "ref-pdnsutil.example.". | (nada) | recusada pelo servidor |
| N2 | NS sem ponto final | add-record ref-pdnsutil.example n2 NS 300 "ns1.exemplo.com.br" | 1 | Error: Name "n2." to add is not part of zone "ref-pdnsutil.example.". | (nada) | recusada pelo servidor |
| N3a | MX sem ponto final | add-record ref-pdnsutil.example n3a MX 300 "10 mail.exemplo.com.br" | 1 | Error: Name "n3a." to add is not part of zone "ref-pdnsutil.example.". | (nada) | recusada pelo servidor |
| N3b | MX sem prioridade | add-record ref-pdnsutil.example n3b MX 300 "mail.exemplo.com.br." | 1 | Error: Name "n3b." to add is not part of zone "ref-pdnsutil.example.". | (nada) | recusada pelo servidor |
| N4 | TXT sem aspas | add-record ref-pdnsutil.example n4 TXT 300 "v=spf1 -all" | 1 | Error: Name "n4." to add is not part of zone "ref-pdnsutil.example.". | (nada) | recusada pelo servidor |
| N5 | SRV sem ponto final no alvo | add-record ref-pdnsutil.example n5 SRV 300 "10 5 5060 sip.exemplo.com.br" | 1 | Error: Name "n5." to add is not part of zone "ref-pdnsutil.example.". | (nada) | recusada pelo servidor |
| N6 | HTTPS sem destino | add-record ref-pdnsutil.example n6 HTTPS 300 "1" | 1 | Error: Name "n6." to add is not part of zone "ref-pdnsutil.example.". | (nada) | recusada pelo servidor |
| N7 | CAA sem flag e tag | add-record ref-pdnsutil.example n7 CAA 300 "letsencrypt.org" | 1 | Error: Name "n7." to add is not part of zone "ref-pdnsutil.example.". | (nada) | recusada pelo servidor |
| N8 | nome relativo | add-record ref-pdnsutil.example www8 A 300 "192.0.2.8" | 1 | Error: Name "www8." to add is not part of zone "ref-pdnsutil.example.". | (nada) | recusada pelo servidor |
| N9 | nome completo sem ponto final | add-record ref-pdnsutil.example www9.ref-pdnsutil.example A 300 "192.0.2.9" | 0 | New rrset: www9.ref-pdnsutil.example. 300 IN A 192.0.2.9 | www9.ref-pdnsutil.example 300 192.0.2.9 | aceita e reescrita |
| N10 | nome já qualificado | add-record ref-pdnsutil.example www10.ref-pdnsutil.example. A 300 "192.0.2.10" | 0 | New rrset: www10.ref-pdnsutil.example. 300 IN A 192.0.2.10 | www10.ref-pdnsutil.example 300 192.0.2.10 | aceita e reescrita |
| N11 | endereço IPv4 inválido | add-record ref-pdnsutil.example n11 A 300 "999.1.1.1" | 1 | Error: Name "n11." to add is not part of zone "ref-pdnsutil.example.". | (nada) | recusada pelo servidor |
| N12 | endereço IPv4 em AAAA | add-record ref-pdnsutil.example n12 AAAA 300 "192.0.2.1" | 1 | Error: Name "n12." to add is not part of zone "ref-pdnsutil.example.". | (nada) | recusada pelo servidor |
| N13 | TTL de 30 segundos | add-record ref-pdnsutil.example n13 A 30 "192.0.2.13" | 1 | Error: Name "n13." to add is not part of zone "ref-pdnsutil.example.". | (nada) | recusada pelo servidor |
