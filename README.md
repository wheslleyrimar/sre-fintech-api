# SRE Fintech API

SRE-style API focused on **fintech**: accounts, tariff adjustments, and reports. It forwards requests to the [fintech-api-failures](../fintech-api-failures) backend, which simulates failures for resilience and chaos-engineering exercises.

## Requirements

- **Go 1.21** or higher
- **K6** (for validation tests): `brew install k6`
- **jq** (optional, for scripting): `brew install jq`
- **netcat** (for validation script port checks): usually pre-installed on macOS

## Quick start

### 1. Build binaries (recommended)

From this project root (the `sre` directory), run:

```bash
./install.sh
```

This clones [fintech-api-failures](https://github.com/wheslleyrimar/fintech-api-failures) from the remote repo, builds it, and compiles the SRE app from this repo. **All binaries are placed in `./bin/`** at this project root:

- `./bin/fintech-api-failures` — backend (listens on port 8080 by default)
- `./bin/sre` — SRE API (use `PORT=8081` to avoid clashing)

The install script does **not** start any process.

### 2. Run manually (optional)

Start the backend, then the SRE API:

```bash
./bin/fintech-api-failures
```

In another terminal:

```bash
BACKEND_URL=http://localhost:8080 PORT=8081 ./bin/sre
```

The SRE API listens on **http://localhost:8081**. Use `BACKEND_URL` to point at wherever the fintech-api-failures backend is running.

## Install script

From the project root:

```bash
./install.sh
```

- Clones **fintech-api-failures** from [GitHub](https://github.com/wheslleyrimar/fintech-api-failures), builds it, and writes the binary to **`./bin/fintech-api-failures`**.
- Builds the **SRE** app and writes the binary to **`./bin/sre`**.
- Does **not** start any process. Run the binaries manually or use the validation script (see below).

## Environment variables

| Variable       | Description                          | Default            |
|----------------|--------------------------------------|--------------------|
| `BACKEND_URL`  | Base URL of the fintech-api-failures | `http://localhost:8080`  |
| `SRE_BASE_URL` | Base URL of this SRE API (callbacks)  | `http://localhost:8080`  |
| `PORT`         | Port for this SRE API                | `8080`             |

## API endpoints (v1)

All routes are under `/v1`.

| Method | Path                                   | Description                    |
|--------|----------------------------------------|--------------------------------|
| GET    | `/v1/search`                           | Search accounts (term optional) |
| GET    | `/v1/report`                           | Fintech summary report         |
| GET    | `/v1/accounts/{id}`                    | Get account by ID              |
| GET    | `/v1/accounts/{id}/tariff-adjustments` | Tariff adjustment history      |
| POST   | `/v1/accounts/{id}/tariff-adjustments` | Create tariff adjustment       |
| POST   | `/v1/accounts/notifications`           | Callback for adjustment result |

## Validation

The validation script runs k6 tests against the SRE API and the backend. It **uses the binaries in `./bin/`** produced by `./install.sh`.

### Run validation

From the project root:

```bash
./validation.sh <app_name> <group_number:1|2|3|4|5|6|7|8> local <case:1|2|3|4>
```

**Arguments:**

| Argument       | Values                 | Description                                  |
|----------------|------------------------|----------------------------------------------|
| `app_name`     | any string             | Application name                             |
| `group_number` | 1, 2, 3, 4, 5, 6, 7, 8 | Group number                                 |
| `env`          | `local`                | Environment (only `local` is supported)      |
| `case`         | 1, 2, 3, 4             | Test case (uses `validations/case_N.js`)     |

**Example:**

```bash
./validation.sh sre-finance 1 local 1
```

**How it works:**

- Runs **`./bin/fintech-api-failures`** on port **8080** (or `$FINTECH_BINARY` if set).
- Runs **`./bin/sre`** on port **8081** with `BACKEND_URL` and `PORT` set.
- Executes the chosen k6 script from `validations/`, then tears down both processes.

**Prerequisite:** Run **`./install.sh`** first so that `./bin/fintech-api-failures` and `./bin/sre` exist.

## Project structure

```
.
├── bin/                  # Binaries (created by ./install.sh): fintech-api-failures, sre
├── cmd/api/              # Entrypoint
├── internal/
│   ├── domain/           # Account, TariffAdjustmentRequest, Report
│   ├── http/             # Chi handlers (accounts, report, search)
│   ├── httpClient/       # HTTP client for backend calls
│   ├── integrations/    # AccountsApi, SearchEngine, AdjustmentFlowProcessor
│   ├── usecases/         # Account, Report, Search services
│   └── utils/            # Helpers
├── validations/          # K6 scripts (case_1.js, ...)
├── install.sh            # Builds binaries into ./bin/
├── validation.sh         # Runs validation using ./bin/ binaries
└── README.md
```

## License

Educational / demonstration use.
