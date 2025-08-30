# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Architecture

Okusuri (お薬管理アプリケーション) is a monorepo containing a medication tracking application with three main components:

- **`okusuri-backend`** - Go REST API backend
- **`okusuri-frontend`** - Next.js 15 frontend 
- **`okusuri-v2`** - Vite + React frontend (lightweight PWA version)

## Development Commands

### Root Level (All Projects)
- **Development**: `pnpm dev` (starts all projects in parallel)
- **Build**: `pnpm build` (builds all projects)
- **Lint**: `pnpm lint` (lints all projects)
- **Lint Fix**: `pnpm lint:fix` (auto-fixes lint issues)
- **Test**: `pnpm test` (runs tests for all projects)
- **Clean**: `pnpm clean` (cleans all projects)

### Individual Projects
- **Backend only**: `pnpm dev:backend`
- **Frontend only**: `pnpm dev:frontend` 
- **V2 only**: `pnpm dev:v2`
- **Build individual**: `pnpm build:backend`, `pnpm build:frontend`, `pnpm build:v2`

### Backend Specific (okusuri-backend/)
- **Development with hot reload**: `make dev` (uses air)
- **Build**: `make build` (outputs to ./bin/server)
- **Test**: `make test`
- **Clean**: `make clean`

### Frontend Projects
- **okusuri-frontend**: Uses Biome for linting (`biome lint ./src`)
- **okusuri-v2**: Uses Biome for linting (`biome lint ./src`, `biome format ./src --write`)

## Technology Stack & Architecture

### Backend (okusuri-backend)
- **Language**: Go 1.24
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL with GORM ORM
- **Architecture**: Clean architecture with handler/service/repository layers
- **Authentication**: JWT-based with Better Auth integration
- **Key Models**: MedicationLog (tracks intake + bleeding), NotificationSetting, User
- **API Prefix**: All endpoints under `/api`

### Frontend (okusuri-frontend) - Next.js
- **Framework**: Next.js 15 with App Router
- **Authentication**: Better Auth + Google OAuth
- **Database**: Direct PostgreSQL connection
- **Styling**: Tailwind CSS + Radix UI
- **Features**: PWA capabilities, push notifications
- **Dev Server**: `next dev --turbopack --experimental-https`

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

### Environment Setup
- **Node.js**: ≥18.0.0
- **pnpm**: ≥8.0.0 (required package manager)
- **Go**: 1.24+
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
- Run linting and tests before commits
- Use monorepo workspace commands for cross-project operations
- Backend hot reload with air, frontend with Vite/Next.js HMR