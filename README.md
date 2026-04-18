# envcmp

> CLI tool to diff and validate `.env` files across environments with secret masking

## Installation

```bash
go install github.com/yourusername/envcmp@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/yourusername/envcmp/releases).

## Usage

```bash
# Compare two .env files
envcmp diff .env.development .env.production

# Validate a .env file against a template
envcmp validate .env.example .env

# Diff with secrets masked
envcmp diff .env.staging .env.production --mask-secrets
```

**Example output:**

```
~ DB_HOST        dev.db.local  →  prod.db.example.com
+ SENTRY_DSN     [masked]
- DEBUG           true
```

### Flags

| Flag             | Description                          |
|------------------|--------------------------------------|
| `--mask-secrets` | Redact values for sensitive keys     |
| `--strict`       | Exit with non-zero code on any diff  |
| `--format`       | Output format: `text`, `json`        |

## Why envcmp?

- Catch missing or mismatched environment variables before deployment
- Safe to use in CI pipelines — secrets are never exposed in logs
- Zero dependencies, single binary

## License

MIT © [yourusername](https://github.com/yourusername)