<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=venom&height=300&color=2FC851&text=SanchezDNS&section=header&fontColor=ffffff"/>
</p>

---

## 💡 Why Choose SanchezDNS

> _A next-generation DNS management platform that merges simplicity, security, and performance._

SanchezDNS is crafted for professionals who demand both **efficiency** and **power** in DNS management — whether you’re a hosting provider, enterprise admin, or infrastructure engineer.

---

### 🚀 Key Highlights

#### 🧭 Unified Control Panel

Manage **all zones, records, and users** in one elegant interface, powered by real-time PowerDNS integration.

#### 🔒 Enterprise-Grade Security

- AES‑256 encryption for sensitive data (API keys, credentials)
- Role-based access control and user session management
- Detailed audit logs for accountability

#### ⚙️ Native PowerDNS API Integration

Seamless communication with the **PowerDNS Authoritative Server API** — instant zone synchronization and zero downtime updates.

#### 💻 Built with Modern Web Tech

Designed with **Nuxt 4** and **Nuxt UI 4**, ensuring a fast, minimal, and responsive user experience with dark mode support.

#### 🤖 Intelligent Automation

Automatically manages:

- Zone propagation and synchronization
- Record updates and statistics collection
- Server status and uptime reporting

#### 🌐 Scalable & Flexible

Ready for **multi-server environments** — perfect for distributed networks, ISPs, and enterprise infrastructures.

---

## 🔍 Technical Overview

### ⚙️ Backend

- **Gin (Go)** — ultra-fast REST API framework.
- **MongoDB** — persistent storage for users, zones, and logs.
- **PowerDNS Authoritative API** — enables real-time DNS management and monitoring.
- Includes: JWT authentication, access control, and audit tracking.

### 🖥️ Frontend

- **Nuxt 4 + Nuxt UI 4** — modern Vue 3 ecosystem for fluid UX.
- Intuitive dashboards for zones, users, logs, and server statistics.

### 🌍 DNS Architecture

- Central control layer for PowerDNS authoritative instances.
- Real-time updates across **primary** and **secondary** servers.
- Supports reverse zones, dynamic record types, and DNSSEC-ready setups.

---

## 🧭 Application Pages Overview

SanchezDNS provides a clean, modular interface — each page designed for precision and simplicity in DNS management.

### 🌐 **Zones**

The heart of SanchezDNS — manage your **DNS zones** and all related records (A, AAAA, MX, TXT, SRV, etc.).  
Easily create, edit, and delete zones with full PowerDNS API synchronization.

<p align="center">
  <img src="https://github.com/rafinhacuri/SanchezDNS/blob/main/public/system.png"/>
</p>

### 📊 **Statistics**

Real-time overview of your DNS infrastructure, including:

- Number of zones, records, and users
- Server uptime, latency, and QPS (queries per second)
- TCP/UDP query distribution
- Active status indicators for connected servers

### ⚙️ **Configuration**

Customize system behavior and connection parameters:

- Define supported record types for your current connection
- Manage connection credentials (host, server ID, and API key)
- Adjust operational preferences without restarting services

### 👥 **Users**

(Administrator only)  
Manage user access within the system:

- Add or remove users with defined roles
- Control permissions and visibility across zones and records
- Perfect for teams managing shared DNS environments

### 🪵 **Logs**

(Administrator only)  
Centralized audit history tracking every action performed in the system:

- Zone creation, record changes, and deletions
- Login activity and API interactions
- Time-stamped records for compliance and transparency

### 🔗 **Connections**

(Administrator only)  
View and manage all PowerDNS server connections:

- Configure multiple authoritative servers
- Securely store API keys with AES-256 encryption
- Monitor connection health and latency in real-time

---

## ✅ Quick Start Checklist

### 🔧 Step 1 — Install Prerequisites

- [ ] Install **Docker**
- [ ] Install **Docker Compose**  
      👉 Official guide: [Get Docker](https://docs.docker.com/get-started/get-docker/)

---

### 📦 Step 2 — Get the `docker-compose.yaml`

Choose one of the options below to download the configuration file:

<details>
<summary>🔽 Using curl</summary>

```bash
curl -L -o docker-compose.yaml https://raw.githubusercontent.com/rafinhacuri/sanchezdns/main/docker-compose.yaml
```

</details>

<details>
<summary>🔽 Using wget</summary>

```bash
wget -O docker-compose.yaml https://raw.githubusercontent.com/rafinhacuri/sanchezdns/main/docker-compose.yaml
```

</details>

Alternatively, copy it directly from the [example file](https://github.com/rafinhacuri/sanchezdns/blob/main/docker-compose.yaml).

---

### 📝 Step 3 — Configure Environment

Place your `.env` file in the project root directory. The following environment variables are required for proper operation:

```
SITE_URL="https://dns.example.com"
MONGO_URL="mongodb://mongo.example.com:27017"
MONGO_USERNAME="mongouser"
MONGO_PASSWORD="exp"
MONGO_DB_NAME="sanchezdns"
JWT_SECRET="exmp"
CRYPT_KEY="exp"
```

> **Note:**
>
> - The `JWT_SECRET` must be **at least 32 characters long**.
> - The `CRYPT_KEY` must be a **base64-encoded string that decodes to 32 bytes (AES‑256)**.
>
> ✅ To generate a valid encryption key run:
>
> ```bash
> openssl rand -base64 32
> ```

---

### 🚀 Step 4 — Launch Services

Run the following commands:

```bash
docker compose pull
docker compose up -d --force-recreate
```

The backend API and frontend UI are preconfigured and will start automatically via Docker.

---

### 🔍 Step 5 — Verify Installation

Check running containers:

```bash
docker compose ps
```

If all services show as `Up`, you’re ready! 🎉

---

## 🤝 Contribution

Contributions, issues, and feature requests are welcome!  
Feel free to check [issues page](https://github.com/rafinhacuri/sanchezdns/issues) and submit pull requests.

---

## 📜 License

> Licensed under the [MIT License](https://github.com/rafinhacuri/sanchezdns/blob/main/LICENSE)  
> © 2025 [Rafael Curi Leonardo](https://github.com/rafinhacuri)

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
