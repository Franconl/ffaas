# FFaaS (Feature Flags as a Service) — Go

FFaaS is a service that allows you to **manage feature flags** and supports **progressive delivery**: enable or disable features on the fly, perform canary releases, and customize user experiences without redeploying.

---

## ✨ Features (MVP)

- REST API for CRUD operations on feature flags (`/api/flags`).
- Health endpoints: `/healthz` and `/readyz`.
- Prometheus metrics endpoint: `/metrics`.
- SDK delivery of flags via HTTP (`/sdk/flags`).

---


