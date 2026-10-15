#!/usr/bin/env bash
# =============================================================================
# Episódio E2 — testes de autorização contra a API (QP3) e cobertura da trilha
# de auditoria (QP4), executados de ponta a ponta contra uma instância real do
# SanchezDNS, com o PowerDNS consultado por `dig` como oráculo externo.
#
# Uso:
#   ADMIN_EMAIL=... ADMIN_SENHA=... ./avaliacao/e2_autorizacao.sh
#
# Rotas, payloads e códigos conferidos em api/routes/routes.go e nos
# controladores de api/controller/ (versão 3.0.2):
#   POST   /login                 {email, senha}            200 | 400 se já há cookie
#   GET    /session                                         200 {email, level}
#   POST   /logout                                          200
#   POST   /cadastro              {nome, email, senha, foto} 200 (cria solicitação)
#   GET    /solicitacoes?filter=  (admin)                   {solicitacoes[], total}
#   PUT    /solicitacoes/status   {id, status} (admin)      200 (status "aprovada" cria o cadastro)
#   GET    /cadastros?filter=     (admin)                   {cadastros[], total}
#   DELETE /cadastro              {id} (admin)              200 (também retira o e-mail das zonas)
#   POST   /user                  {email, permissao, zona} (admin)
#   PUT    /zone                  {domain, type, soa{...}} (admin)
#   DELETE /zone?id=<zona>        (admin)
#   GET    /records?zone=         leitura na zona           403 sem vínculo
#   PUT    /records               Registro                  403 sem escrita
#   PATCH  /records               {oldValue, newValue}      403 sem escrita
#   DELETE /records               Registro                  403 sem escrita
#   GET    /reverses/orphans?zone= leitura na zona          403 sem vínculo
#   DELETE /reverses?zone=&name=  escrita na zona           403 sem escrita
#   GET    /logs?filter=&limit=   (admin)                   {logs[], total}
#
# Distinção sessão × 403 (api/middleware/validate-session.go): o middleware de
# sessão responde 401 a partir da 3.0.3 (400 até a 3.0.2) com "Missing session
# ID" / "Invalid session", antes de qualquer verificação de permissão; 403 só
# sai do controlador, depois que a sessão foi aceita. Um erro de sessão num caso que espera 403 não diz nada
# sobre autorização, por isso é classificado como INCONCLUSIVO, e não como
# APROVADO ou REPROVADO.
#
# Classes do Quadro 1 exercitadas (QP entre parênteses):
#   A0        Atualização dinâmica não autenticada (QP3)
#   A1 A2     Privilégio excessivo do operador (QP3)
#   A3        Divergência entre A/AAAA e PTR (QP2), dentro das condições da §3.4.5 (QP3)
#   A4 A5     Privilégio excessivo pela escrita derivada — condições 2 e 1 da §3.4.5 (QP3)
#   A6 A7     Privilégio excessivo: o operador não lista, não edita e não triagem a zona reversa (QP3)
#   A8        Privilégio excessivo pela escrita derivada na alteração de endereço — §6.2.1 (QP3)
#   A9        Divergência A/PTR: remoção de PTR ainda referenciado — guarda de referência (QP2)
#   A10       Registros pendentes e reversos órfãos: remoção encadeada (QP2)
#   A11       Registros pendentes e reversos órfãos: triagem de órfãos (QP2)
#   A12       Privilégio excessivo: escape de escopo por nome plenamente qualificado (QP3)
#   Trilha    Ausência de rastreabilidade das alterações (QP4)
# =============================================================================

set -o pipefail

# ----------------------------------------------------------------------------
# CONFIGURAÇÃO
# ----------------------------------------------------------------------------
# O backend Go fica atrás do prefixo /go (Caddyfile do Dockerfile: handle_path /go/*).
BASE_URL="${BASE_URL:-https://mizle.cbpf.dev.br:3000/go}"
# Credenciais do administrador: só por variável de ambiente.
ADMIN_EMAIL="${ADMIN_EMAIL:-}"
ADMIN_SENHA="${ADMIN_SENHA:-}"
# PowerDNS autoritativo consultado pelo dig (oráculo externo ao artefato).
DNS_HOST="${DNS_HOST:-152.84.126.7}"

# Zonas descartáveis: domínio reservado (RFC 2606) e faixa de documentação (RFC 5737).
ZD="aval-sanchezdns.example."
ZR="2.0.192.in-addr.arpa."
# Faixa de documentação sem zona reversa provisionada (caso A5).
IP_SEM_ZONA="198.51.100.7"
# PTR "de outra equipe" semeado pelo admin.
IP_OUTRA="192.0.2.50"
PTR_OUTRA="servidor.outra-equipe.example."

TTL=300
# Segundos que o dig pode levar para refletir uma alteração feita pela API.
DIG_ESPERA="${DIG_ESPERA:-15}"
# Pausa antes de ler um estado que NÃO deveria mudar (evita ler antes da escrita).
PAUSA="${PAUSA:-2}"
# Atraso antes de ler a trilha: a gravação do log é assíncrona (go logs.Insert).
ESPERA_LOGS="${ESPERA_LOGS:-4}"

DIR="$(cd "$(dirname "$0")" && pwd)"
OUT_MD="$DIR/resultado_e2.md"
OUT_LOG="$DIR/e2_execucao.log"

RUN="$(date +%Y%m%d%H%M%S)"
TMP="$(mktemp -d)"

ADMIN_JAR="$TMP/admin.jar"
OP_ROT=(opA opB opC)
OP_PERM=(escrita leitura "")
OP_EMAIL=("aval-e2-opa-$RUN@example.com" "aval-e2-opb-$RUN@example.com" "aval-e2-opc-$RUN@example.com")
OP_JAR=("$TMP/opa.jar" "$TMP/opb.jar" "$TMP/opc.jar")
OP_SENHA=("" "" "")
OP_SOLIC=("" "" "")
OP_CAD=("" "" "")
OP_SOLICITADO=(0 0 0)
OP_LOGADO=(0 0 0)

ADMIN_LOGADO=0
ZD_TENTADA=0
ZR_TENTADA=0

CASOS=()      # "caso<TAB>esperado<TAB>obtido<TAB>resultado"
ESCRITAS=()   # "categoria|acao|autor|zona|chave|descrição"  (escritas com HTTP 200)
RECUSADAS=()  # "autor|chave|descrição"                       (tentativas negadas)
DERIVADAS=()  # "nome reverso|descrição"                       (PTR alterado pela automação)

