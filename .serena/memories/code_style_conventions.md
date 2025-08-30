# Code Style and Conventions

## Language & Communication Requirements
- **Primary Language**: Japanese for all communication, code comments, and PR descriptions when working in this repository
- **Response Language**: Always respond in Japanese unless specifically asked to use English
- **Commit Messages**: Write in Japanese with descriptive purpose

## Backend (Go) Conventions
- **Formatting**: Always run `gofmt -l .` before commits to ensure proper formatting
- **Dependencies**: Run `go mod tidy` after adding or removing dependencies
- **Testing**: Use testify framework (`github.com/stretchr/testify`) for tests
- **Architecture**: Clean architecture pattern - handler/service/repository layers
- **Naming**: Use Go standard naming conventions (PascalCase for exported, camelCase for unexported)
- **Error Handling**: Proper Go error handling with descriptive error messages

## Frontend (TypeScript/React) Conventions
- **Linting**: Biome for both frontend projects
- **Indentation**: Tab style (configured in biome.json)
- **Quotes**: Double quotes for JavaScript/TypeScript strings
- **Import Organization**: Enabled and automatic via Biome
- **TypeScript**: Strict typing required, no `any` types allowed
- **Unused Code**: Warn on unused imports and variables

### okusuri-v2 Specific Rules
- **API Calls**: MUST use TanStack Query for ALL API calls, never direct fetch
- **State Management**: TanStack Query + React Context, NO global state libraries
- **Error Handling**: Unified error handling through apiClient
- **Authentication**: Better Auth with automatic token management
- **PWA**: Custom lightweight Service Worker, NO Workbox

## File Structure Conventions
- Backend: Clean architecture with clear separation of concerns
- Frontend: Component-based architecture with clear feature separation
- Shared utilities and types in appropriate directories

## Testing Standards
- Backend: Go tests with testify framework, run `make test`
- Write tests for new features and bug fixes
- Maintain good test coverage

## Commit and PR Standards
- **Commit Messages**: Japanese, descriptive, explain the "why" not just "what"
- **PR Descriptions**: Japanese, include summary and test plan
- **Branch Names**: Descriptive feature branches from main
- **Size**: One feature per PR, keep changes focused and reviewable