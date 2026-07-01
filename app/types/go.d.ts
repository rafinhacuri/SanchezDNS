/* Do not change, this code is generated from Golang structs */

interface GoRes {
  message: string
}
interface Soa {
  startOfAuthority: string
  email: string
  refresh: number
  retry: number
  expire: number
  negativeCacheTtl: number
}
interface CreateZoneRequest {
  domain: string
  soa: Soa
  type: string
}
interface Log {
  zone: string
  username: string
  action: string
  details: string
  createdAt: string
}
interface LogsResponse {
  logs: Log[]
  total: number
}
interface PdnsZone {
  name: string
  id: string
  kind: string
  serial: number
  url: string
  soa_edit_api: string
}
interface SessionRes {
  email: string
  level: string
}
interface StatisticsResponse {
  zones: number
  records: number
  uptime: string
  status: string
  udpQueries: number
  tcpQueries: number
  serverId: string
  startedAt: string
}
interface ZoneFetch {
  name: string
  serial: number
  nivel: string
  dnssec: boolean
}
interface ZonesResponse {
  zones: ZoneFetch[]
}

interface User {
  email: string
  permissao: string
  zona: string
  id: string
}
interface Solicitacao {
  id: string
  nome: string
  email: string
  senha: string
  foto: string
  status: string
  createdAt: string
  updatedAt: string
}
interface SolicitacaoResponse {
  solicitacoes: Solicitacao[]
  total: number
}
interface Cadastro {
  id: string
  nome: string
  email: string
  senha: string
  foto: string
  level: string
  createdAt: string
  updatedAt: string
}
interface CadastroResponse {
  cadastros: Cadastro[]
  total: number
}
