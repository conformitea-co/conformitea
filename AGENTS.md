# Agent Guidelines for ConformiTea

## Build/Test Commands
- Frontend: `cd frontend && npm run dev` (development), `npm run build` (production)
- Go Backend: `air` (hot reload dev), `go build -o ./tmp/main .` (build), `go test ./...` (all tests)
- Single Go test: `go test -run TestName ./path/to/package`
- Lint: `cd frontend && npx eslint src/` (frontend), `gofmt -l .` (Go formatting check)

## Go Code Style
- Use Go modules with workspace (`go.work`), packages: app, cmd, domain, infrastructure, server
- Follow standard Go naming: PascalCase for exports, camelCase for private
- Struct receivers: Use pointer receivers for methods that modify state
- Error handling: Always check errors, return early on error
- Imports: Standard library first, then third-party, then local packages

## Frontend Code Style (React + TypeScript)
- Use TypeScript strict mode, React Router v7, Tailwind CSS, shadcn/ui components
- Import order: `@/types/*`, `@/hooks/*`, `@/stores/*`, `@/lib/*`, `@/components/*`, relative imports
- Prettier: 2 spaces, double quotes, semicolons, 110 char width
- Components: Use forwardRef for UI components, prefer function components
- File naming: kebab-case for files, PascalCase for components
- Props: Always type props interfaces, use `React.ComponentProps` for extending HTML elements

## Error Handling
- Go: Return errors as last value, use custom error types in `cerror` package
- Frontend: Use try/catch for async operations, display user-friendly error messages