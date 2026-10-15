#!/usr/bin/env bash
# =============================================================================
# Condição de referência do E2 (avaliacao/plano.md, seção 4) pelo pdnsutil,
# SEM o SanchezDNS. Roda NO SERVIDOR do PowerDNS, com permissão de usar o
# pdnsutil (em geral, root ou o usuário do pdns).
#
# Uso, no servidor:
#   sudo ./e2_referencia_pdnsutil.sh
#
# Copie o arquivo de resultado (resultado_referencia_pdnsutil.md) de volta para
# avaliacao/. A zona descartável ref-pdnsutil.example. (RFC 2606) é criada e
# removida. Usa os nomes de comando clássicos do pdnsutil (create-zone,
# add-record, list-zone, delete-zone), mantidos como apelidos no PowerDNS 5.
# =============================================================================

set -o pipefail

Z="ref-pdnsutil.example"
PDNSUTIL="${PDNSUTIL:-pdnsutil}"

DIR="$(cd "$(dirname "$0")" && pwd)"
OUT_MD="$DIR/resultado_referencia_pdnsutil.md"
OUT_LOG="$DIR/referencia_pdnsutil_execucao.log"

# O pdnsutil 5 trata o NOME como absoluto ("n1" vira "n1."), por isso os casos de
# conteúdo usam o nome completo; N8 a N10 exercitam justamente a forma do nome.
# id|tipo|nome (como digitado)|ttl|conteúdo|entrada|pretendido (nome conteúdo; "-" = deveria ser recusada)
CASOS=(
  "N1|CNAME|n1.$Z|300|www.destino.com.br|CNAME sem ponto final|n1.$Z. www.destino.com.br."
  "N2|NS|n2.$Z|300|ns1.exemplo.com.br|NS sem ponto final|n2.$Z. ns1.exemplo.com.br."
  "N3a|MX|n3a.$Z|300|10 mail.exemplo.com.br|MX sem ponto final|n3a.$Z. 10 mail.exemplo.com.br."
  "N3b|MX|n3b.$Z|300|mail.exemplo.com.br.|MX sem prioridade|-"
  "N4|TXT|n4.$Z|300|v=spf1 -all|TXT sem aspas|n4.$Z. \"v=spf1 -all\""
  "N5|SRV|n5.$Z|300|10 5 5060 sip.exemplo.com.br|SRV sem ponto final no alvo|n5.$Z. 10 5 5060 sip.exemplo.com.br."
  "N6|HTTPS|n6.$Z|300|1|HTTPS sem destino|-"
  "N7|CAA|n7.$Z|300|letsencrypt.org|CAA sem flag e tag|n7.$Z. 0 issue \"letsencrypt.org\""
  "N8|A|www8|300|192.0.2.8|nome relativo|www8.$Z. 192.0.2.8"
  "N9|A|www9.$Z|300|192.0.2.9|nome completo sem ponto final|www9.$Z. 192.0.2.9"
  "N10|A|www10.$Z.|300|192.0.2.10|nome já qualificado (controle)|www10.$Z. 192.0.2.10"
  "N11|A|n11.$Z|300|999.1.1.1|endereço IPv4 inválido|-"
  "N12|AAAA|n12.$Z|300|192.0.2.1|endereço IPv4 em AAAA|-"
  "N13|A|n13.$Z|30|192.0.2.13|TTL de 30 segundos|-"
)

log() { printf '%s\n' "$*" >>"$OUT_LOG"; }
info() { printf '%s\n' "$*"; log "$*"; }
die() { info "ERRO: $*"; exit 1; }

ZONA_CRIADA=0
limpeza() {
  local rc=$?
  if [[ $ZONA_CRIADA == 1 ]]; then
    "$PDNSUTIL" delete-zone "$Z" >>"$OUT_LOG" 2>&1
    info "  zona $Z removida: rc=$?"
  fi
  exit "$rc"
}
trap limpeza EXIT
trap 'exit 130' INT TERM

: >"$OUT_LOG"
command -v "$PDNSUTIL" >/dev/null || die "pdnsutil não encontrado"
VERSAO_PDNS=$("$PDNSUTIL" --version 2>&1 | head -1)
info "== Referência pelo pdnsutil · $VERSAO_PDNS · zona $Z"

"$PDNSUTIL" list-zone "$Z" >/dev/null 2>&1 && die "a zona $Z já existe; abortando"
"$PDNSUTIL" create-zone "$Z" >>"$OUT_LOG" 2>&1 || die "create-zone falhou (ver $OUT_LOG)"
ZONA_CRIADA=1

LINHAS=()
for c in "${CASOS[@]}"; do
  IFS='|' read -r id tipo nome ttl conteudo descr pretendido <<<"$c"
  saida=$("$PDNSUTIL" add-record "$Z" "$nome" "$tipo" "$ttl" "$conteudo" 2>&1)
  rc=$?
  log "CASO $id rc=$rc: $saida"
  saida=$(tr '\n' ' ' <<<"$saida" | sed 's/  */ /g; s/ *$//')

  gravado=$("$PDNSUTIL" list-zone "$Z" 2>/dev/null | awk -v t="$tipo" '$4==t' |
    awk '{n=$1; sub(/\.$/, "", n); c=""; for (i=5;i<=NF;i++) c=c (i>5?" ":"") $i; print n". "c}' | head -1)

  if [[ $rc -ne 0 || -z "$gravado" ]]; then
    classe="recusada pelo servidor"
    gravado="(nada)"
    confere="-"
  else
    g_nome="${gravado%% *}"
    if [[ "$gravado" == "${nome%.}. $conteudo" ]]; then
      classe="aceita como enviada"
    else
      classe="aceita e reescrita"
    fi
    if [[ "$pretendido" == "-" ]]; then confere="não (deveria ser recusada)"
    elif [[ "$gravado" == "$pretendido" ]]; then confere="sim"
    else confere="não"; fi
    "$PDNSUTIL" delete-rrset "$Z" "$g_nome" "$tipo" >>"$OUT_LOG" 2>&1
  fi

  info "  [$classe] $id $descr → rc=$rc"
  LINHAS+=("$id"$'\t'"$descr"$'\t'"add-record $Z $nome $tipo $ttl \"$conteudo\""$'\t'"$rc"$'\t'"${saida:--}"$'\t'"$gravado"$'\t'"$classe"$'\t'"$confere")
done

md() { local s="${1//|/\\|}"; echo "| ${s//$'\t'/ | } |"; }
{
  echo "# Condição de referência — pdnsutil"
  echo
  echo "Execução: $(date '+%Y-%m-%d %H:%M %z') · $VERSAO_PDNS · zona descartável \`$Z.\` · critérios em \`avaliacao/plano.md\`, seção 4."
  echo
  echo "| Caso | Entrada | Comando | Saída (rc) | Mensagem | Gravado (nome conteúdo) | Classificação | Corresponde ao pretendido |"
  echo "| --- | --- | --- | --- | --- | --- | --- | --- |"
  for l in "${LINHAS[@]}"; do md "$l"; done
} >"$OUT_MD"

info ""
info "  tabela: $OUT_MD"
