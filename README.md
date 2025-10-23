<p align="center">
  <img src="https://github.com/rafinhacuri/SanchezDNS/blob/main/public/logo.png" alt="SanchezDNS Logo" width="180">
</p>


<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=venom&height=300&color=2FC851&text=SanchezDNS&section=header&fontColor=ffffff"/>
</p>

---

SanchezDNS is a modern, secure, and scalable DNS management platform designed to simplify authoritative DNS administration. Built with **Gin** (Go) backend, **MongoDB** for data persistence, **PowerDNS Authoritative API** integration, and a sleek **Nuxt 4** frontend UI with **Nuxt UI 4**, SanchezDNS empowers operators to manage DNS zones, records, users, and audit logs efficiently.

---

## 💡 Why Choose SanchezDNS

SanchezDNS is built for professionals who need both simplicity and power in DNS management.  
Here’s why it stands out:

- **Unified Management:** Manage all DNS zones, records, and users from a single, intuitive interface.  
- **Secure by Design:** AES‑256 encryption for credentials, role‑based access, and full audit logging.  
- **PowerDNS Integration:** Native integration with PowerDNS Authoritative API ensures real‑time synchronization and reliability.  
- **Modern Frontend:** Built with Nuxt 4 + Nuxt UI 4, offering a smooth and responsive experience on any device.  
- **Automated Operations:** Automatically handles updates, zone synchronization, and statistics collection.  
- **Scalable Infrastructure:** Ideal for organizations managing multiple domains or distributed DNS servers.  

Whether you are a hosting provider, enterprise network admin, or independent operator, SanchezDNS brings clarity, control, and performance to DNS management.

## 🔍 Technical Overview

### Backend
- **Gin**: High-performance Go web framework powering the RESTful API.
- **MongoDB**: Stores application data including users, zones, records, audit logs, and statistics.
- **PowerDNS Authoritative API**: Interfaces with PowerDNS to synchronize DNS zones and records, enabling real-time DNS management.
- Features include authentication (JWT-based), role-based access control, task management, and detailed audit logging.

### Frontend
- **Nuxt 4** with **Nuxt UI 4**: Vue 3-based modern UI framework delivering a responsive and intuitive user experience.
- Provides dashboards for DNS statistics, zone and record management, user administration, and system logs.

### DNS Architecture
- Acts as a centralized management layer for PowerDNS authoritative servers.
- Synchronizes DNS data via PowerDNS API, ensuring consistency and reliability.
- Supports multiple DNS zones, dynamic record updates, and real-time monitoring.

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