# SwipeLearn – Spaced Repetition App

SwipeLearn is a **production-oriented spaced repetition web platform** built to demonstrate **DevOps, Platform, and Cloud engineering skills**.

The goal of this project is not feature breadth, but to **design, deploy, and operate a real, stateful system** using Kubernetes, Infrastructure as Code, and CI/CD — from local development to a low-cost production environment.

---

## 🎯 Purpose

This project exists to answer one question clearly:

> *Can this engineer operate a real system in production?*

Focus areas:
- Kubernetes-based architecture
- Infrastructure as Code (Terraform)
- CI/CD pipelines
- Stateful workloads (PostgreSQL)
- Environment parity (local → production)
- Reliability and observability

---

## 🧠 Product Overview

Users review flashcards using a **spaced repetition algorithm (SM-2)** through a **swipe-based interface**.

Key characteristics:
- Time-based scheduling
- Data correctness matters
- Low-latency user interactions
- Persistent user progress

---

## 🏗️ Architecture

```
[ Browser (React) ]
        |
     [ Ingress ]
        |
      [ Go API ]
        |
   [ PostgreSQL ]
```

Design principles:
- Stateless application layer
- Stateful database with persistent volumes
- Declarative configuration
- Health checks and graceful shutdowns

---

## 🧩 Tech Stack

**Frontend**
- React (mobile-first)

**Backend**
- Go (REST API)
- Health & readiness probes
- Structured logging

**Database**
- PostgreSQL
- SM-2 scheduling logic

**Infrastructure**
- Kubernetes (k3d locally, k3s in production)
- Terraform (IaC)
- Helm (infra add-ons)
- Docker (multi-stage builds)

**CI/CD**
- GitHub Actions

**Observability (planned)**
- Prometheus
- Grafana

---

## 📂 Repository Structure

```
swipelearn/
├── app/              # Go backend
├── frontend/         # React frontend
├── infra/            # Terraform (local & prod)
├── k8s/              # Kubernetes manifests
├── .github/workflows # CI pipelines
├── docs/             # Architecture & decisions
└── README.md
```

---

## 🚀 Local Development

### Prerequisites
- Docker
- kubectl
- k3d
- Terraform
- Go
- Node.js

### Create local cluster
```bash
k3d cluster create dev   --servers 1   --agents 1   --port "8080:80@loadbalancer"
```

---

## 🧱 Infrastructure as Code

All infrastructure is managed with **Terraform**:
- Kubernetes namespaces
- Helm releases
- Base platform components

No manual kubectl drift. Environments are reproducible.

---

## 🔄 CI/CD Approach

- Tests and builds on every pull request
- Container images published via GitHub Actions
- Deployment automation evolves incrementally

Manual steps early on are **intentional and realistic**.

---

## 💰 Cost Model

| Environment | Cost |
|---|---|
| Local (k3d) | $0 |
| Production (VPS) | $5–10 / month |
| Optional second node | +$5 |

Designed to stay under **$15/month**.

---

## 📄 License

MIT License.

---

## 🧠 Why this project matters

This is not a demo app.

It is a **deliberate exercise in operating software**, covering:
- Deployment
- Configuration
- Failure handling
- State management
- Observability
- Cost awareness
