# 🧭 Zones

## Overview
The **Zones** section in SanchezDNS manages DNS zones, each representing a distinct domain. Users can create, edit, or delete zones directly connected to PowerDNS through the API, enabling seamless synchronization between the SanchezDNS interface and the authoritative DNS servers.

## SOA (Start of Authority)
Every new zone must include a SOA record, which defines the authoritative information about the zone and controls DNS replication behavior between primary and secondary servers.

The SOA record consists of the following parameters:

- **Primary Nameserver (MNAME)** — The authoritative server responsible for the zone.
- **Email (RNAME)** — The contact email for the zone administrator, written as `rafael.curi.dev.br` in the interface but interpreted as `rafael@curi.dev.br`.
- **Serial** — A version number incremented on each update to signal changes to secondary servers.
- **Refresh** — The interval (in seconds) for secondary servers to check for zone updates, typically set between 900 and 3600 seconds.
- **Retry** — The time secondary servers wait before retrying a failed zone transfer, usually between 600 and 1800 seconds.
- **Expire** — The time after which secondary servers stop responding if unable to contact the primary server, often set between 86400 and 604800 seconds.
- **Negative TTL** — The duration for which negative responses (e.g., NXDOMAIN) are cached, generally ranging from 60 to 3600 seconds.

Together, these parameters establish the authority of the zone and control how DNS data is replicated and refreshed across servers.

## Creating and Managing Records
When adding DNS records, the interface provides intuitive behavior to simplify management:

- If the zone is `teste.curi.dev.br` and the user enters a record name such as `rafael`, the system automatically expands it to `rafael.teste.curi.dev.br`.
- If no name is provided, the record is stored at the root of the zone and displayed as `@` in the records table.

Users can edit TTL, content, comments, and other record details directly within the interface.

SanchezDNS supports the following DNS record types:

```
A, AAAA, ALIAS, CAA, CNAME, HTTPS, MX, NS, TXT, SRV
```

- **A** — Maps a hostname to an IPv4 address.
- **AAAA** — Maps a hostname to an IPv6 address.
- **ALIAS** — Provides aliasing functionality similar to CNAME but at the zone apex.
- **CAA** — Specifies which certificate authorities are allowed to issue certificates for the domain.
- **CNAME** — Creates an alias from one hostname to another.
- **HTTPS** — Provides HTTPS service binding information for the domain.
- **MX** — Defines mail exchange servers for the domain.
- **NS** — Specifies authoritative name servers for the zone.
- **TXT** — Holds arbitrary text data, often used for verification or policy records.
- **SRV** — Defines service location records for specific protocols.

## Notes
SanchezDNS is designed for DNS professionals and does not provide DNS concept tutorials. It automates record management by interfacing directly with the PowerDNS Authoritative API, ensuring that all changes are applied immediately and accurately.

## Best Practices
- Verify SOA values carefully and ensure that secondary servers synchronize correctly.
- Use TTL values thoughtfully to balance DNS propagation speed and cache efficiency.
- Always confirm DNS changes with external lookups before deploying them live to avoid disruptions.