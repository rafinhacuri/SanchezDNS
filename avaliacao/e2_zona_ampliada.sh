#!/usr/bin/env bash
# =============================================================================
# R4 do plano (avaliacao/plano.md, seção 6): acrescenta à zona já delegada
# registros de tipos estruturados criados PELA PLATAFORMA, com as entradas que
# o operador digitaria (sem ponto final, sem aspas, campos separados), e
# confere com o dig o que foi publicado. Depois, rode o Zonemaster na zona e,
# se disponível, um validador independente sobre a zona exportada.
#
# Uso:
#   BASE_URL=https://<instância>/go VERSAO=3.0.3 ADMIN_EMAIL=... ADMIN_SENHA=... ./avaliacao/e2_zona_ampliada.sh adicionar
#   ... rode o Zonemaster ...
#   BASE_URL=... ADMIN_EMAIL=... ADMIN_SENHA=... ./avaliacao/e2_zona_ampliada.sh remover
# =============================================================================

set -o pipefail

ACAO="${1:-}"
BASE_URL="${BASE_URL:-}"
ADMIN_EMAIL="${ADMIN_EMAIL:-}"
ADMIN_SENHA="${ADMIN_SENHA:-}"
DNS_HOST="${DNS_HOST:-152.84.126.7}"
ZONA="${ZONA:-teste.curi.dev.br.}"
IP_MAIL="${IP_MAIL:-152.84.126.7}"
VERSAO="${VERSAO:-$(jq -r .version "$(dirname "$0")/../package.json" 2>/dev/null)}"

DIR="$(cd "$(dirname "$0")" && pwd)"
OUT_MD="$DIR/resultado_zona_ampliada_v$VERSAO.md"
TMP="$(mktemp -d)"
JAR="$TMP/admin.jar"
Z="${ZONA%.}"

# id|tipo|nome digitado|json com os campos do formulário|valor canônico esperado no dig
# O HTTPS é consultado como TYPE65 porque o dig do macOS (9.10) não conhece o tipo;
# o esperado é a forma genérica da RFC 3597 de "1 . alpn=h2" (prioridade 1, alvo
# raiz, chave 1 = alpn, valor "h2").
REGISTROS=(
  "Z1|A|mail|{\"vl\":\"$IP_MAIL\"}|$IP_MAIL"
  "Z2|MX|$Z|{\"vl\":\"mail.$Z\",\"priority\":10}|10 mail.$Z."
  "Z3|SRV|_sip._tcp|{\"priority\":10,\"weight\":5,\"port\":5060,\"target\":\"mail.$Z\"}|10 5 5060 mail.$Z."
  "Z4|CAA|$Z|{\"vl\":\"letsencrypt.org\"}|0 issue \"letsencrypt.org\""
  "Z5|TYPE65|$Z|{\"svcPriority\":1,\"targetName\":\"\",\"svcParams\":\"\"}|\\# 10 00010000010003026832"
  "Z6|CNAME|www|{\"vl\":\"mail.$Z\"}|mail.$Z."
)

die() { echo "ERRO: $*"; exit 1; }

req() { # req <método> <caminho> [json] → ST e BODY
  local -a args=(-sS -m 40 -X "$1" -o "$TMP/body" -w '%{http_code}' -b "$JAR" -c "$JAR" -H 'Accept: application/json')
  [[ -n "${3-}" ]] && args+=(-H 'Content-Type: application/json' --data-binary @-)
  ST=$(printf '%s' "${3-}" | curl "${args[@]}" "$BASE_URL$2" 2>/dev/null) || ST="000"
  BODY=$(cat "$TMP/body" 2>/dev/null)
}

tipo_api() { [[ "$1" == TYPE65 ]] && echo HTTPS || echo "$1"; }

corpo() { # corpo <tipo> <nome> <json dos campos>
  jq -cn --arg z "$ZONA" --arg t "$(tipo_api "$1")" --arg n "$2" --argjson f "$3" '{zone:$z, type:$t, name:$n, ttl:3600} + $f'
}

