# ⚙️ Setup Guide

Welcome to **SanchezDNS** — a powerful and elegant platform for managing PowerDNS authoritative servers.  
This guide will help you get the system running smoothly.

---

## 🧩 Prerequisites

Before you begin, make sure you have:

- **Docker & Docker Compose** installed on your system.  
- A **PowerDNS Authoritative Server** (version 4.7 or higher).  
- Access to a **MongoDB instance** (local or remote).  
- Basic knowledge of DNS records and PowerDNS API.

---

## 🧱 Environment Configuration

Create a `.env` file in the root directory with the following structure:

```bash
SITE_URL="https://dns.example.com"
MONGO_URL="mongodb://mongo.example.com:27017"
MONGO_USERNAME="mongouser"
MONGO_PASSWORD="essa senha é mt dificil de ser quebrada :o"
MONGO_DB_NAME="expo-go"

# Must have at least 32 characters
JWT_SECRET="essa senha é mt dificil de ser quebrada :o"

# This must be a base64-encoded string that decodes to 32 bytes (AES-256)
CRYPT_KEY="essa senha é mt dificil de ser quebrada :o meu deussss"
```

✅ To generate a valid encryption key, run:

```bash
openssl rand -base64 32
```

> ⚠️ Keep your environment file safe. Never share your `JWT_SECRET` or `CRYPT_KEY`.

---

## 🐳 Run with Docker

Simply run the following commands to build and start SanchezDNS:

```bash
docker compose build
docker compose up -d
```

The system will automatically start both backend and frontend containers.

Once running, open your browser and go to:

```
http://localhost:3000
```

---

## 🌍 Accessing the Interface

When you first open SanchezDNS, you'll be guided through:

1. Creating the **first admin user**  
2. Adding your **PowerDNS connection**  
3. Syncing zones and viewing server statistics  

The admin user automatically gains access to:
- Logs
- Connections
- System configuration tools

---

## 🧠 Notes

- The backend and API are **fully preconfigured** — no manual edits are needed.  
- DNS records are synced automatically with your PowerDNS server.  
- You can safely update via Docker without losing data (persistent volumes).  

> _SanchezDNS is designed to be ready out of the box — focus on your zones, not the setup._

---
