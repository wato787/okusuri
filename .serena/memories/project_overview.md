# Okusuri Project Overview

## Project Purpose
Okusuri (お薬管理アプリケーション) is a medication tracking application built as a monorepo with multiple frontend and backend components. The application helps users track their medication intake, manage notification settings, and monitor their medication status.

## Project Structure
This is a monorepo containing three main components:

- **`okusuri-backend`** - Go REST API backend (primary server)
- **`okusuri-frontend`** - Next.js 15 frontend application
- **`okusuri-v2`** - Vite + React lightweight PWA version

## Tech Stack

### Backend (okusuri-backend)
- **Language**: Go 1.24
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL with GORM ORM
- **Architecture**: Clean architecture with handler/service/repository layers
- **Authentication**: JWT-based with Better Auth integration
- **Hot Reload**: Air for development

### Frontend (okusuri-frontend) - Next.js
- **Framework**: Next.js 15 with App Router
- **Authentication**: Better Auth + Google OAuth
- **Database**: Direct PostgreSQL connection
- **Styling**: Tailwind CSS + Radix UI
- **Features**: PWA capabilities, push notifications
- **Dev Server**: Turbopack with experimental HTTPS

### V2 Frontend (okusuri-v2) - Vite + React
- **Framework**: Vite + React 19 + TypeScript
- **State Management**: TanStack Query (React Query) - MANDATORY for all API calls
- **Authentication**: Better Auth context
- **Styling**: Tailwind CSS + Radix UI with class-variance-authority
- **PWA**: Custom Service Worker (no Workbox)
- **Pages**: Home, Calendar, Settings only (lightweight design)

## Core Features
- Medication intake tracking with bleeding status
- User notification management per platform
- Authentication and user management
- PWA capabilities for mobile experience
- Push notification support