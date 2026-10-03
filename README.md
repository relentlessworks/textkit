# textkit

Agentic-first text analysis and manipulation service. Count words/characters/lines/sentences, wrap text, number lines, deduplicate, sort, trim, pad, reverse, head/tail, find/replace, extract lines, squeeze blank lines, grep, join lines, and get text statistics. Plain text API, agent-driven, single Go binary.

## Quick Start

```bash
# Build
make build

# Run
./textkit

# Or with custom address
./textkit -addr :3000

# Use it
curl -X POST localhost:8080/stats -d "Hello world. This is a test."
curl -X POST "localhost:8080/wrap?width=20" -d "The quick brown fox jumps over the lazy dog"
curl -X POST "localhost:8080/sort?desc=true" -d "banana
apple
cherry"
curl -X POST "localhost:8080/head?n=3" -d "$(cat file.txt)"
curl -X POST "localhost:8080/replace?old=foo&new=bar" -d "foo bar foo"
```

## API Reference

All endpoints accept POST with text in the request body. Parameters are passed as query strings.

| Endpoint | Description | Parameters |
|----------|-------------|------------|
| `POST /stats` | Text statistics | — |
| `POST /count` | Count words/chars/bytes/lines | `type` |
| `POST /wrap` | Wrap text to width | `width` |
| `POST /number` | Prepend line numbers | `start` |
| `POST /dedup` | Remove duplicate lines | — |
| `POST /sort` | Sort lines | `desc` |
| `POST /trim` | Trim whitespace | `mode` |
| `POST /pad` | Pad text to width | `width`, `align`, `fill` |
| `POST /reverse` | Reverse characters | — |
| `POST /reverse-lines` | Reverse line order | — |
| `POST /head` | First N lines | `n` |
| `POST /tail` | Last N lines | `n` |
| `POST /extract` | Extract line range | `start`, `end` |
| `POST /find` | Find lines containing string | `q` |
| `POST /replace` | Replace text | `old`, `new` |
| `POST /squeeze` | Collapse blank lines | — |
| `POST /grep` | Filter matching lines | `q`, `invert` |
| `POST /unique` | Remove duplicates (alias) | — |
| `POST /join` | Join lines with separator | `sep` |
| `POST /mcp` | MCP JSON-RPC 2.0 | — |
| `GET /help` | Operating manual | — |

### Response Format

- **Plain text** by default (one labeled, grepable line per record)
- **JSON** via `Accept: application/json` header or `?format=json` query param
- **Errors**: `error: message | hint: what to do next`

### MCP

The service speaks Model Context Protocol at `POST /mcp` with 18 tools.

## Configuration

| Flag | Env | Default | Description |
|------|-----|---------|-------------|
| `-addr` | `TEXTKIT_ADDR` | `:8080` | Listen address |
| `-no-auth` | `TEXTKIT_NO_AUTH` | `false` | Disable auth (dev mode) |

## Build

```bash
make build    # CGO_ENABLED=0, single static binary
make test     # go test -race
make vet      # go vet
```

## License

MIT