# ----------------------------------------------------------------------------
# UTILITÁRIOS
# ----------------------------------------------------------------------------
log() { printf '%s\n' "$*" >>"$OUT_LOG"; }
info() { printf '%s\n' "$*"; log "$*"; }
die() { info "ERRO: $*"; exit 1; }

q() { jq -rn --arg v "$1" '$v|@uri'; }

senha_aleatoria() { openssl rand -hex 12 2>/dev/null || LC_ALL=C tr -dc 'a-z0-9' </dev/urandom | head -c 24; }

# req <jar|-> <método> <caminho> [json]  → ST (status) e BODY (resposta)
# O corpo vai por stdin para não expor senhas na lista de processos.
req() {
  local jar="$1" metodo="$2" caminho="$3" dados="${4-}"
  local -a args=(-sS -m 40 -X "$metodo" -o "$TMP/body" -w '%{http_code}' -H 'Accept: application/json')
  [[ "$jar" != "-" ]] && args+=(-b "$jar" -c "$jar")
  [[ -n "$dados" ]] && args+=(-H 'Content-Type: application/json' --data-binary @-)
  ST=$(printf '%s' "$dados" | curl "${args[@]}" "$BASE_URL$caminho" 2>>"$OUT_LOG") || ST="000"
  BODY=$(cat "$TMP/body" 2>/dev/null)
  : >"$TMP/body"
  if [[ "$caminho" == "/login" ]]; then
    log "  $metodo $caminho → $ST"
  else
    log "  $metodo $caminho → $ST ${BODY:0:400}"
  fi
}

msg() { jq -r '.message // .error // empty' <<<"$BODY" 2>/dev/null; }

# 401 (ou 400 até a 3.0.2) vindo do middleware de sessão, e não de validação de corpo.
erro_sessao() { [[ "$ST" == 401 || "$ST" == 400 ]] && [[ "$(msg)" == *session* ]]; }

# Descrição do status para a coluna "Obtido".
desc() {
  case "$ST" in
    400|401) if erro_sessao; then echo "$ST (sessão: $(msg))"; else echo "$ST ($(msg))"; fi ;;
    000) echo "000 (sem resposta)" ;;
    *) echo "$ST" ;;
  esac
}

# Veredito sobre o status: esperado → APROVADO; erro de sessão, sem resposta
# ou 5xx → INCONCLUSIVO (não mede autorização); qualquer outro → REPROVADO.
v_status() {
  if [[ "$ST" == "$1" ]]; then echo APROVADO
  elif erro_sessao || [[ "$ST" == 000 || "$ST" == 5* ]]; then echo INCONCLUSIVO
  else echo REPROVADO; fi
}

pior() {
  if [[ "$1" == REPROVADO || "$2" == REPROVADO ]]; then echo REPROVADO
  elif [[ "$1" == INCONCLUSIVO || "$2" == INCONCLUSIVO ]]; then echo INCONCLUSIVO
  else echo APROVADO; fi
}

caso() { # caso <id e descrição> <esperado> <obtido> <resultado>
  CASOS+=("$1"$'\t'"$2"$'\t'"$3"$'\t'"$4")
  printf '  %-12s %s\n' "[$4]" "$1"
  printf '  %-12s obtido: %s\n' "" "$3"
  log "CASO $1 | $2 | $3 | $4"
}

escrita_ok() { ESCRITAS+=("$1|$2|$3|$4|$5|$6"); }
recusada() { RECUSADAS+=("$1|$2|$3"); }
derivada() { DERIVADAS+=("$1|$2"); }

nome_rev() { local IFS=.; set -- $1; echo "$4.$3.$2.$1.in-addr.arpa."; }

reg() { # reg <zona> <tipo> <nome> <valor>
  jq -cn --arg z "$1" --arg t "$2" --arg n "$3" --arg v "$4" --argjson ttl "$TTL" \
    '{zone:$z, type:$t, name:$n, vl:$v, ttl:$ttl}'
}

# Leitura do estado sem tocar em ST/BODY.
api_admin_get() { curl -sS -m 40 -b "$ADMIN_JAR" -c "$ADMIN_JAR" "$BASE_URL$1" 2>>"$OUT_LOG"; }

# Valor de um RRset segundo a API (estado no PowerDNS, sem cache).
rr_api() {
  api_admin_get "/records?zone=$(q "$1")" |
    jq -r --arg n "$2" --arg t "$3" '[.record[]? | select(.type==$t and .name==$n) | .vl] | sort | join(",")' 2>/dev/null
}

# Valor de um RRset segundo o dig (o que o servidor publica).
rr_dig() {
  dig +short +time=2 +tries=2 @"$DNS_HOST" "$1" "$2" 2>/dev/null | grep -v '^;' | sort | paste -sd, -
}

dig_status() {
  dig +time=2 +tries=2 @"$DNS_HOST" "$1" "$2" 2>/dev/null | sed -n 's/.*status: \([A-Z]*\).*/\1/p' | head -1
}

# Repete o comando até devolver o valor esperado ou esgotar DIG_ESPERA.
espera_igual() {
  local esperado="$1" fim v
  shift
  fim=$(($(date +%s) + DIG_ESPERA))
  while :; do
    v=$("$@")
    [[ "$v" == "$esperado" ]] && break
    (($(date +%s) >= fim)) && break
    sleep 1
  done
  printf '%s' "$v"
}

# Lê o estado estável de um RRset: espera o dig convergir para o que a API diz.
# Define V_API e V_DIG. Se divergirem ao fim da espera, há cache ou falha de publicação.
le_estavel() { # <zona> <nome> <tipo>
  sleep "$PAUSA"
  V_API=$(rr_api "$1" "$2" "$3")
  V_DIG=$(espera_igual "$V_API" rr_dig "$2" "$3")
}

# Veredito sobre o valor publicado (dig), com a API como contraprova.
v_valor() { # <esperado>
  if [[ "$V_API" != "$V_DIG" ]]; then echo INCONCLUSIVO
  elif [[ "$V_DIG" == "$1" ]]; then echo APROVADO
  else echo REPROVADO; fi
}

