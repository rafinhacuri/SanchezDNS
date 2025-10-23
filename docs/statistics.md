# 📊 Statistics

## Overview
This page provides real-time operational metrics from the connected PowerDNS Authoritative Server. It helps administrators monitor performance, uptime, and server activity, enabling effective management of the DNS infrastructure.

## Metrics Available
The main metrics displayed include:

- **Zones** – Total number of DNS zones managed by the server.  
- **Records** – Total number of resource records across all zones.  
- **Users** – Number of system users assigned to the connection.  
- **Uptime** – Time since the server started or was last restarted.  
- **QPS (Queries Per Second)** – Current query rate handled by the DNS server.  
- **UDP Queries / TCP Queries** – Breakdown of queries by protocol type.  
- **Latency P50 / P95** – Median and 95th percentile response latency (if available).  

## Data Source
All statistics are collected via the PowerDNS Authoritative API endpoint `/api/v1/servers/{server-id}/statistics`, ensuring accurate, real-time information is retrieved directly from the server.

## Usage
Statistics data refresh automatically when viewing the dashboard, providing instant insight into system health and workload distribution without manual intervention.

## Best Practices
- Monitor **QPS** and **latency** metrics to identify potential overload or configuration issues.  
- Use statistics to plan scaling efforts or detect abnormal query spikes early.  
- Correlate high query load periods with DNS cache settings and TTL configurations to optimize performance.  

> SanchezDNS provides clear visibility into your DNS infrastructure, helping you maintain performance, stability, and trust in every query.
