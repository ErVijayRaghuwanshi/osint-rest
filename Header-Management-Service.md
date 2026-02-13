# 🔐 Header Management Service (Gin + Kubernetes)

## 📌 Overview

This module provides a **dynamic, platform-aware HTTP header management system** for a Go (Gin)–based service running in Kubernetes.

It enables:

* Centralized header storage
* Platform-specific header configurations
* Multiple header sets per platform
* Round-robin rotation
* Zero code changes for header updates

This is particularly useful for **OSINT collectors**, **social platform integrations**, and **anti-bot–sensitive APIs** such as Snapchat, Instagram, and X (Twitter).

---

## 🧠 Problem Statement

External platforms often require:

* Frequent header updates (tokens, cookies, UA strings)
* Header rotation to avoid rate limits or bans
* Platform-specific header logic

Hardcoding headers leads to:

* Frequent rebuilds
* Operational risk
* Poor scalability

---

## ✅ Solution

Headers are:

* Stored as JSON in **persistent storage**
* Loaded dynamically by the Gin service
* Selected per platform using **round-robin rotation**
* Updated independently of application code

---

## 🏗️ Initial Architecture

```
┌─────────────────────────────┐
│ Kubernetes Cluster          │
│                             │
│  ┌──────────────┐           │
│  │ Header JSON  │◄──────────┼── Persistent Volume / ConfigMap
│  └──────────────┘           │
│          │                  │
│          ▼                  │
│  ┌───────────────────────┐  │
│  │ Gin API Service       │  │
│  │                       │  │
│  │  HeaderManager        │  │
│  │   ├── Load config     │  │
│  │   ├── Platform lookup │  │
│  │   ├── Round robin     │  │
│  │                       │  │
│  └─────────┬─────────────┘  │
│            ▼                │
│    External Platforms       │
│   (Snapchat / X / Instagram)│
└─────────────────────────────┘
```

---

## 📁 Header Configuration File

### File Location

Mounted inside the container via:

* **PersistentVolume** (recommended)
* or **ConfigMap** (read-only)

Example path:

```
/config/headers.json
```

---

## 📄 JSON Structure

```json
{
  "version": "1.0",
  "platforms": {
    "snapchat": {
      "enabled": true,
      "rotation": "round_robin",
      "headers": [
        {
          "id": "snap-1",
          "headers": {
            "User-Agent": "Mozilla/5.0 ...",
            "Authorization": "Bearer token1",
            "Accept": "application/json"
          }
        },
        {
          "id": "snap-2",
          "headers": {
            "User-Agent": "Mozilla/5.0 ...",
            "Authorization": "Bearer token2",
            "Accept": "application/json"
          }
        }
      ]
    },
    "instagram": {
      "enabled": true,
      "rotation": "round_robin",
      "headers": [
        {
          "id": "insta-1",
          "headers": {
            "User-Agent": "Mozilla/5.0 ...",
            "Cookie": "sessionid=..."
          }
        }
      ]
    }
  }
}
```

---

## 🧩 Core Components

### 1️⃣ HeaderManager

Responsible for:

* Loading JSON configuration
* Maintaining per-platform counters
* Selecting headers using rotation strategy

### 2️⃣ Platform Configuration

Each platform defines:

* Enable/disable flag
* Rotation strategy
* One or more header sets

### 3️⃣ Header Rotation

Current strategy:

* `round_robin`

Each API request fetches the **next header set** for the platform.

---

## 🔁 Request Flow

1. Incoming API request hits Gin service
2. Platform is identified (e.g. `snapchat`)
3. HeaderManager selects a header set
4. Headers are injected into outbound request
5. Request is sent to external platform

---

## 🚀 Kubernetes Integration

### Storage Options

#### Option 1: Persistent Volume (Recommended)

* Allows runtime updates
* Supports admin tooling / hot reload

#### Option 2: ConfigMap

* Simple and safe
* Requires pod restart for updates

---

## 🔒 Security Considerations

* Header files should **not be committed** to source control
* Mount secrets via:

  * Encrypted PV
  * Sealed Secrets (optional)
* Limit file access permissions inside container

---

## 🧪 Operational Benefits

* No redeployments for header updates
* Easy header revocation
* Platform-specific isolation
* Production-safe concurrency handling

---

## 🛣️ Future Enhancement Plan

### 🔄 1. Hot Reloading

* Watch file changes using `fsnotify`
* Reload headers without restarting pods

### 📊 2. Header Health & Scoring

* Track failures per header
* Auto-disable failing headers

```json
"health": {
  "failures": 3,
  "last_used": "2026-01-01T10:00:00Z"
}
```

---

### 🔁 3. Advanced Rotation Strategies

* Weighted round robin
* Random with cooldown
* Failure-aware selection

---

### 🌐 4. Proxy + Header Pairing

```json
{
  "headers": { ... },
  "proxy": "http://proxy-1"
}
```

---

### 🔐 5. Secret Manager Integration

* AWS Secrets Manager
* HashiCorp Vault
* GCP Secret Manager

---

### 🧠 6. Admin APIs

* `/admin/headers`
* `/admin/platforms`
* Enable/disable headers at runtime

---

### 📈 7. Metrics & Observability

* Prometheus counters
* Header usage stats
* Platform error rates

---

## 🎯 Design Principles

* **Config over code**
* **Platform isolation**
* **Operational flexibility**
* **Production-first mindset**

---

## 📌 Summary

This header management system provides a **scalable, extensible, and production-grade foundation** for interacting with sensitive third-party platforms while minimizing operational risk and deployment overhead.

---