mostra() { # mostra o valor publicado; acusa divergência com a API
  local v="${V_DIG:-(vazio)}"
  [[ "$V_API" != "$V_DIG" ]] && v="$v [API: ${V_API:-(vazio)}]"
  printf '%s' "$v"
}

zonas_reversas() {
  api_admin_get "/zones?type=reverse" | jq -r '[.zones[]?.name] | sort | join(",")' 2>/dev/null
}

# ----------------------------------------------------------------------------
# LIMPEZA (sempre roda, inclusive após erro ou Ctrl-C)
# ----------------------------------------------------------------------------
LIMPEZA_FEITA=0
limpeza() {
  local rc=$? i corpo
  [[ $LIMPEZA_FEITA == 1 ]] && return
  LIMPEZA_FEITA=1
  set +e
  echo
  info "== LIMPEZA"

  for i in 0 1 2; do
    if [[ ${OP_LOGADO[$i]} == 1 ]]; then
      req "${OP_JAR[$i]}" POST /logout
      info "  logout ${OP_ROT[$i]}: $ST"
    fi
  done

  if [[ $ADMIN_LOGADO == 1 ]]; then
    for i in 0 1 2; do
      if [[ -z "${OP_CAD[$i]}" && ${OP_SOLICITADO[$i]} == 1 ]]; then
        # A solicitação pode ter sido aprovada sem que o id do cadastro tenha sido lido.
        req "$ADMIN_JAR" GET "/cadastros?limit=50&filter=$(q "${OP_EMAIL[$i]}")"
        OP_CAD[$i]=$(jq -r --arg e "${OP_EMAIL[$i]}" '.cadastros[]? | select(.email==$e) | .id' <<<"$BODY" 2>/dev/null | head -1)
      fi
      if [[ -n "${OP_CAD[$i]}" ]]; then
        corpo=$(jq -cn --arg id "${OP_CAD[$i]}" '{id:$id}')
        req "$ADMIN_JAR" DELETE /cadastro "$corpo"
        info "  cadastro ${OP_ROT[$i]} removido (e permissões retiradas): $ST"
      fi
    done

    if [[ $ZR_TENTADA == 1 ]]; then
      req "$ADMIN_JAR" DELETE "/zone?id=$(q "$ZR")"
      info "  zona $ZR removida: $ST (dig: $(dig_status "$ZR" SOA))"
    fi
    if [[ $ZD_TENTADA == 1 ]]; then
      req "$ADMIN_JAR" DELETE "/zone?id=$(q "$ZD")"
      info "  zona $ZD removida: $ST (dig: $(dig_status "$ZD" SOA))"
    fi

    req "$ADMIN_JAR" POST /logout
    info "  logout admin: $ST"
  elif [[ $ZD_TENTADA == 1 || $ZR_TENTADA == 1 || ${OP_SOLICITADO[0]} == 1 ]]; then
    info "  ATENÇÃO: sem sessão de admin, nada pôde ser removido."
  fi

  # Resíduos que a API não permite remover.
  local resid=0
  for i in 0 1 2; do [[ ${OP_SOLICITADO[$i]} == 1 ]] && resid=1; done
  if [[ $resid == 1 ]]; then
    info "  Resíduos sem rota de remoção na API:"
    for i in 0 1 2; do
      [[ ${OP_SOLICITADO[$i]} == 1 ]] && info "    - solicitação de ${OP_EMAIL[$i]} (coleção solicitacoes)"
    done
    info "    - documento da zona $ZD na coleção users (vetores vazios após DELETE /cadastro)"
    info "    - entradas da trilha (coleção logs), mantidas de propósito"
    info "  Para remover as solicitações direto no MongoDB:"
    info "    db.solicitacoes.deleteMany({email: /^aval-e2-op[abc]-$RUN@example\\.com\$/})"
    info "    db.users.deleteOne({zona: \"$ZD\"})"
  fi

  rm -rf "$TMP"
  exit "$rc"
}
trap limpeza EXIT
trap 'exit 130' INT TERM

# ----------------------------------------------------------------------------
# PRÉ-VERIFICAÇÃO
# ----------------------------------------------------------------------------
: >"$OUT_LOG"
info "== E2 — autorização (QP3) e trilha de auditoria (QP4) · rodada $RUN"
info "   API: $BASE_URL · DNS: $DNS_HOST"

for b in curl jq dig; do command -v "$b" >/dev/null || die "falta o utilitário $b"; done
[[ -n "$ADMIN_EMAIL" && -n "$ADMIN_SENHA" ]] || die "defina ADMIN_EMAIL e ADMIN_SENHA no ambiente"

req - GET /healthcheck
[[ "$ST" == 200 ]] || die "healthcheck respondeu $(desc)"

# Instante inicial pelo relógio do servidor (cabeçalho Date), para filtrar a trilha.
DATA_SRV=$(curl -sS -m 20 -D - -o /dev/null "$BASE_URL/healthcheck" 2>/dev/null | awk 'tolower($1)=="date:"{$1=""; print}' | tr -d '\r' | sed 's/^ *//')
T0=$(jq -rn --arg d "$DATA_SRV" '$d | strptime("%a, %d %b %Y %H:%M:%S GMT") | mktime' 2>/dev/null) || T0=""
[[ "$T0" =~ ^[0-9]+$ ]] || T0=$(date +%s)
T0=$((T0 - 5))

[[ -n "$(dig_status . SOA)" ]] || die "o dig não obteve resposta de $DNS_HOST"

: >"$ADMIN_JAR"
req "$ADMIN_JAR" POST /login "$(jq -cn --arg e "$ADMIN_EMAIL" --arg s "$ADMIN_SENHA" '{email:$e, senha:$s}')"
[[ "$ST" == 200 ]] || die "login do admin falhou: $(desc)"
ADMIN_LOGADO=1
req "$ADMIN_JAR" GET /session
[[ "$(jq -r .level <<<"$BODY")" == admin ]] || die "a conta informada não é admin"

# Não tocar em zonas que já existam: a limpeza apagaria dados reais.
req "$ADMIN_JAR" GET "/zones?type=normal"
jq -e --arg z "$ZD" '.zones[]? | select(.name==$z)' <<<"$BODY" >/dev/null && die "a zona $ZD já existe; abortando"
REV_ANTES=$(zonas_reversas)
[[ ",$REV_ANTES," == *",$ZR,"* ]] && die "a zona $ZR já existe; abortando"
[[ "$(dig_status "$ZD" SOA)" == NOERROR ]] && die "o PowerDNS já responde por $ZD; abortando"
[[ "$(dig_status "$ZR" SOA)" == NOERROR ]] && die "o PowerDNS já responde por $ZR; abortando"
info "   pré-verificação ok (zonas de teste ausentes, dig respondendo)"

