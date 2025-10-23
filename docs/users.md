# 👥 User Management

SanchezDNS includes a simple but secure user system to control access to your DNS infrastructure.

## 🧩 Overview

- The **first registered user** in the system automatically becomes the **Administrator**.  
- The Administrator has full access to every part of the platform, including:
  - Viewing and managing all DNS connections  
  - Creating, editing, or deleting zones and records  
  - Managing other users  
  - Viewing logs and audit activity  

All other users who register afterward will **not** have access to any DNS connections until the Administrator explicitly grants it.

---

## 🔐 Access Control

- New users can freely create accounts, but by default they have **no permissions** to view or modify any zones or servers.  
- The Administrator must manually associate them with one or more **connections** using the **Users page** (available only to administrators).  
- Once added to a connection, a user can:
  - View and manage zones within that connection  
  - Create and update DNS records (depending on their assigned access level)  

---

## ⚙️ Administrator Panel

Accessible only to the first registered user, the **Users page** allows the admin to:
- View a list of all registered users  
- Add or remove user access to specific DNS connections  
- Deactivate or delete accounts if needed  

This ensures tight control over who can interact with your authoritative DNS servers.

---

## ⚠️ Security Notes

- Only assign access to trusted users — each connection includes sensitive PowerDNS API credentials.  
- Rotate API keys and passwords periodically.  
- Always use strong passwords when creating accounts.  

---

With this approach, SanchezDNS ensures that your DNS infrastructure remains secure and organized, with full transparency over user permissions and actions.

---
