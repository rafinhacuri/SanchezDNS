# ⚙️ Configuration

This page explains how to properly configure your PowerDNS server and how SanchezDNS connects to it securely and reliably.

---

## 🧩 Supported Servers

SanchezDNS currently supports **only PowerDNS Authoritative Servers**.  
Other DNS software such as BIND or CoreDNS is not compatible.

Your PowerDNS instance must be configured to expose its **API** and **webserver** endpoints so that SanchezDNS can communicate with it.

---

## 🌐 Enabling the PowerDNS API

To allow SanchezDNS to connect, make sure the following parameters are set in your PowerDNS configuration file:

**File:** `/etc/powerdns/pdns.conf`


```ini
api=yes
api-key=YourSecureAPIKeyHere
webserver=yes
webserver-address=0.0.0.0
webserver-allow-from=0.0.0.0/0
webserver-port=8081
server-id=localhost
```

> ⚠️ **Security Warning:**  
> Avoid using `webserver-allow-from=0.0.0.0/0` in production.  
> This allows anyone to access your PowerDNS API.  
> Instead, restrict access to your SanchezDNS server IP only:
> ```ini
> webserver-allow-from=YOUR_SYSTEM_IP/32
> ```

> 🔐 **Tip:** The `api-key` must match the one you register inside SanchezDNS when creating a new connection.

---

## 🔒 Enabling DNSSEC (Optional)

If you plan to use **DNSSEC** (Domain Name System Security Extensions),  
its activation depends on the backend configured in your PowerDNS server.  
For example, if you are using the **SQLite** backend, you must enable DNSSEC like this:

**File:** `/etc/powerdns/pdns.conf`

```ini
launch=gsqlite3
gsqlite3-database=/var/lib/powerdns/pdns.sqlite3
gsqlite3-dnssec=yes
```

You can then manage DNSSEC-enabled zones directly from the SanchezDNS interface.  
Once DNSSEC is active, SanchezDNS will automatically display and track DNSSEC status for your zones.

---

## 🔗 Creating a Connection

In SanchezDNS, go to the **Connections** page and click **Add New Connection**.  
Fill in the following details:

| Field | Description |
|-------|--------------|
| **Name** | A friendly name to identify your server (e.g. “Authoritative DNS - Primary”). |
| **Host** | The IP or hostname of your PowerDNS server (e.g. `152.84.120.200`). |
| **Server ID** | Typically `localhost`, unless you use a custom setup. |
| **API Key** | The same API key you defined in `/etc/powerdns/pdns.conf`. |

---

## 🧠 Connectivity Check

SanchezDNS automatically tests the connection before saving it.  
If the system can **ping** the target host and validate the **PowerDNS API**,  
the connection will be saved successfully. Otherwise, you’ll receive an error message.

> ✅ Tip: Make sure the PowerDNS API port (default **8081**) is open in your firewall.

---

## 🧱 Recommended Configuration

To ensure stable operation, make sure your PowerDNS server includes:

```ini
local-address=0.0.0.0
local-port=53
launch=gsqlite3
gsqlite3-database=/var/lib/powerdns/pdns.sqlite3
default-ttl=3600
```

These settings allow your server to respond to DNS queries and store zone data correctly.

---

## 🧾 Summary

- SanchezDNS supports **PowerDNS Authoritative Servers only**.  
- Ensure **API** and **webserver** are enabled in `/etc/powerdns/pdns.conf`.  
- Enable **DNSSEC** manually if desired.  
- The connection will only be established if the host responds to ping and the API key is valid.  
- Once connected, SanchezDNS provides full control over zones, records, users, and logs — all from one place.