# ----------------------------------------------------------------------------
# SETUP
# ----------------------------------------------------------------------------
echo
info "== SETUP"

cria_zona() { # <domínio> <tipo>
  local corpo
  corpo=$(jq -cn --arg d "$1" --arg t "$2" --arg ns "ns1.$ZD" --arg em "hostmaster.$ZD" \
    '{domain:$d, type:$t, soa:{startOfAuthority:$ns, email:$em, refresh:3600, retry:600, expire:604800, negativeCacheTtl:3600}}')
  req "$ADMIN_JAR" PUT /zone "$corpo"
  [[ "$ST" == 200 ]] || die "criação da zona $1 falhou: $(desc)"
  escrita_ok dns create_zone "$ADMIN_EMAIL" "$1" "$1" "admin cria a zona $1"
  [[ "$(espera_igual NOERROR dig_status "$1" SOA)" == NOERROR ]] || die "o PowerDNS não publicou $1"
  info "  zona $1 criada ($2)"
}

ZD_TENTADA=1
cria_zona "$ZD" normal
ZR_TENTADA=1
cria_zona "$ZR" reverse

for i in 0 1 2; do
  rot="${OP_ROT[$i]}"
  email="${OP_EMAIL[$i]}"
  OP_SENHA[$i]=$(senha_aleatoria)

  # Fluxo real de entrada: solicitação pública → aprovação pelo admin.
  corpo=$(jq -cn --arg n "Avaliação E2 $rot" --arg e "$email" --arg s "${OP_SENHA[$i]}" \
    '{nome:$n, email:$e, senha:$s, foto:"avaliacao-e2"}')
  req - POST /cadastro "$corpo"
  [[ "$ST" == 200 ]] || die "POST /cadastro de $rot falhou: $(desc)"
  OP_SOLICITADO[$i]=1
  escrita_ok conta '*' "$email" - "$email" "$rot envia solicitação de cadastro"

  req "$ADMIN_JAR" GET "/solicitacoes?limit=50&filter=$(q "$email")"
  OP_SOLIC[$i]=$(jq -r --arg e "$email" '.solicitacoes[]? | select(.email==$e) | .id' <<<"$BODY" | head -1)
  [[ -n "${OP_SOLIC[$i]}" ]] || die "solicitação de $rot não encontrada"

  req "$ADMIN_JAR" PUT /solicitacoes/status "$(jq -cn --arg id "${OP_SOLIC[$i]}" '{id:$id, status:"aprovada"}')"
  [[ "$ST" == 200 ]] || die "aprovação de $rot falhou: $(desc)"
  escrita_ok conta '*' "$ADMIN_EMAIL" - "$email" "admin aprova a solicitação de $rot"

  req "$ADMIN_JAR" GET "/cadastros?limit=50&filter=$(q "$email")"
  OP_CAD[$i]=$(jq -r --arg e "$email" '.cadastros[]? | select(.email==$e) | .id' <<<"$BODY" | head -1)
  [[ -n "${OP_CAD[$i]}" ]] || die "cadastro de $rot não encontrado após aprovação"

  # Cookie jar próprio: /login recusa (400) se a requisição já traz cookie de sessão.
  : >"${OP_JAR[$i]}"
  req "${OP_JAR[$i]}" POST /login "$(jq -cn --arg e "$email" --arg s "${OP_SENHA[$i]}" '{email:$e, senha:$s}')"
  [[ "$ST" == 200 ]] || die "login de $rot falhou: $(desc)"
  OP_LOGADO[$i]=1
  req "${OP_JAR[$i]}" GET /session
  [[ "$(jq -r .level <<<"$BODY")" == member ]] || die "$rot não ficou com nível member"

  if [[ -n "${OP_PERM[$i]}" ]]; then
    req "$ADMIN_JAR" POST /user "$(jq -cn --arg e "$email" --arg p "${OP_PERM[$i]}" --arg z "$ZD" '{email:$e, permissao:$p, zona:$z}')"
    [[ "$ST" == 200 ]] || die "permissão de $rot falhou: $(desc)"
    escrita_ok permissao insert_user "$ADMIN_EMAIL" "$ZD" "$email" "admin concede ${OP_PERM[$i]} em $ZD a $rot"
    info "  $rot criado, aprovado e logado; ${OP_PERM[$i]} em $ZD"
  else
    info "  $rot criado, aprovado e logado; sem vínculo"
  fi
done
OPA="${OP_JAR[0]}" OPB="${OP_JAR[1]}" OPC="${OP_JAR[2]}"
EA="${OP_EMAIL[0]}" EB="${OP_EMAIL[1]}"

# PTR "de outra equipe" em 192.0.2.50, gravado pelo admin.
REV_OUTRA=$(nome_rev "$IP_OUTRA")
req "$ADMIN_JAR" PUT /records "$(reg "$ZR" PTR 50 "$PTR_OUTRA")"
[[ "$ST" == 200 ]] || die "semeadura do PTR falhou: $(desc)"
escrita_ok dns insert_record "$ADMIN_EMAIL" "$ZR" "$REV_OUTRA" "admin semeia PTR de outra equipe em $IP_OUTRA"
V_API=$(rr_api "$ZR" "$REV_OUTRA" PTR)
V_DIG=$(espera_igual "$PTR_OUTRA" rr_dig "$REV_OUTRA" PTR)
[[ "$V_DIG" == "$PTR_OUTRA" ]] || die "o dig não mostrou o PTR semeado (obtido: ${V_DIG:-vazio})"
info "  PTR semeado: $IP_OUTRA → $PTR_OUTRA"

# ----------------------------------------------------------------------------
# CASOS DE TESTE
# ----------------------------------------------------------------------------
echo
info "== CASOS"

