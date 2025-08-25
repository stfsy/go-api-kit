---
applyTo: "**/*.go"
---

# Libraries
- Use https://github.com/kelseyhightower/envconfig for environment configuration
- Use https://github.com/stfsy/go-api-kit for http handler configuration, endpoints, middlewares, sending http responses
- Use https://github.com/stfsy/go-api-key for API key creation and valiation
- Use https://github.com/stfsy/go-argon2id for secure hashing and hash verification
- Use https://github.com/stretchr/testify for test assertions

# Testing
- Creates test cases for all new functions, test happy paths and edge cases.

# Error Handling
- Handle all errors from functions that are called and return a wrapped error with a concise error message to the callee.
- For clarity, do not use single-line error handling, always separate the declaration of the error and its check

## Preferred Error Handling Style
```go
result, err := someMethod()
if err != nil {
    // handle error
}
```
- **Do NOT use:**
```go
var err error
var result SomeType
result, err = someMethod()
if err != nil {
    // handle error
}
```

# Code Style
- Use idiomatic Go style: follow effective Go and Go community conventions.
- Use gofmt for formatting and goimports for import management.
- Use clear, descriptive names for variables, functions, and types.
- Keep functions small and focused; prefer composition over inheritance.
- Add comments for exported functions, types, and complex logic.
- Avoid global variables; use dependency injection where possible.
- Group related code into packages; avoid circular dependencies.
- Use error wrapping and context for error handling.
- Prefer explicitness over cleverness; optimize for readability and maintainability.

# Context Management
- Only create a new context with `context.Background()` in exceptional situations or in test cases. In all other cases, use the parent context, e.g., the context of the http request. If beneficial, create a new child context to be able to cancel only the child request.