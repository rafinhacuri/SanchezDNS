# 📦 Zones

## Overview
The **Zones** page in SanchezDNS provides a centralized interface for managing DNS zones and their associated records across your connected PowerDNS servers.  
Each zone represents a distinct DNS domain, and managing zones here allows you to define and update DNS records efficiently.

Creating a zone requires defining a Start of Authority (SOA) record, which is essential for DNS zone management and delegation.

---

## 📝 What is an SOA Record?
The SOA (Start of Authority) record is a critical DNS record that indicates the authoritative information about a DNS zone. It contains metadata about the zone and controls how DNS data is propagated and refreshed.  

Each SOA record includes the following fields:

- **MNAME:** The primary master name server for the zone. This is the hostname of the server that contains the original zone data.
- **RNAME (Email):** The email address of the person responsible for the zone. In DNS format, the `@` symbol is replaced by a dot (`.`). For example, `admin.example.com` corresponds to `admin@example.com`.
- **REFRESH:** The time interval (in seconds) that secondary name servers wait before querying the master server for zone updates.
- **RETRY:** The time interval (in seconds) that secondary servers wait before retrying a failed zone transfer.
- **EXPIRE:** The time interval (in seconds) that secondary servers keep the zone data before discarding it if unable to reach the master.
- **NEGATIVE TTL:** The time interval (in seconds) that DNS resolvers cache negative responses (e.g., non-existent domain).

---

## ⚙️ Creating and Managing Records
Within the **Zones** page, you can add, edit, or delete DNS records for each zone. The interface is designed to simplify DNS management:

- When creating a record, domain names are automatically completed based on the zone's root domain.
- Use the symbol “@” to represent the root of the zone (e.g., `example.com`).
- Supported record types include:  
  `['A', 'AAAA', 'ALIAS', 'CAA', 'CNAME', 'HTTPS', 'MX', 'NS', 'TXT', 'SRV']`.
- The UI guides you through entering valid record data but assumes familiarity with DNS concepts; it is intended for DNS professionals managing zones, not for DNS education.

---

## 🔧 Best Practices
- Always verify DNS propagation after making changes to ensure updates have taken effect globally.
- Consider appropriate TTL (Time To Live) values to balance between update speed and caching efficiency.
- Maintain authoritative consistency by ensuring SOA and NS records are correctly configured and synchronized across your DNS infrastructure.

---
SanchezDNS streamlines DNS zone management by providing a professional, unified interface to control your DNS data across multiple PowerDNS servers.