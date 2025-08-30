# Suggested Development Commands

## Root Level (All Projects)
- **Development**: `pnpm dev` - Starts all projects in parallel
- **Build**: `pnpm build` - Builds all projects
- **Lint**: `pnpm lint` - Lints all projects
- **Lint Fix**: `pnpm lint:fix` - Auto-fixes lint issues across all projects
- **Test**: `pnpm test` - Runs tests for all projects
- **Clean**: `pnpm clean` - Cleans all projects
- **Install**: `pnpm install` - Install dependencies for all workspaces

## Individual Project Commands
- **Backend only**: `pnpm dev:backend`
- **Frontend only**: `pnpm dev:frontend`
- **V2 only**: `pnpm dev:v2`
- **Build individual**: `pnpm build:backend`, `pnpm build:frontend`, `pnpm build:v2`

## Backend Specific (okusuri-backend/)
Navigate to okusuri-backend/ directory first:
- **Development with hot reload**: `make dev` - Uses air for hot reload
- **Build**: `make build` - Outputs to ./bin/server
- **Run**: `make run` - Builds then runs the server
- **Test**: `make test` - Runs all Go tests with verbose output
- **Clean**: `make clean` - Removes ./bin and ./tmp directories
- **Install dependencies**: `make install-deps` - Downloads Go modules

## Frontend Projects
### okusuri-frontend
Navigate to okusuri-frontend/ directory:
- **Dev**: `npm run dev` or `pnpm dev` - Next.js with Turbopack and HTTPS
- **Build**: `npm run build` or `pnpm build`
- **Lint**: `biome lint ./src`
- **Lint Fix**: `biome lint ./src --fix`

### okusuri-v2
Navigate to okusuri-v2/ directory:
- **Dev**: `npm run dev` or `pnpm dev` - Vite dev server
- **Build**: `npm run build` or `pnpm build` - TypeScript build + Vite build
- **Lint**: `biome lint ./src`
- **Lint Fix**: `biome lint ./src --fix`
- **Format**: `biome format ./src --write`
- **Check**: `biome check ./src` - Combined lint + format check

## Essential System Commands (Darwin/macOS)
- **Git**: `git status`, `git add`, `git commit`, `git push`
- **File operations**: `ls`, `cd`, `find`, `grep` (prefer `rg` ripgrep if available)
- **Process management**: `ps aux`, `kill`, `pkill`
- **Network**: `lsof -i :8080` (check port usage), `netstat -an`