# A0 — Atualização dinâmica não autenticada (QP3).
# Escrita sem cookie: o middleware de sessão deve barrar antes do controlador.
N="a0-$RUN.$ZD"
req - PUT /records "$(reg "$ZD" A "a0-$RUN" 192.0.2.3)"
r="APROVADO"; ob="$(desc)"
if [[ "$ST" == 000 ]]; then r=INCONCLUSIVO; elif ! erro_sessao; then r=REPROVADO; fi
le_estavel "$ZD" "$N" A
r=$(pior "$r" "$(v_valor "")")
caso "A0 — escrita sem sessão (PUT /records)" "erro de sessão; nada gravado" "$ob; A=$(mostra)" "$r"

# A1 — Privilégio excessivo (QP3): operador sem vínculo lê a zona.
req "$OPC" GET "/records?zone=$(q "$ZD")"
caso "A1 — opC (sem vínculo) lê a zona direta" "403" "$(desc)" "$(v_status 403)"

# A2 — Privilégio excessivo (QP3): permissão de leitura usada para escrever.
N="a2-$RUN.$ZD"
req "$OPB" PUT /records "$(reg "$ZD" A "a2-$RUN" 192.0.2.2)"
r=$(v_status 403); ob="$(desc)"
[[ "$ST" != 200 ]] && recusada "$EB" "a2-$RUN" "A2 opB PUT /records"
le_estavel "$ZD" "$N" A
r=$(pior "$r" "$(v_valor "")"); ob="$ob; A=$(mostra)"
le_estavel "$ZR" "$(nome_rev 192.0.2.2)" PTR
r=$(pior "$r" "$(v_valor "")"); ob="$ob; PTR=$(mostra)"
caso "A2 — opB (leitura) cria A na zona direta" "403; nada gravado (dig)" "$ob" "$r"

# A3 — Divergência A/PTR (QP2) dentro das condições da §3.4.5 (QP3):
# zona reversa existe e não há PTR → o PTR deve ser derivado.
N3="a3-$RUN.$ZD"
REV10=$(nome_rev 192.0.2.10)
req "$OPA" PUT /records "$(reg "$ZD" A "a3-$RUN" 192.0.2.10)"
r=$(v_status 200); ob="$(desc)"
[[ "$ST" == 200 ]] && escrita_ok dns insert_record "$EA" "$ZD" "$N3" "A3 opA cria A $N3 → 192.0.2.10"
V_API=$(rr_api "$ZD" "$N3" A); V_DIG=$(espera_igual "192.0.2.10" rr_dig "$N3" A)
r=$(pior "$r" "$(v_valor 192.0.2.10)"); ob="$ob; A=$(mostra)"
V_API=$(espera_igual "$N3" rr_api "$ZR" "$REV10" PTR); V_DIG=$(espera_igual "$N3" rr_dig "$REV10" PTR)
r=$(pior "$r" "$(v_valor "$N3")"); ob="$ob; PTR=$(mostra)"
[[ "$V_API" == "$N3" ]] && derivada "$REV10" "A3: PTR 192.0.2.10 criado pela automação"
caso "A3 — opA cria A 192.0.2.10 sem PTR prévio" "200; PTR criado → $N3" "$ob" "$r"

# A4 — Privilégio excessivo pela escrita derivada (QP3, §3.4.5 condição 2):
# já existe PTR de outra equipe → não pode ser sobrescrito.
N4="a4-$RUN.$ZD"
req "$OPA" PUT /records "$(reg "$ZD" A "a4-$RUN" "$IP_OUTRA")"
r=$(v_status 200); ob="$(desc)"
[[ "$ST" == 200 ]] && escrita_ok dns insert_record "$EA" "$ZD" "$N4" "A4 opA cria A $N4 → $IP_OUTRA"
le_estavel "$ZR" "$REV_OUTRA" PTR
r=$(pior "$r" "$(v_valor "$PTR_OUTRA")"); ob="$ob; PTR=$(mostra)"
caso "A4 — opA cria A no IP do PTR de outra equipe" "200; PTR intacto → $PTR_OUTRA" "$ob" "$r"

# A5 — Privilégio excessivo pela escrita derivada (QP3, §3.4.5 condição 1):
# faixa sem zona reversa → nenhuma zona e nenhum PTR criados.
N5="a5-$RUN.$ZD"
REV_SEM=$(nome_rev "$IP_SEM_ZONA")
COBERTA=0
for z in $(tr ',' ' ' <<<"$REV_ANTES"); do
  [[ "$REV_SEM" == *".$z" || "$REV_SEM" == "$z" ]] && COBERTA=1
done
req "$OPA" PUT /records "$(reg "$ZD" A "a5-$RUN" "$IP_SEM_ZONA")"
r=$(v_status 200); ob="$(desc)"
[[ "$ST" == 200 ]] && escrita_ok dns insert_record "$EA" "$ZD" "$N5" "A5 opA cria A $N5 → $IP_SEM_ZONA"
sleep "$PAUSA"
REV_DEPOIS=$(zonas_reversas)
NOVAS=$(comm -13 <(tr ',' '\n' <<<"$REV_ANTES,$ZR" | sort -u) <(tr ',' '\n' <<<"$REV_DEPOIS" | sort -u) | paste -sd, -)
P5=$(rr_dig "$REV_SEM" PTR); S5=$(dig_status "$REV_SEM" PTR)
[[ -n "$NOVAS" || -n "$P5" ]] && r=REPROVADO
[[ $COBERTA == 1 ]] && r=INCONCLUSIVO
ob="$ob; zonas novas: ${NOVAS:-nenhuma}; PTR=${P5:-(vazio)} ($S5)"
[[ $COBERTA == 1 ]] && ob="$ob; faixa já coberta por zona existente"
caso "A5 — opA cria A em $IP_SEM_ZONA (faixa sem zona reversa)" "200; nenhuma zona e nenhum PTR" "$ob" "$r"

# A6 — Privilégio excessivo (QP3): o operador não lista nem edita a zona reversa.
req "$OPA" GET "/records?zone=$(q "$ZR")"
caso "A6a — opA lista a zona reversa (GET /records)" "403" "$(desc)" "$(v_status 403)"

REV60=$(nome_rev 192.0.2.60)
req "$OPA" PUT /records "$(reg "$ZR" PTR 60 "a6-$RUN.$ZD")"
r=$(v_status 403); ob="$(desc)"
[[ "$ST" != 200 ]] && recusada "$EA" "$REV60" "A6b opA PUT /records na zona reversa"
le_estavel "$ZR" "$REV60" PTR
r=$(pior "$r" "$(v_valor "")"); ob="$ob; PTR=$(mostra)"
caso "A6b — opA cria PTR na zona reversa (PUT /records)" "403; nada gravado" "$ob" "$r"

