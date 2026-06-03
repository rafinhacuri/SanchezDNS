# 🧾 Logs

The **Logs** page in **SanchezDNS** provides a complete audit trail of actions performed within the selected DNS connection.

## 🔍 Overview

Logs record every significant action taken by users on a specific DNS connection, ensuring transparency and accountability.

Each entry includes:

- **Username** — who performed the action
- **Action** — what was done (e.g., created a zone, added a record)
- **Details** — additional information about the event
- **Timestamp** — when the action occurred

## ⚙️ Behavior

- Logs are automatically generated for all changes made through the system.
- They are **connection‑specific**, meaning you only see logs related to the currently selected DNS server.
- Data is displayed in real time and stored securely in MongoDB for persistence.

## 🔒 Access Control

Only **administrators** can access the Logs page. Regular users cannot view or modify logs.

---

This feature ensures you can always trace what happened, who did it, and when — keeping your DNS operations auditable and secure.
