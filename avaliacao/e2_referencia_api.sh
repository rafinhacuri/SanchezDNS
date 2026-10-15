#!/usr/bin/env bash
# =============================================================================
# Condição de referência do E2 (avaliacao/plano.md, seção 4) pela API REST do
# PowerDNS, SEM o SanchezDNS. As entradas dos casos N1–N13 do Quadro 12 são
# enviadas como um operador as escreveria direto na API, e cada uma é
# classificada como recusada pelo servidor, aceita e reescrita, ou aceita como
# enviada.
#
# Uso:
#   PDNS_URL=http://152.84.126.7:8081 PDNS_API_KEY=... ./avaliacao/e2_referencia_api.sh
#
# A zona descartável ref-sanchezdns.example. (RFC 2606) é criada e removida.
# =============================================================================

set -o pipefail

PDNS_URL="${PDNS_URL:-}"
PDNS_API_KEY="${PDNS_API_KEY:-}"
PDNS_SERVER_ID="${PDNS_SERVER_ID:-localhost}"
DNS_HOST="${DNS_HOST:-152.84.126.7}"

Z="ref-sanchezdns.example."

DIR="$(cd "$(dirname "$0")" && pwd)"
OUT_MD="$DIR/resultado_referencia_api.md"
OUT_LOG="$DIR/referencia_api_execucao.log"
TMP="$(mktemp -d)"

# id|tipo|nome|ttl|conteúdo|entrada como o operador a escreveria
CASOS=(
  "N1|CNAME|n1.$Z|300|www.destino.com.br|CNAME sem ponto final"
  "N2|NS|n2.$Z|300|ns1.exemplo.com.br|NS sem ponto final"
  "N3a|MX|n3a.$Z|300|10 mail.exemplo.com.br|MX sem ponto final"
  "N3b|MX|n3b.$Z|300|mail.exemplo.com.br.|MX sem prioridade"
  "N4|TXT|n4.$Z|300|v=spf1 -all|TXT sem aspas"
  "N5|SRV|n5.$Z|300|10 5 5060 sip.exemplo.com.br|SRV sem ponto final no alvo"
  "N6|HTTPS|n6.$Z|300|1|HTTPS sem destino"
  "N7|CAA|n7.$Z|300|letsencrypt.org|CAA sem flag e tag"
  "N8|A|www|300|192.0.2.8|nome relativo"
  "N9|A|www.${Z%.}|300|192.0.2.9|nome completo sem ponto final"
  "N10|A|n10.$Z|300|192.0.2.10|nome já qualificado (controle)"
  "N11|A|n11.$Z|300|999.1.1.1|endereço IPv4 inválido"
  "N12|AAAA|n12.$Z|300|192.0.2.1|endereço IPv4 em AAAA"
  "N13|A|n13.$Z|30|192.0.2.13|TTL de 30 segundos"
)

log() { printf '%s\n' "$*" >>"$OUT_LOG"; }
info() { printf '%s\n' "$*"; log "$*"; }
die() { info "ERRO: $*"; exit 1; }

# api <método> <caminho> [json] → ST e BODY; a chave vai por stdin (-K -), fora da lista de processos.
api() {
  local metodo="$1" caminho="$2" dados="${3-}"
  local -a args=(-sS -m 30 -X "$metodo" -o "$TMP/body" -w '%{http_code}' -K "$TMP/cfg")
  [[ -n "$dados" ]] && args+=(-H 'Content-Type: application/json' --data-binary "$dados")
  ST=$(curl "${args[@]}" "$PDNS_URL/api/v1/servers/$PDNS_SERVER_ID$caminho" 2>>"$OUT_LOG") || ST="000"
  BODY=$(cat "$TMP/body" 2>/dev/null)
  log "  $metodo $caminho → $ST ${BODY:0:300}"
}

ZONA_CRIADA=0
limpeza() {
  local rc=$?
  if [[ $ZONA_CRIADA == 1 ]]; then
    api DELETE "/zones/$Z"
    info "  zona $Z removida: $ST"
  fi
  rm -rf "$TMP"
  exit "$rc"
}
trap limpeza EXIT
trap 'exit 130' INT TERM

: >"$OUT_LOG"
for b in curl jq dig; do command -v "$b" >/dev/null || die "falta o utilitário $b"; done
[[ -n "$PDNS_URL" && -n "$PDNS_API_KEY" ]] || die "defina PDNS_URL e PDNS_API_KEY no ambiente"
printf 'header = "X-API-Key: %s"\nheader = "Accept: application/json"\n' "$PDNS_API_KEY" >"$TMP/cfg"
chmod 600 "$TMP/cfg"