req "$OPA" DELETE /records "$(reg "$ZR" PTR 50 "$PTR_OUTRA")"
r=$(v_status 403); ob="$(desc)"
[[ "$ST" != 200 ]] && recusada "$EA" "$REV_OUTRA" "A6c opA DELETE /records na zona reversa"
le_estavel "$ZR" "$REV_OUTRA" PTR
r=$(pior "$r" "$(v_valor "$PTR_OUTRA")"); ob="$ob; PTR=$(mostra)"
caso "A6c — opA apaga o PTR de outra equipe (DELETE /records)" "403; PTR intacto" "$ob" "$r"

# A7 — Privilégio excessivo (QP3): triagem e remoção de órfãos exigem permissão na zona reversa.
req "$OPA" GET "/reverses/orphans?zone=$(q "$ZR")"
caso "A7 — opA consulta órfãos da zona reversa" "403" "$(desc)" "$(v_status 403)"

req "$OPA" DELETE "/reverses?zone=$(q "$ZR")&name=$(q "$REV_OUTRA")"
r=$(v_status 403); ob="$(desc)"
[[ "$ST" != 200 ]] && recusada "$EA" "$REV_OUTRA" "A7b opA DELETE /reverses"
le_estavel "$ZR" "$REV_OUTRA" PTR
r=$(pior "$r" "$(v_valor "$PTR_OUTRA")"); ob="$ob; PTR=$(mostra)"
caso "A7b — opA remove o PTR como órfão (DELETE /reverses)" "403; PTR intacto" "$ob" "$r"

# A8 — Privilégio excessivo pela escrita derivada na ALTERAÇÃO de endereço (QP3, §6.2.1).
# UpdateReverso grava o PTR do novo IP com REPLACE, sem a checagem de PTR existente
# que InsertReverso faz. O critério do Quadro 3 é "zero sobrescritas".
N8="a8-$RUN.$ZD"
REV20=$(nome_rev 192.0.2.20)
req "$OPA" PUT /records "$(reg "$ZD" A "a8-$RUN" 192.0.2.20)"
[[ "$ST" == 200 ]] || die "preparação do A8 falhou: $(desc)"
escrita_ok dns insert_record "$EA" "$ZD" "$N8" "A8 (preparo) opA cria A $N8 → 192.0.2.20"
V_API=$(espera_igual "$N8" rr_api "$ZR" "$REV20" PTR)
[[ "$V_API" == "$N8" ]] && derivada "$REV20" "A8 (preparo): PTR 192.0.2.20 criado pela automação"
corpo=$(jq -cn --argjson o "$(reg "$ZD" A "a8-$RUN" 192.0.2.20)" --argjson n "$(reg "$ZD" A "a8-$RUN" "$IP_OUTRA")" '{oldValue:$o, newValue:$n}')
req "$OPA" PATCH /records "$corpo"
r=$(v_status 200); ob="$(desc)"
[[ "$ST" == 200 ]] && escrita_ok dns edit_record "$EA" "$ZD" "a8-$RUN" "A8 opA altera $N8 de 192.0.2.20 para $IP_OUTRA"
le_estavel "$ZR" "$REV_OUTRA" PTR
r=$(pior "$r" "$(v_valor "$PTR_OUTRA")")
if [[ "$V_API" == "$PTR_OUTRA" ]]; then
  ob="$ob; PTR $IP_OUTRA mantido → $(mostra)"
else
  ob="$ob; PTR $IP_OUTRA SOBRESCRITO → $(mostra)"
  [[ "$V_API" == "$N8" ]] && derivada "$REV_OUTRA" "A8: PTR $IP_OUTRA da outra equipe sobrescrito pela automação"
fi
le_estavel "$ZR" "$REV20" PTR
if [[ -z "$V_API" ]]; then
  ob="$ob; PTR 192.0.2.20 removido"
  derivada "$REV20" "A8: PTR 192.0.2.20 removido pela automação"
else
  ob="$ob; PTR 192.0.2.20 → $(mostra)"
fi
caso "A8 — opA muda o IP de um host para $IP_OUTRA (PATCH /records)" "200; PTR de outra equipe mantido (QP3)" "$ob" "$r"

# A9 — Divergência A/PTR (QP2): guarda de referência.
# a9 passa a usar 192.0.2.10; apaga-se a3, o nome para o qual o PTR aponta.
# O PTR deve permanecer porque o IP ainda tem registro direto.
N9="a9-$RUN.$ZD"
req "$OPA" PUT /records "$(reg "$ZD" A "a9-$RUN" 192.0.2.10)"
[[ "$ST" == 200 ]] || die "preparação do A9 falhou: $(desc)"
escrita_ok dns insert_record "$EA" "$ZD" "$N9" "A9 (preparo) opA cria A $N9 → 192.0.2.10"
req "$OPA" DELETE /records "$(reg "$ZD" A "a3-$RUN" 192.0.2.10)"
r=$(v_status 200); ob="$(desc)"
[[ "$ST" == 200 ]] && escrita_ok dns delete_record "$EA" "$ZD" "$N3" "A9 opA apaga A $N3 (192.0.2.10 segue em $N9)"
V_API=$(rr_api "$ZD" "$N3" A); V_DIG=$(espera_igual "" rr_dig "$N3" A)
r=$(pior "$r" "$(v_valor "")"); ob="$ob; A de a3=$(mostra)"
le_estavel "$ZR" "$REV10" PTR
if [[ -n "$V_API" && "$V_API" == "$V_DIG" ]]; then
  ob="$ob; PTR mantido → $(mostra)"
  [[ "$V_DIG" == "$N3" ]] && ob="$ob (aponta para o nome recém-apagado)"
elif [[ "$V_API" != "$V_DIG" ]]; then
  r=$(pior "$r" INCONCLUSIVO); ob="$ob; PTR=$(mostra)"
else
  r=REPROVADO; ob="$ob; PTR removido"
fi
caso "A9 — opA apaga A cujo IP (192.0.2.10) segue em uso por outro nome" "200; PTR permanece" "$ob" "$r"

