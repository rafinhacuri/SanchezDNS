# Zona ampliada para os oráculos externos (R4)

Execução: 2026-10-15 16:08 -0300 · versão 3.0.3 · zona `teste.curi.dev.br.` · critérios em `avaliacao/plano.md`, seção 6.

| Id | Tipo | Nome | Campos enviados | HTTP | Esperado | Publicado (dig) | Resultado |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Z1 | A | mail.teste.curi.dev.br. | `{"vl":"152.84.126.7"}` | 200 | `152.84.126.7` | `152.84.126.7` | publicado na forma canônica |
| Z2 | MX | teste.curi.dev.br. | `{"vl":"mail.teste.curi.dev.br","priority":10}` | 200 | `10 mail.teste.curi.dev.br.` | `10 mail.teste.curi.dev.br.` | publicado na forma canônica |
| Z3 | SRV | _sip._tcp.teste.curi.dev.br. | `{"priority":10,"weight":5,"port":5060,"target":"mail.teste.curi.dev.br"}` | 200 | `10 5 5060 mail.teste.curi.dev.br.` | `10 5 5060 mail.teste.curi.dev.br.` | publicado na forma canônica |
| Z4 | CAA | teste.curi.dev.br. | `{"vl":"letsencrypt.org"}` | 200 | `0 issue "letsencrypt.org"` | `0 issue "letsencrypt.org"` | publicado na forma canônica |
| Z5 | TYPE65 | teste.curi.dev.br. | `{"svcPriority":1,"targetName":"","svcParams":""}` | 200 | `\# 10 00010000010003026832` | `\# 10 00010000010003026832` | publicado na forma canônica |
| Z6 | CNAME | www.teste.curi.dev.br. | `{"vl":"mail.teste.curi.dev.br"}` | 200 | `mail.teste.curi.dev.br.` | `mail.teste.curi.dev.br.` | publicado na forma canônica |

## Oráculos externos

- Zonemaster (https://zonemaster.net), motor 9.0.1, executado em 15/10/2026 durante a primeira execução deste script, com Z1 a Z6 publicados (`zonemaster_teste.curi.dev.br_v3.0.3.json`). Resultado: 58 informações, 7 notas, 5 avisos e 4 erros, todos de origem já descrita na seção 5.3.2 (um único servidor de nomes, PTR do endereço do NS fora do servidor, DS ausente no pai, refresh e retry abaixo do recomendado). O MX passou no Syntax08 e nos testes Zone08 e Zone09. O Zonemaster não tem casos de teste para SRV, CAA ou HTTPS.
- Validador independente (named-checkzone ou kzonecheck): não executado; nenhum dos dois estava disponível. Para SRV, CAA e HTTPS, o oráculo é o `dig` contra o servidor autoritativo, com o HTTPS conferido na forma genérica da RFC 3597.
- Na primeira execução, o HTTPS não foi verificado (o `dig` 9.10 consultou o A) e, na segunda, o MX foi recusado como duplicado porque a remoção anterior não o tinha apagado; ver `plano.md`, seção 8. A tabela acima é a da terceira execução.
