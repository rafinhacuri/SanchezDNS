# 📊 Estatísticas

A página de estatísticas dá uma visão operacional da instância PowerDNS configurada no ambiente. Esta seção explica **o que é exibido, de onde cada número vem e como o backend o calcula**.

## Métricas exibidas

- **quantidade de zonas**
- **quantidade total de registros**
- **uptime** formatado (ex.: `3d 4h 12m`)
- **status** da instância
- **UDP Queries** e **TCP Queries**
- **Server ID**
- **momento de início** (calculado a partir do uptime)

## De onde vêm os dados

O backend combina **duas fontes da API do PowerDNS** ([statistics.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/controller/fetch/statistics.go)):

1. **`GET /api/v1/servers/{id}/statistics`** — a lista de métricas internas do servidor (uptime, `udp-queries`, `tcp-queries`, entre muitas outras).
2. **`GET /api/v1/servers/{id}/zones`** + um `GET` por zona — para **contar zonas e registros**.

### Por que percorrer todas as zonas para contar registros

O endpoint de estatísticas do PowerDNS **não entrega uma contagem pronta de registros por zona** no formato que o painel precisa. Então o backend lista as zonas e, para cada uma, busca seus detalhes e **soma o tamanho de cada RRset**:

```go
for _, rr := range zd.RRsets {
    records += len(rr.Records)
}
```

Isso é coerente com o modelo de **RRset** explicado em [PowerDNS a fundo](/powerdns#rrset-o-conceito-central-da-api): como cada RRset agrupa vários registros do mesmo nome+tipo, a contagem real de registros é a soma dos `records` de todos os RRsets. Se a busca de uma zona específica falhar, o backend simplesmente a ignora e continua, para não derrubar a página inteira por causa de uma zona problemática.

### Tratamento dos valores das métricas

As estatísticas do PowerDNS vêm com valores de tipos variados (string, número, booleano) dentro de um JSON. O backend faz um *parsing* tolerante: tenta interpretar cada valor como string, depois booleano, depois inteiro, depois float, guardando tudo normalizado em um mapa. Assim ele consegue extrair com segurança números como `uptime` e `udp-queries` sem quebrar se o formato de uma métrica variar entre versões do PowerDNS.

### Uptime e "momento de início"

O uptime vem em segundos e é formatado para leitura humana (`humanUptime`). O **momento de início** é derivado subtraindo o uptime do horário atual (`startedAtFromNow`) — ou seja, não é um dado bruto do servidor, é um cálculo feito no backend para exibir "desde quando o servidor está no ar".

## Atualização automática

A interface busca os dados no backend e se **atualiza sozinha a cada 60 segundos**, dando uma visão quase em tempo real sem exigir refresh manual.

## Limite atual

A tela reflete **uma única instância PowerDNS** por ambiente — coerente com o modelo do projeto (ver [Configuração](/configuration#o-modelo-uma-instancia-powerdns-por-ambiente)). Não há agregação de métricas entre vários servidores.
