# Copilot Coding Agent Instructions for go-api-kit

## Project Overview
- **go-api-kit** is a modular Go framework for building secure, production-grade HTTP APIs, with a focus on out-of-the-box best practices and OWASP mitigations.
- Major components: `server/` (core server, handlers, middlewares), `config/` (env/config loading), `utils/` (utility types).
- The framework is designed for extensibility and composability, supporting SaaS, client/server, and general API products.

## Key Architectural Patterns
- **Handlers**: All HTTP logic is in `server/handlers/`. Use generic handlers (e.g., `ValidatingHandler`) for request validation and decoding.
- **Middlewares**: Pluggable, composable, and compatible with Negroni. See `server/middlewares/` for access log, security headers, content-type, and body size enforcement.
- **Validation**: Uses go-playground/validator. Validation errors are mapped to JSON fields using custom logic in `server/handlers/validation/`.
- **Error Responses**: Standardized error responses via `HttpError` and `ErrorDetails` (see `response-error-sender.go`). All error responses use `application/problem+json`.
- **Field Mapping Cache**: Reflection-based field-to-JSON mapping is cached for performance (see `struct-json-name-map.go`).

## Developer Workflows
- **Build**: Standard Go build (`go build ./...`).
- **Test**: Use `./test.sh` for full test suite (runs `go test ./...` with extra checks). Individual test files use testify/assert.
- **Lint**: Run `./lint.sh` or see `.github/workflows/lints.yml` for CI linting.
- **Release**: Automated via GitHub Actions (`.github/workflows/release.yml`).

## Project-Specific Conventions
- **ErrorDetails**: Always use `map[string]map[string]string` for field-level errors. Example:
  ```go
  handlers.SendBadRequest(w, handlers.ErrorDetails{
      "zip_code": {"pattern": "must match the required pattern"},
  })
  ```
- **Response Senders**: Use `SendJsonStruct` (or `SendStructAsJson`) for marshaling structs to JSON responses. Use `SendJson` for raw JSON bytes.
- **Validation**: Add validator tags to struct fields. Validation errors are auto-mapped to JSON field names.
- **Testing**: Use JSON-based assertions for error responses to avoid type mismatches.
- **Configuration**: All config is via environment variables (see README for full list).

## Integration Points
- **Negroni**: All middlewares are compatible with Negroni for flexible chaining.
- **go-playground/validator**: Used for struct validation.

## Key Files & Directories
- `server/handlers/response-error-sender.go`: Error response logic and types
- `server/handlers/validating_handler.go`: Generic request validation handler
- `server/middlewares/`: All middleware implementations
- `server/handlers/validation/`: Field mapping and validation helpers
- `config/`: Environment/config loading
- `README.md`: Usage, conventions, and examples

## Examples
- See `README.md` for up-to-date usage patterns and error response formats.
- Example error response:
  ```json
  {
    "status": 400,
    "title": "Bad Request",
    "details": {
      "zip_code": {"pattern": "must match the required pattern"}
    }
  }
  ```

---


## Commit Guidelines
- Use Conventional Commits for all commit messages. 
- The format is:
```
<type>[scope]: <description>
[optional body]
[optional footer(s)]
```

- The description should be concise and in the imperative mood (e.g., "add new endpoint" not "added new endpoint").
- Use the body to explain what and why vs. how, if necessary.
- Reference issues or PRs in the footer when relevant (e.g., "Closes #123").

### Common scopes:
- feat:     A new feature
- fix:      A bug fix
- docs:     Documentation only changes
- style:    Changes that do not affect the meaning of the code (white-space, formatting, etc)
- refactor: A code change that neither fixes a bug nor adds a feature
- perf:     A code change that improves performance
- test:     Adding or correcting tests
- chore:    Other changes that don't modify src or test files
