/* Do not change, this code is generated from Golang structs */

export interface GoRes {
  message: string
}
export interface Soa {
  startOfAuthority: string
  email: string
  refresh: number
  retry: number
  expire: number
  negativeCacheTtl: number
}
export interface CreateZoneRequest {
  domain: string
  soa: Soa
  type: string
}
export interface Log {
  zone: string
  username: string
  action: string
  details: string
  createdAt: string
}
export interface LogsResponse {
  logs: Log[]
  total: number
}
export interface PdnsZone {
  name: string
  id: string
  kind: string
  serial: number
  url: string
  soa_edit_api: string
}
export interface SessionRes {
  email: string
  level: string
}
export interface StatisticsResponse {
  zones: number
  records: number
  uptime: string
  status: string
  udpQueries: number
  tcpQueries: number
  serverId: string
  startedAt: string
}
export interface ZoneFetch {
  name: string
  serial: number
  nivel: string
  dnssec: boolean
}
export interface ZonesResponse {
  zones: ZoneFetch[]
}

export interface User {
  email: string
  permissao: string
  zona: string
  id: string
}
export interface Solicitacao {
  id: string
  nome: string
  email: string
  senha: string
  foto: string
  status: string
  createdAt: string
  updatedAt: string
}
export interface SolicitacaoResponse {
  solicitacoes: Solicitacao[]
  total: number
}
export interface Cadastro {
  id: string
  nome: string
  email: string
  senha: string
  foto: string
  level: string
  createdAt: string
  updatedAt: string
}
export interface CadastroResponse {
  cadastros: Cadastro[]
  total: number
}
