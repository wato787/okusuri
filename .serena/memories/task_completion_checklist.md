# Task Completion Checklist

## Before Committing Code Changes

### Backend (okusuri-backend)
1. **Format Code**: Run `gofmt -l .` and ensure no files are listed
   - If files are listed, run `gofmt -w .` to format them
2. **Dependencies**: Run `go mod tidy` if dependencies were added/removed
3. **Build**: Run `make build` to ensure code compiles successfully
4. **Test**: Run `make test` to ensure all tests pass
5. **Lint**: Check for any Go linting issues (golangci-lint if configured)

### Frontend Projects
#### okusuri-frontend
1. **Lint**: Run `biome lint ./src` and fix any issues
2. **Build**: Run `pnpm build` to ensure successful build
3. **Type Check**: Ensure TypeScript compilation succeeds

#### okusuri-v2
1. **Lint**: Run `biome lint ./src` and fix any issues
2. **Format**: Run `biome format ./src --write` for consistent formatting
3. **Check**: Run `biome check ./src` for combined checks
4. **Build**: Run `pnpm build` (includes TypeScript compilation)

### Root Level (Monorepo)
1. **Lint All**: Run `pnpm lint` to lint all projects
2. **Build All**: Run `pnpm build` to build all projects
3. **Test All**: Run `pnpm test` to run tests across all projects

## Quality Assurance Steps
1. **Manual Testing**: Test the changed functionality manually
2. **API Testing**: For backend changes, test API endpoints
3. **Cross-Platform**: For frontend changes, test on different devices/browsers
4. **Authentication**: Verify auth flows if authentication-related changes

## Pre-Commit Requirements
1. Code must be properly formatted according to project standards
2. All tests must pass
3. No linting errors or warnings
4. Code must build successfully
5. Changes must be tested (manual or automated)

## Commit Standards
1. Write commit messages in Japanese
2. Include descriptive purpose explaining the "why"
3. Link to issues when applicable (`Closes #XX`)
4. Include Claude Code signature:
   ```
   🤖 Generated with [Claude Code](https://claude.ai/code)
   
   Co-Authored-By: Claude <noreply@anthropic.com>
   ```

## Post-Commit Actions
1. Verify the commit was successful
2. Check that all changes are properly tracked
3. Push to remote repository if ready for collaboration
4. Create PR if changes are ready for review