# A remoção (DELETE /records) espera o MX como a tela o envia: prioridade dentro do
# valor E no campo próprio. Sem o campo, a API responde 200 sem remover nada.
corpo_remocao() {
  if [[ "$1" == MX ]]; then
    jq -cn --arg z "$ZONA" --arg n "$2" --argjson f "$3" '{zone:$z, type:"MX", name:$n, ttl:3600, vl:"\($f.priority) \($f.vl)", priority:$f.priority}'
  else
    corpo "$@"
  fi
}

nome_q() { [[ "$1" == "$Z" ]] && echo "$Z." || echo "$1.$Z."; }

trap 'req POST /logout >/dev/null 2>&1; rm -rf "$TMP"' EXIT

[[ "$ACAO" == adicionar || "$ACAO" == remover ]] || die "uso: $0 adicionar|remover"
[[ -n "$BASE_URL" && -n "$ADMIN_EMAIL" && -n "$ADMIN_SENHA" ]] || die "defina BASE_URL, ADMIN_EMAIL e ADMIN_SENHA no ambiente"
for b in curl jq dig; do command -v "$b" >/dev/null || die "falta o utilitário $b"; done

: >"$JAR"
req POST /login "$(jq -cn --arg e "$ADMIN_EMAIL" --arg s "$ADMIN_SENHA" '{email:$e, senha:$s}')"
[[ "$ST" == 200 ]] || die "login falhou: $ST"

if [[ "$ACAO" == remover ]]; then
  for ((i = ${#REGISTROS[@]} - 1; i >= 0; i--)); do
    IFS='|' read -r id tipo nome campos _ <<<"${REGISTROS[$i]}"
    req DELETE /records "$(corpo_remocao "$tipo" "$nome" "$campos")"
    echo "  $id $tipo $(nome_q "$nome") removido: $ST"
  done
  exit 0
fi

LINHAS=()
for r in "${REGISTROS[@]}"; do
  IFS='|' read -r id tipo nome campos esperado <<<"$r"
  req PUT /records "$(corpo "$tipo" "$nome" "$campos")"
  st="$ST"
  sleep 2
  publicado=$(dig +short +time=2 +tries=2 @"$DNS_HOST" "$(nome_q "$nome")" "$tipo" 2>/dev/null | paste -sd, -)
  if [[ "$st" == 200 && "${publicado// /}" == *"${esperado// /}"* ]]; then res="publicado na forma canônica"
  elif [[ "$st" == 200 ]]; then res="publicado com valor diferente do esperado"
  else res="recusado ($st)"; fi
  echo "  [$res] $id $tipo $(nome_q "$nome") → ${publicado:-(vazio)}"
  LINHAS+=("$id | $tipo | $(nome_q "$nome") | \`$(jq -c . <<<"$campos")\` | $st | \`$esperado\` | \`${publicado:-(vazio)}\` | $res")
done

{
  echo "# Zona ampliada para os oráculos externos (R4)"
  echo
  echo "Execução: $(date '+%Y-%m-%d %H:%M %z') · versão $VERSAO · zona \`$ZONA\` · critérios em \`avaliacao/plano.md\`, seção 6."
  echo
  echo "| Id | Tipo | Nome | Campos enviados | HTTP | Esperado | Publicado (dig) | Resultado |"
  echo "| --- | --- | --- | --- | --- | --- | --- | --- |"
  for l in "${LINHAS[@]}"; do echo "| $l |"; done
  echo
  echo "## Oráculos externos"
  echo
  echo "- Zonemaster (https://zonemaster.net), executado em: _preencher_ · resultado: _preencher_"
  echo "- Validador independente (named-checkzone ou kzonecheck) sobre \`pdnsutil list-zone $Z\`: _preencher ou registrar que não foi executado_"
} >"$OUT_MD"

echo
echo "  tabela: $OUT_MD"
echo "  agora rode o Zonemaster em $Z e depois: $0 remover"