# A10 — Registros pendentes e reversos órfãos (QP2): remoção encadeada do último A do IP.
req "$OPA" DELETE /records "$(reg "$ZD" A "a9-$RUN" 192.0.2.10)"
r=$(v_status 200); ob="$(desc)"
[[ "$ST" == 200 ]] && escrita_ok dns delete_record "$EA" "$ZD" "$N9" "A10 opA apaga o último A de 192.0.2.10 ($N9)"
V_API=$(espera_igual "" rr_api "$ZR" "$REV10" PTR); V_DIG=$(espera_igual "" rr_dig "$REV10" PTR)
r=$(pior "$r" "$(v_valor "")"); ob="$ob; PTR=$(mostra)"
[[ -z "$V_API" ]] && derivada "$REV10" "A10: PTR 192.0.2.10 removido pela automação"
caso "A10 — opA apaga o último A de 192.0.2.10" "200; PTR removido" "$ob" "$r"

# A11 — Registros pendentes e reversos órfãos (QP2): a triagem aponta o PTR sem A.
REV99=$(nome_rev 192.0.2.99)
req "$ADMIN_JAR" PUT /records "$(reg "$ZR" PTR 99 "orfao-$RUN.$ZD")"
r=$(v_status 200); ob="$(desc)"
[[ "$ST" == 200 ]] && escrita_ok dns insert_record "$ADMIN_EMAIL" "$ZR" "$REV99" "A11 admin cria PTR órfão em 192.0.2.99"
sleep "$PAUSA"
req "$ADMIN_JAR" GET "/reverses/orphans?zone=$(q "$ZR")"
ob="$ob; GET orphans $(desc)"
LISTA=$(jq -r '[.[]?.ip] | join(",")' <<<"$BODY" 2>/dev/null)
if [[ "$ST" == 200 ]] && jq -e --arg n "$REV99" 'any(.[]?; .nomeReverso==$n)' <<<"$BODY" >/dev/null 2>&1; then
  ob="$ob; órfãos: $LISTA"
else
  r=$(pior "$r" REPROVADO); ob="$ob; órfãos: ${LISTA:-(nenhum)}"
fi
caso "A11 — admin cria PTR sem A e consulta órfãos" "192.0.2.99 listado como órfão" "$ob" "$r"

# A12 — Privilégio excessivo (QP3): escape de escopo por FQDN.
# NomeCompleto preserva nomes terminados em ponto; com escrita só na zona direta,
# opA tenta gravar um PTR sob in-addr.arpa. Espera-se recusa (o PowerDNS rejeita
# nome fora da zona) e nada publicado.
REV61=$(nome_rev 192.0.2.61)
req "$OPA" PUT /records "$(reg "$ZD" PTR "$REV61" "a12-$RUN.$ZD")"
ob="$(desc)"
if [[ "$ST" == 200 ]]; then r=REPROVADO
elif erro_sessao || [[ "$ST" == 000 ]]; then r=INCONCLUSIVO
else r=APROVADO; recusada "$EA" "$REV61" "A12 opA PUT /records com FQDN fora da zona"; fi
le_estavel "$ZR" "$REV61" PTR
r=$(pior "$r" "$(v_valor "")"); ob="$ob; PTR=$(mostra)"
caso "A12 — opA grava PTR sob in-addr.arpa via FQDN na zona direta" "recusado (≠200); nada gravado" "$ob" "$r"

# ----------------------------------------------------------------------------
# TRILHA DE AUDITORIA (QP4) — Ausência de rastreabilidade das alterações
# ----------------------------------------------------------------------------
echo
info "== TRILHA DE AUDITORIA"
sleep "$ESPERA_LOGS"

# Junta as entradas que mencionam as zonas de teste ou a rodada, a partir de T0.
: >"$TMP/logs.ndjson"
for f in "aval-sanchezdns" "2.0.192.in-addr.arpa" "$RUN"; do
  api_admin_get "/logs?page=1&limit=1000&filter=$(q "$f")" | jq -c '.logs[]?' >>"$TMP/logs.ndjson" 2>/dev/null
done
JQ_EP='def ep: if (.|type) != "string" or length < 19 then null else
  (.[0:19] + "Z" | fromdateiso8601) as $b | (.[19:]) as $o |
  if $o == "" or $o == "Z" then $b
  else (($o[1:3]|tonumber)*3600 + ($o[4:6]|tonumber)*60) as $d | (if $o[0:1] == "+" then $b - $d else $b + $d end) end end;'
jq -s --argjson t0 "$T0" "$JQ_EP"' unique | map(select((.createdAt|ep) != null and (.createdAt|ep) >= $t0))' \
  "$TMP/logs.ndjson" >"$TMP/logs.json"
info "  entradas da rodada na trilha: $(jq length "$TMP/logs.json")"