info "== Referência pela API do PowerDNS · $PDNS_URL · zona $Z"

api GET /zones
[[ "$ST" == 200 ]] || die "a API do PowerDNS respondeu $ST"
VERSAO_PDNS=$(curl -sS -m 10 -K "$TMP/cfg" "$PDNS_URL/api/v1/servers/$PDNS_SERVER_ID" 2>/dev/null | jq -r '.version // "?"')
jq -e --arg z "$Z" '.[] | select(.name==$z)' <<<"$BODY" >/dev/null && die "a zona $Z já existe; abortando"

api POST /zones "$(jq -cn --arg z "$Z" '{name:$z, kind:"Native"}')"
[[ "$ST" == 201 ]] || die "criação da zona falhou: $ST $BODY"
ZONA_CRIADA=1
info "  zona criada (PowerDNS $VERSAO_PDNS)"

LINHAS=()
for c in "${CASOS[@]}"; do
  IFS='|' read -r id tipo nome ttl conteudo descr <<<"$c"
  corpo=$(jq -cn --arg n "$nome" --arg t "$tipo" --argjson ttl "$ttl" --arg v "$conteudo" \
    '{rrsets:[{name:$n, type:$t, ttl:$ttl, changetype:"REPLACE", records:[{content:$v, disabled:false}]}]}')
  api PATCH "/zones/$Z" "$corpo"
  st="$ST"
  erro=$(jq -r '.error // empty' <<<"$BODY" 2>/dev/null)

  api GET "/zones/$Z"
  gravado=$(jq -r --arg t "$tipo" --arg id "$(tr 'A-Z' 'a-z' <<<"$id")" \
    '[.rrsets[]? | select(.type==$t and (.name | startswith($id + ".") or startswith("www")))
      | "\(.name) \(.ttl) \(.records | map(.content) | join(","))"] | join("; ")' <<<"$BODY")

  if [[ "$st" == 204 ]]; then
    gravado_conteudo=$(jq -r --arg t "$tipo" --arg id "$(tr 'A-Z' 'a-z' <<<"$id")" \
      '[.rrsets[]? | select(.type==$t and (.name | startswith($id + ".") or startswith("www"))) | .records[].content] | first // ""' <<<"$BODY")
    if [[ "$gravado_conteudo" == "$conteudo" && "$id" != N8 && "$id" != N9 ]]; then
      classe="aceita como enviada"
    else
      classe="aceita e reescrita"
    fi
  else
    classe="recusada pelo servidor"
    gravado="(nada)"
  fi

  publicado="-"
  if [[ "$st" == 204 ]]; then
    nome_q="$nome"; [[ "$nome_q" != *. ]] && nome_q="$nome_q."
    publicado=$(dig +short +time=2 +tries=2 @"$DNS_HOST" "$nome_q" "$tipo" 2>/dev/null | paste -sd, -)
    api PATCH "/zones/$Z" "$(jq -cn --arg n "$nome" --arg t "$tipo" '{rrsets:[{name:$n, type:$t, changetype:"DELETE"}]}')"
  fi

  info "  [$classe] $id $descr → HTTP $st ${erro:+($erro)}"
  LINHAS+=("$id"$'\t'"$descr"$'\t'"$tipo $nome $ttl $conteudo"$'\t'"$st"$'\t'"${erro:--}"$'\t'"${gravado:--}"$'\t'"${publicado:--}"$'\t'"$classe")
done

md() { local s="${1//|/\\|}"; echo "| ${s//$'\t'/ | } |"; }
{
  echo "# Condição de referência — API do PowerDNS"
  echo
  echo "Execução: $(date '+%Y-%m-%d %H:%M %z') · PowerDNS $VERSAO_PDNS · zona descartável \`$Z\` · critérios em \`avaliacao/plano.md\`, seção 4."
  echo
  echo "| Caso | Entrada | Enviado (tipo nome ttl conteúdo) | HTTP | Mensagem do PowerDNS | Gravado | Publicado (dig) | Classificação |"
  echo "| --- | --- | --- | --- | --- | --- | --- | --- |"
  for l in "${LINHAS[@]}"; do md "$l"; done
} >"$OUT_MD"

info ""
info "  tabela: $OUT_MD"
