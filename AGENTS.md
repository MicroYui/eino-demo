# Repository Guidelines

## Project Structure & Module Organization
- `main.go` boots the GoFrame HTTP server.
- `api/` holds GoFrame API contracts; generated interfaces live in `api/chat/chat.go` (do not edit).
- `internal/controller/` and `internal/logic/` contain request handlers and core business logic.
- `internal/ai/` contains agent pipelines, tools, and CLI entry points under `internal/ai/cmd/`.
- `utility/` hosts shared middleware and clients (Milvus, config helpers).
- `SuperBizAgentFrontend/` contains the static UI (`index.html`, `app.js`, `styles.css`).
- `manifest/` stores runtime config and Docker resources; `hack/` holds GoFrame build/codegen recipes; `docs/` is project documentation.

## Build, Test, and Development Commands
- `go run .` starts the API server on port 6872 (see `main.go`).
- `make -f hack/hack.mk build` builds the backend binary via GoFrame (`gf build -ew`).
- `make -f hack/hack.mk ctrl` regenerates controllers from `api/` contracts.
- `docker compose -f manifest/docker/docker-compose.yml up` starts Milvus and its dependencies for vector search.
- `cd SuperBizAgentFrontend && npm run start` (or `python3 -m http.server 8080`) serves the frontend.

## Coding Style & Naming Conventions
- Go: run `gofmt` (tabs) and keep package names lowercase; versioned handlers follow `chat_v1_*.go`.
- Generated files: avoid editing `api/chat/chat.go`; update `api/chat/v1/*.go` and rerun GoFrame codegen.
- CLI commands: follow the existing `*_cmd` naming under `internal/ai/cmd/`.
- Frontend: vanilla JS/HTML/CSS with 4-space indentation as in `SuperBizAgentFrontend/app.js`.

## Testing Guidelines
- No `*_test.go` files are currently present; add tests beside the package they cover.
- Run `go test ./...` for backend checks; use `curl http://localhost:6872/api/...` for endpoint smoke tests.
- The frontend has no automated tests; do a quick browser smoke test after changes.

## Commit & Pull Request Guidelines
- Git history is not included in this release; no commit convention was discoverable. Use short, imperative summaries (e.g., `feat: add upload validation`) and keep commits scoped.
- PRs should include a brief description, commands run, and note config changes. Include screenshots for `SuperBizAgentFrontend/` UI updates.

## Security & Configuration Tips
- `manifest/config/config.yaml` and `internal/ai/cmd/knowledge_cmd/config/config.yaml` include API keys and `file_dir` paths; keep secrets out of git and localize values.
- Set `file_dir` to your local docs directory before running the knowledge indexer (`internal/ai/cmd/knowledge_cmd`).
