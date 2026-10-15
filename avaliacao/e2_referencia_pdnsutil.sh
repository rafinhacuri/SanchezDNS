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

# id|tipo|nome (como digitado)|ttl|conteúdo|entrada como o operador a escreveria
CASOS=(
  "N1|CNAME|n1|300|www.destino.com.br|CNAME sem ponto final"
  "N2|NS|n2|300|ns1.exemplo.com.br|NS sem ponto final"
  "N3a|MX|n3a|300|10 mail.exemplo.com.br|MX sem ponto final"
  "N3b|MX|n3b|300|mail.exemplo.com.br.|MX sem prioridade"
  "N4|TXT|n4|300|v=spf1 -all|TXT sem aspas"
  "N5|SRV|n5|300|10 5 5060 sip.exemplo.com.br|SRV sem ponto final no alvo"
  "N6|HTTPS|n6|300|1|HTTPS sem destino"
  "N7|CAA|n7|300|letsencrypt.org|CAA sem flag e tag"
  "N8|A|www8|300|192.0.2.8|nome relativo"
  "N9|A|www9.$Z|300|192.0.2.9|nome completo sem ponto final"
  "N10|A|www10.$Z.|300|192.0.2.10|nome já qualificado"
  "N11|A|n11|300|999.1.1.1|endereço IPv4 inválido"
  "N12|AAAA|n12|300|192.0.2.1|endereço IPv4 em AAAA"
  "N13|A|n13|30|192.0.2.13|TTL de 30 segundos"
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
  IFS='|' read -r id tipo nome ttl conteudo descr <<<"$c"
  saida=$("$PDNSUTIL" add-record "$Z" "$nome" "$tipo" "$ttl" "$conteudo" 2>&1)
  rc=$?
  log "CASO $id rc=$rc: $saida"
  saida=$(tr '\n' ' ' <<<"$saida" | sed 's/  */ /g; s/ *$//')

  gravado=$("$PDNSUTIL" list-zone "$Z" 2>/dev/null | awk -v t="$tipo" '$4==t' | grep -v 'SOA' |
    awk '{c=""; for (i=5;i<=NF;i++) c=c (i>5?" ":"") $i; print $1" "$2" "c}' | paste -sd';' -)

  if [[ $rc -ne 0 || -z "$gravado" ]]; then
    classe="recusada pelo servidor"
    gravado="(nada)"
  else
    g_cont=$(cut -d' ' -f3- <<<"${gravado%%;*}")
    g_nome=$(cut -d' ' -f1 <<<"${gravado%%;*}")
    esperado="$nome.$Z."; [[ "$nome" == *. ]] && esperado="$nome"
    if [[ "$g_cont" == "$conteudo" && "$g_nome" == "$esperado" ]]; then
      classe="aceita como enviada"
    else
      classe="aceita e reescrita"
    fi
    # Remove o RRset gravado para o próximo caso começar da zona vazia.
    "$PDNSUTIL" delete-rrset "$Z" "$nome" "$tipo" >>"$OUT_LOG" 2>&1 ||
      "$PDNSUTIL" delete-rrset "$Z" "${g_nome%.$Z.}" "$tipo" >>"$OUT_LOG" 2>&1
  fi

  info "  [$classe] $id $descr → rc=$rc"
  LINHAS+=("$id"$'\t'"$descr"$'\t'"add-record $Z $nome $tipo $ttl \"$conteudo\""$'\t'"$rc"$'\t'"${saida:--}"$'\t'"$gravado"$'\t'"$classe")
done

md() { local s="${1//|/\\|}"; echo "| ${s//$'\t'/ | } |"; }
{
  echo "# Condição de referência — pdnsutil"
  echo
  echo "Execução: $(date '+%Y-%m-%d %H:%M %z') · $VERSAO_PDNS · zona descartável \`$Z.\` · critérios em \`avaliacao/plano.md\`, seção 4."
  echo
  echo "| Caso | Entrada | Comando | Saída (rc) | Mensagem | Gravado (nome ttl conteúdo) | Classificação |"
  echo "| --- | --- | --- | --- | --- | --- | --- |"
  for l in "${LINHAS[@]}"; do md "$l"; done
} >"$OUT_MD"

info ""
info "  tabela: $OUT_MD"
