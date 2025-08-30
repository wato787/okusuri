# Architecture Details

## Backend Architecture (okusuri-backend)

### Project Structure
```
okusuri-backend/
├── cmd/server/main.go          # Entry point, server setup
├── internal/                   # Core business logic
│   ├── dto/                   # Data Transfer Objects
│   ├── handler/               # HTTP request handlers
│   ├── middleware/            # Auth, CORS, logging
│   ├── model/                 # Database models (GORM)
│   ├── repository/            # Data access layer
│   ├── service/               # Business logic services
│   └── routes.go             # Router configuration
├── pkg/config/                # Database configuration
├── migrations/                # Database migrations
└── Makefile                   # Build and dev commands
```

### Key Models
- **MedicationLog**: Tracks medication intake with bleeding status and timestamps
- **NotificationSetting**: Manages user notification preferences per platform
- **User**: User authentication and profile information
- **Session**: User session management
- **Account**: User account linking for OAuth

### API Structure
- **Base Path**: `/api`
- **Health**: `/api/health` (public)
- **Notifications**: `/api/notification` (send), `/api/notification/setting` (settings, auth required)
- **Medication**: `/api/medication-log` (logging, auth required)
- **Status**: `/api/medication-status` (status, auth required)

### Authentication Flow
1. JWT-based authentication with Better Auth integration
2. Protected routes use auth middleware
3. User ID extracted from JWT tokens for data access
4. Google OAuth integration for user registration/login

## Frontend Architecture

### okusuri-frontend (Next.js 15)
- **App Router**: Modern Next.js routing system
- **Auth**: Better Auth with Google OAuth integration
- **Database**: Direct PostgreSQL connection for some operations
- **PWA**: Progressive Web App capabilities
- **Push Notifications**: Firebase/web push integration
- **Styling**: Tailwind CSS + Radix UI components

### okusuri-v2 (Vite + React 19)
- **Lightweight Design**: Only Home, Calendar, Settings pages
- **State Management**: TanStack Query (mandatory for all API calls)
- **Routing**: React Router DOM
- **PWA**: Custom Service Worker implementation
- **Auth Context**: Better Auth context management
- **API Client**: Unified API client with error handling

## Development Environment

### Prerequisites
- **Node.js**: ≥18.0.0
- **pnpm**: ≥8.0.0 (required package manager)
- **Go**: 1.24+
- **PostgreSQL**: Database server
- **Air**: Go hot reload tool (auto-installed via make dev)

### Environment Variables
#### Backend
- `DATABASE_URL`: PostgreSQL connection string (required)

#### okusuri-v2
- `VITE_API_URL`: Backend API URL (e.g., http://localhost:8080)
- `VITE_VAPID_PUBLIC_KEY`: VAPID public key for push notifications

### Development Workflow
1. **Monorepo Structure**: Uses pnpm workspaces
2. **Parallel Development**: `pnpm dev` starts all projects
3. **Hot Reload**: Air for backend, Vite/Next.js HMR for frontends
4. **Database**: Migrations run automatically on backend startup
5. **Cross-Project**: Changes can affect multiple components