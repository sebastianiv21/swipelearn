# PROJECT_CONTEXT.md
# SwipeLearn – DevOps Learning & Build Intent

This file exists to give any collaborator, reviewer, or AI agent clear context on **what is being built**, **why it is being built**, and **what skills are intentionally being practiced**.

This is not a product roadmap.  
It is an **engineering intent document**.

---

## 1. High-level Goal

Build and operate a **production-style web platform** that demonstrates real-world **DevOps / Platform Engineering** skills while staying within a low personal budget.

The project prioritizes:
- Operability over feature richness
- Infrastructure correctness over UI polish
- Realistic constraints over demos

---

## 2. What is being built

**SwipeLearn** is a spaced repetition web application with a swipe-based interface.

At a product level:
- Users create decks and cards
- Reviews follow the SM-2 spaced repetition algorithm
- Progress is persistent and time-dependent

At a systems level:
- A stateless backend API
- A stateful PostgreSQL database
- A frontend web client
- All deployed on Kubernetes

The application exists to support infrastructure and operational learning.

---

## 3. Target Architecture

```
[ Browser (React) ]
        |
     [ Ingress ]
        |
      [ Go API ]
        |
   [ PostgreSQL ]
```

Key properties:
- Stateless application layer
- Stateful database with persistent volumes
- Declarative configuration
- Local-to-production parity

---

## 4. Technology Choices (Intentional)

### Backend
- **Language:** Go
- Reason: strong DevOps/SRE market demand, fast builds, common infra tooling language

### Frontend
- **Framework:** React
- Mobile-first, minimal complexity

### Database
- **PostgreSQL**
- Chosen for realism, state management, and operational learning

### Containerization
- Docker (multi-stage builds, non-root images)

### Orchestration
- Kubernetes
  - Local: k3d (k3s in Docker)
  - Production: k3s on a VPS

### Infrastructure as Code
- Terraform
- Manages Kubernetes infrastructure and add-ons

### CI/CD
- GitHub Actions
- Build, test, image publishing, deployment evolution

### Observability (later)
- Prometheus
- Grafana

---

## 5. What this project is NOT

- Not a startup MVP
- Not optimized for monetization
- Not feature-complete learning software
- Not a framework showcase

This is a **systems and operations project**.

---

## 6. Learning Objectives

This project is explicitly designed to practice:

- Kubernetes fundamentals (Deployments, Services, Ingress, StatefulSets)
- Infrastructure as Code with Terraform
- CI/CD pipeline design
- Managing stateful workloads
- Environment parity (local vs production)
- Failure modes and recovery
- Cost-aware infrastructure decisions
- Production-style repo structure
- Operational thinking

---

## 7. Constraints

- Budget cap: **$0 locally, $15/month max in production**
- Solo developer (team simulated via process)
- Time horizon: months, not days

---

## 8. Development Philosophy

- Infrastructure before application features
- Small, reversible changes
- Commit early, commit clearly
- Document decisions
- Prefer boring, proven tech

---

## 9. How to work on this repo

Assumptions for collaborators or agents:
- Treat this as a production system
- Favor clarity over cleverness
- Avoid unnecessary abstractions
- Every tool choice should be justifiable

---

## 10. Current Status

- Repository initialized
- README written
- License chosen (MIT)

Next milestones:
- Local Kubernetes cluster (k3d)
- Terraform providers and infra skeleton
- Minimal backend service deployed

---

## 11. Success Criteria

This project is successful if the author can confidently discuss:
- How it is deployed
- How it fails
- How it is monitored
- How it is recovered
- How it scales
- How much it costs

Not how pretty the UI is.
