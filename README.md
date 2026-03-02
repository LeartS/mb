# mb — Metabase CLI

Work with [Metabase](https://www.metabase.com) from the command line. Create
dashboards, manage saved questions, run SQL queries, explore database schemas,
and more.

Designed for both human and LLM usage — structured JSON output, composable
commands, and a raw API escape hatch for anything not covered by dedicated
subcommands.

## Install

### From source (requires Go 1.26+)

```bash
go install github.com/LeartS/mb@latest
```

### Homebrew

```bash
brew install LeartS/tap/mb
```

### Binary releases

Download a prebuilt binary from the
[releases page](https://github.com/LeartS/mb/releases).

## Quick start

```bash
# Authenticate (API key — recommended)
mb auth login --host https://metabase.example.com --api-key mb_XXXX

# Authenticate (username + interactive password prompt)
mb auth login --host https://metabase.example.com --username admin@example.com

# Explore
mb database list
mb database schemas 1
mb database tables 1 public
mb search revenue

# Run an ad-hoc query
mb dataset query --database 1 --native-query "SELECT count(*) FROM orders"

# Create a saved question
mb card create \
  --name "Revenue by Month" \
  --database 1 \
  --native-query "SELECT date_trunc('month', created_at) AS month, SUM(total) FROM orders GROUP BY 1" \
  --display line \
  --collection 5

# Create a dashboard and add cards
mb dashboard create --name "Sales Overview" --collection 5
mb dashboard add-card 42 --from-json cards.json
```

## Commands

| Command | Description |
|---|---|
| `mb auth login` | Authenticate with a Metabase instance |
| `mb auth status` | Show current authentication status |
| `mb auth logout` | Remove stored credentials |
| `mb config set\|get\|list` | Manage CLI configuration |
| `mb card list\|get\|create\|update\|delete\|query` | Manage saved questions |
| `mb dashboard list\|get\|create\|update\|delete\|add-card\|copy` | Manage dashboards |
| `mb collection list\|tree\|get\|create\|items` | Manage collections |
| `mb database list\|get\|metadata\|schemas\|tables\|sync` | Explore databases |
| `mb dataset query` | Run ad-hoc SQL queries |
| `mb table get\|fields` | Inspect tables |
| `mb search <query>` | Search across all objects |
| `mb api <METHOD> <path> [body]` | Raw API escape hatch |

Run `mb <command> --help` for details on any command.

## Configuration

Config is stored in `~/.config/mb/config.json`. The following environment
variables override the config file:

| Variable | Description |
|---|---|
| `MB_HOST` | Metabase instance URL |
| `MB_API_KEY` | API key (takes precedence over session token) |
| `MB_SESSION_TOKEN` | Session token |
| `MB_CONFIG_DIR` | Override config directory |

## Authentication

Two methods are supported:

- **API key** (recommended): create one in Metabase Admin > Settings >
  Authentication > API Keys. Pass it with `--api-key` during login or set
  `MB_API_KEY`.
- **Session token**: authenticate with `--username` and an interactive password
  prompt. The session token is stored in the config file. Sessions expire after
  14 days by default.

## Output

All commands output JSON by default. Use `--json` to force JSON output even
when other formats might be available in the future.

## Raw API access

The `mb api` command is an escape hatch for any Metabase API endpoint not
covered by a dedicated command:

```bash
mb api GET /api/user/current
mb api POST /api/card/ '{"name":"test","dataset_query":...}'
mb api PUT /api/card/42 --input-file payload.json
mb api DELETE /api/card/42
```

## Dashboard creation workflow

A typical workflow for building a dashboard programmatically:

1. **Discover**: find the database ID and relevant tables.
   ```bash
   mb database list
   mb database tables 1 public
   mb table fields 42
   ```

2. **Create questions**: one `mb card create` per metric.
   ```bash
   mb card create --name "Total Revenue" --database 1 \
     --native-query "SELECT SUM(total) FROM orders" \
     --display scalar --collection 5
   ```

3. **Create the dashboard**:
   ```bash
   mb dashboard create --name "Sales Overview" --collection 5
   ```

4. **Lay out the cards**: build a JSON file with grid positions and apply it.
   The dashboard grid is 24 columns wide.
   ```bash
   mb dashboard add-card <dashboard-id> --from-json layout.json
   ```

See `mb dashboard add-card --help` for the JSON format.

## Development

```bash
# Build
go build -o mb .

# Run tests
go test ./...

# Vet
go vet ./...

# Install locally
go install .

# Build with version info
go build -ldflags "-X github.com/LeartS/mb/cmd/root.version=0.1.0 -X github.com/LeartS/mb/cmd/root.commit=$(git rev-parse HEAD) -X github.com/LeartS/mb/cmd/root.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o mb .
```

## License

[MIT](LICENSE)
