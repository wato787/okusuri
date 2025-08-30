# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Architecture

Okusuri (お薬管理アプリケーション) is a medication tracking application with multiple components:

- **`backend`** - Go REST API backend
- **`frontend`** - Next.js 15 frontend 
- **`okusuri-v2`** - Vite + React frontend (lightweight PWA version)
- **`docs`** - Design and infrastructure documentation
- **`notification`** - Notification-related documentation
- **`infra`** - Infrastructure-related documentation

## Development Commands

This project uses [Task](https://taskfile.dev/) for project management. Install Task first:

```bash
# macOS (Homebrew)
brew install go-task
```

### Primary Development Commands
- **Install dependencies**: `task install:all` (installs all project dependencies)
- **Frontend development**: `task dev:frontend` (Next.js with hot reload)
- **V2 development**: `task dev:v2` (Vite with hot reload) 
- **Backend development**: `task dev:backend` (Go server with hot reload)

### Individual Install Commands
- **Frontend dependencies**: `task install:frontend`
- **V2 dependencies**: `task install:v2` (corrected from original docs)
- **Backend dependencies**: `task install:backend`

### Build Commands
- **Build frontend**: `task build:frontend`
- **Build V2**: `task build:v2`
- **Build backend**: `task build:backend`
- **Build Lambda API**: `task build:lambda:api` (Docker container)
- **Build Lambda notification**: `task build:lambda:notification` (ZIP package)
- **Build all Lambda**: `task build:lambda`

### Production Commands
- **Start production frontend**: `task start:frontend`
- **Preview V2 build**: `task preview:v2`
- **Docker build backend**: `task docker:build:backend`
- **Docker run backend**: `task docker:run:backend`

### Linting & Formatting
- **Lint frontend**: `task lint:frontend`, `task lint:fix:frontend`
- **Lint V2**: `task lint:v2`, `task lint:fix:v2`

### Testing & Utilities  
- **Test backend**: `task test:backend`
- **Test with coverage**: `task test:coverage:backend`
- **Clean builds**: `task clean`
- **Clean backend only**: `task clean:backend`
- **Project status**: `task status`
- **List all tasks**: `task --list`

### Backend Specific (backend/ directory)
- **Development with Makefile**: `make dev` (uses air for hot reload)
- **Direct Go run**: `go run cmd/server/main.go`
- **Tests**: `go test ./...`

## Technology Stack & Architecture

### Backend (backend/)
- **Language**: Go 1.24
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL with GORM ORM
- **Architecture**: Clean architecture with handler/service/repository layers
- **Authentication**: JWT-based with Better Auth integration
- **Key Models**: MedicationLog (tracks intake + bleeding), NotificationSetting, User
- **API Prefix**: All endpoints under `/api`
- **Entry Point**: `cmd/server/main.go`
- **Lambda Support**: `cmd/lambda/` for AWS Lambda deployment

### Frontend (frontend/) - Next.js
- **Framework**: Next.js 15 with App Router
- **Authentication**: Better Auth + Google OAuth
- **Database**: Direct PostgreSQL connection
- **Styling**: Tailwind CSS + Radix UI
- **Features**: PWA capabilities, push notifications
- **Dev Server**: Uses Turbopack with experimental HTTPS

### V2 Frontend (okusuri-v2) - Vite + React
- **Framework**: Vite + React 19 + TypeScript
- **State Management**: TanStack Query (React Query) - MANDATORY for all API calls
- **Authentication**: Better Auth context
- **Styling**: Tailwind CSS + Radix UI with class-variance-authority
- **PWA**: Custom Service Worker (no Workbox)
- **Pages**: Home, Calendar, Settings only (lightweight design)

## Key Development Rules

### Language & Communication
- **Primary Language**: Japanese for all communication, code comments, and PR descriptions when working in this repository
- **Response Language**: Always respond in Japanese unless specifically asked to use English
- **Commit Messages**: Write in Japanese with descriptive purpose

### okusuri-v2 Specific Architecture Rules
- **API Communication**: MUST use TanStack Query for ALL API calls, never direct fetch
- **State Management**: TanStack Query + React Context, NO global state libraries
- **Error Handling**: Unified error handling through apiClient
- **Authentication**: Better Auth with automatic token management
- **PWA**: Custom lightweight Service Worker, NO Workbox

### Code Quality Standards
- **TypeScript**: Strict typing required, no `any` types
- **Backend**: Run `gofmt -l .` before commits, use `go mod tidy` after dependency changes
- **Frontend**: Use Biome for linting and formatting
- **Testing**: Use testify framework for Go backend tests
- **Coverage**: Use `task test:coverage:backend` to generate coverage reports

### Environment Setup
- **Node.js**: ≥18.0.0
- **pnpm**: ≥8.0.0 (package manager for frontend projects)
- **Go**: 1.24+
- **Task**: Task runner for project management (`brew install go-task`)
- **Database**: PostgreSQL (requires DATABASE_URL env var for backend)

### V2 Environment Variables (Required)
```env
VITE_API_URL=http://localhost:8080
VITE_VAPID_PUBLIC_KEY=your_vapid_public_key
```

## Database & Migrations
- Backend automatically runs migrations on startup
- Uses GORM for ORM operations
- Database connection configured in `pkg/config/database.go`

## Authentication Flow
- Google OAuth integration via Better Auth
- JWT tokens stored in localStorage (frontend) / context management
- Protected routes require authentication middleware
- User ID extraction from JWT for backend operations

## PWA Features (V2)
- Custom Service Worker for offline functionality
- Web Push notifications with VAPID keys
- Install prompts and PWA capabilities
- Lightweight caching strategy

## Development Workflow Best Practices
- Create feature branches from main
- Run linting and tests before commits: `task lint:frontend lint:v2 test:backend`
- Use Task commands for consistent development workflow
- Backend hot reload with air (via `task dev:backend` or `make dev`), frontend with Vite/Next.js HMR
- Use `task status` to check overall project health
- For backend: Use `make dev` in backend/ directory for air hot reloading
- For Lambda deployment: Use `task build:lambda` to build both API container and notification ZIP

## Project Structure Overview
```
okusuri/
├── backend/           # Go REST API with clean architecture
├── frontend/          # Next.js 15 app with App Router
├── okusuri-v2/        # Lightweight React PWA
├── docs/              # Design and infrastructure docs
├── notification/      # Notification documentation
├── infra/             # Infrastructure documentation
└── Taskfile.yml       # Task automation configuration
```