TRILHA=()  # "descrição<TAB>autor<TAB>ação<TAB>zona<TAB>instante<TAB>registrada"
X=0; Y=0; XD=0; YD=0; XC=0; YC=0
i=0
while [[ $i -lt ${#ESCRITAS[@]} ]]; do
  IFS='|' read -r cat acao autor zona chave descr <<<"${ESCRITAS[$i]}"
  i=$((i + 1))
  ent=$(jq -c --arg a "$acao" --arg u "$autor" --arg k "$chave" '
    [ .[] | select(.username == $u)
          | select(((.details // "") | contains($k)) or ((.zone // "") | contains($k)))
          | select(if $a == "*" then ((.action // "") | test("_user$") | not) else .action == $a end) ]
    | first // empty' "$TMP/logs.json")
  if [[ -n "$ent" ]]; then
    c_aut=$(jq -r 'if (.username // "") != "" then "sim" else "não" end' <<<"$ent")
    c_aca=$(jq -r 'if (.action // "") != "" then "sim (\(.action))" else "não" end' <<<"$ent")
    if [[ "$zona" == "-" ]]; then
      c_zon="n/a"
    else
      c_zon=$(jq -r --arg z "${zona%.}" '
        if ((.zone // "") | contains($z)) then "sim (campo zone)"
        elif ((.details // "") | contains($z)) then "só no texto de details (zone=\(.zone))"
        else "não" end' <<<"$ent")
    fi
    c_ins=$(jq -r "$JQ_EP"' if (.createdAt|ep) != null then "sim" else "não" end' <<<"$ent")
    ok=sim
    [[ "$c_aut" == não || "$c_aca" == não || "$c_zon" == não || "$c_ins" == não ]] && ok=não
  else
    c_aut="—"; c_aca="—"; c_zon="—"; c_ins="—"; ok=não
  fi
  [[ "$cat" == conta ]] && c_aca="${c_aca/—/— (sem código no Quadro 9)}"
  TRILHA+=("$descr"$'\t'"$c_aut"$'\t'"$c_aca"$'\t'"$c_zon"$'\t'"$c_ins"$'\t'"$ok")
  Y=$((Y + 1)); [[ $ok == sim ]] && X=$((X + 1))
  if [[ "$cat" == conta ]]; then
    YC=$((YC + 1)); [[ $ok == sim ]] && XC=$((XC + 1))
  else
    YD=$((YD + 1)); [[ $ok == sim ]] && XD=$((XD + 1))
  fi
  printf '  %-4s %s\n' "$ok" "$descr"
done

# Alterações derivadas: o PTR mudou no servidor sem requisição própria do operador.
# Procura qualquer entrada não-admin que cite o nome reverso.
DER=(); XV=0; YV=0
i=0
while [[ $i -lt ${#DERIVADAS[@]} ]]; do
  IFS='|' read -r rev descr <<<"${DERIVADAS[$i]}"
  i=$((i + 1))
  n=$(jq --arg k "$rev" --arg adm "$ADMIN_EMAIL" \
    '[.[] | select(.username != $adm) | select(((.details // "") | contains($k)) or ((.zone // "") | contains($k)))] | length' "$TMP/logs.json")
  YV=$((YV + 1))
  if [[ "$n" -gt 0 ]]; then XV=$((XV + 1)); DER+=("$descr"$'\t'"sim"); else DER+=("$descr"$'\t'"não"); fi
done

# Tentativas recusadas não podem aparecer na trilha.
REC=(); NR=0
i=0
while [[ $i -lt ${#RECUSADAS[@]} ]]; do
  IFS='|' read -r autor chave descr <<<"${RECUSADAS[$i]}"
  i=$((i + 1))
  n=$(jq --arg u "$autor" --arg k "$chave" \
    '[.[] | select(.username == $u) | select(((.details // "") | contains($k)) or ((.zone // "") | contains($k)))] | length' "$TMP/logs.json")
  [[ "$n" -gt 0 ]] && NR=$((NR + 1))
  REC+=("$descr"$'\t'"$n")
done

info ""
info "  Trilha: $X de $Y alterações registradas"
info "    zonas, registros e permissões: $XD de $YD"
info "    contas (cadastro e aprovação): $XC de $YC"
info "  Alterações derivadas (PTR pela automação): $XV de $YV registradas"
info "  Tentativas recusadas presentes na trilha: $NR de ${#RECUSADAS[@]} (esperado 0)"

# ----------------------------------------------------------------------------
# SAÍDA
# ----------------------------------------------------------------------------
APR=0; REP=0; INC=0
for l in "${CASOS[@]}"; do
  case "${l##*$'\t'}" in APROVADO) APR=$((APR + 1)) ;; REPROVADO) REP=$((REP + 1)) ;; *) INC=$((INC + 1)) ;; esac
done

echo
info "== RESULTADO"
{
  printf 'Caso\tEsperado\tObtido\tResultado\n'
  for l in "${CASOS[@]}"; do printf '%s\n' "$l"; done
} | column -t -s $'\t' | tee -a "$OUT_LOG"
info ""
info "  $APR aprovados · $REP reprovados · $INC inconclusivos"

md_linha() { local s="${1//|/\\|}"; echo "| ${s//$'\t'/ | } |"; }

{
  echo "# Resultado do E2 — autorização (QP3) e trilha de auditoria (QP4)"
  echo
  echo "Execução: $(date '+%Y-%m-%d %H:%M %z') · rodada \`$RUN\` · zonas descartáveis \`$ZD\` e \`$ZR\` (RFC 2606 e RFC 5737) · verificação externa por \`dig\` contra o servidor autoritativo."
  echo
  echo "## Casos de teste"
  echo
  echo "| Caso | Esperado | Obtido | Resultado |"
  echo "| --- | --- | --- | --- |"
  for l in "${CASOS[@]}"; do md_linha "$l"; done
  echo
  echo "Total: $APR aprovados, $REP reprovados, $INC inconclusivos."
  echo
  echo "## Classes do Quadro 1 por caso"
  echo
  echo "| Caso | Classe de erro | QP |"
  echo "| --- | --- | --- |"
  echo "| A0 | Atualização dinâmica não autenticada | QP3 |"
  echo "| A1, A2, A6a–A6c, A7, A7b, A12 | Privilégio excessivo do operador | QP3 |"
  echo "| A3 | Divergência entre A/AAAA e PTR (dentro das condições da §3.4.5) | QP2, QP3 |"
  echo "| A4, A5, A8 | Privilégio excessivo pela escrita derivada (§3.4.5, §6.2.1) | QP3 |"
  echo "| A9 | Divergência entre A/AAAA e PTR (guarda de referência) | QP2 |"
  echo "| A10, A11 | Registros pendentes e reversos órfãos | QP2 |"
  echo "| Trilha | Ausência de rastreabilidade das alterações | QP4 |"
  echo
  echo "## Trilha de auditoria"
  echo
  echo "**$X de $Y alterações registradas** com autor, ação, zona e instante (zonas, registros e permissões: $XD de $YD; contas: $XC de $YC)."
  echo "Alterações derivadas (PTR gravado ou removido pela automação): $XV de $YV registradas."
  echo "Tentativas recusadas presentes na trilha: $NR de ${#RECUSADAS[@]}."
  echo
  echo "| Alteração | Autor | Ação | Zona | Instante | Registrada |"
  echo "| --- | --- | --- | --- | --- | --- |"
  for l in "${TRILHA[@]}"; do md_linha "$l"; done
  if [[ ${#DER[@]} -gt 0 ]]; then
    echo
    echo "| Alteração derivada | Registrada |"
    echo "| --- | --- |"
    for l in "${DER[@]}"; do md_linha "$l"; done
  fi
  if [[ ${#REC[@]} -gt 0 ]]; then
    echo
    echo "| Tentativa recusada | Entradas na trilha |"
    echo "| --- | --- |"
    for l in "${REC[@]}"; do md_linha "$l"; done
  fi
} >"$OUT_MD"

info ""
info "  tabela: $OUT_MD"
info "  registro bruto: $OUT_LOG"
