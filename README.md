# simulation-catalogue

## Overview

Source code for physics simulation catalogue application used to run Fortran-powered physics simulations. Simulations are presented in a catalogue-style interface, where users can browse, run, and visualize results. Simulations are all written in Fortran, and ran using Golang.

Simulation metadata and executables are stored in a PostgreSQL database, and runs are implemented using a job queue. All components are containerized using Docker, and are deployed using Kubernetes and Helm. Infrastructure is managed using Terraform.

Note that the simulations themselves are not included in this repository. See links below for more information.

### Links

- [Simulations GitHub Repository](https://github.com/psauerborn/simulations)
- [Simulation Catalogue Application](https://simulation-catalogue.s31-software.com)

## Table of Contents

- [Overview](#simulation-catalogue)
- [Project Structure](#project-structure)
- [Components](#components)
  - [API](#api)
  - [Web](#web)
  - [Database Migrations](#database-migrations)
  - [CLI Tool](#cli-tool)
  - [Infrastructure](#infrastructure)
- [Client System](#client-system)

## Project Structure

```
simulation-catalogue/
├── api/                    # Go REST API server
│   ├── controller.go       # HTTP route handlers
│   ├── db.go               # PostgreSQL database operations
│   ├── runner.go           # Simulation execution engine
│   ├── openapi.yaml        # OpenAPI 3.0 specification
│   ├── Makefile            # Build, test, and deploy tasks
│   └── Dockerfile          # Container image definition
├── web/                    # Quasar (Vue 3) frontend with SSR
│   ├── src/
│   │   ├── components/     # Reusable Vue components
│   │   ├── pages/          # Route page components
│   │   └── stores/         # Pinia state management
│   ├── Makefile            # Build and deploy tasks
│   └── Dockerfile          # Container image definition
├── alembic/                # Database migrations (Python/Alembic)
│   └── versions/           # Migration scripts
├── scripts/
│   └── simctl/             # CLI tool for simulation management (Go)
├── terraform/              # Infrastructure as Code
│   └── modules/            # Kubernetes/Helm deployments
│       ├── api/            # API service Helm chart
│       ├── web/            # Web service Helm chart
│       └── helm/           # Shared Helm configurations
├── docker-compose.yml      # Local PostgreSQL development database
├── Makefile                # Root-level tasks (secret scanning)
└── LICENSE                 # MIT License
```

## Components

### API

The REST API is written in Go and provides endpoints for:
- Browsing simulation metadata and parameters
- Running simulations with custom inputs
- Retrieving simulation outputs (JSON or ZIP)
- Admin operations (create, update, delete simulations)

The API runs simulations by fetching the relevant binary from the database, creating a temporary file system for the simulation, and running the binary with the provided input parameters. All simulation outputs are then stored as a zip file in the blob store.

Simulations are currently run within the actual API container on a separate set of go-routines, but will be moved to a separate container in the future and orchestrated using RabbitMQ.

See [`api/openapi.yaml`](api/openapi.yaml) for the full API specification.

### Web

The frontend is built with [Quasar Framework](https://quasar.dev/) (Vue 3) and features:
- Server-side rendering (SSR) for SEO
- LaTeX rendering for mathematical models
- Interactive plots for simulation results (Plotly.js)

### Database Migrations

Database schema is managed using [Alembic](https://alembic.sqlalchemy.org/). Migrations are located in `alembic/versions/`.

### Infrastructure

Deployment is managed via Terraform with Helm charts for Kubernetes:
- AWS ECR for container images
- Kubernetes (bare-metal) deployments and services
- Ingress configuration for public access

## Client System

The application uses a lightweight, anonymous session system to track simulation runs without requiring user authentication:

1. **Client Initialization**: On first visit, the frontend calls `POST /v1/public/client/init` to receive a unique Client ID (32-character hex string)
2. **Local Storage**: The Client ID is persisted in the browser's `localStorage`, allowing the session to survive page refreshes and browser restarts
3. **Request Authentication**: For simulation-related endpoints, the Client ID is sent via the `X-Client-Id` header
4. **Session Tracking**: The API validates the Client ID against the database and associates simulation runs with the client

This approach provides session continuity without the overhead of traditional authentication, while still allowing users to:
- Track their active simulation runs
- Retrieve outputs from completed simulations
- Maintain state across browser sessions

It also avoids the need to store or handle any personal data; users can use the application without creating an account.
