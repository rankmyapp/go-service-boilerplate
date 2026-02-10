# Export Microservice Strategy (MIDEV-6329)

## Goals
- Support export requests for `csv`, `jpeg`, and `pdf`.
- Support two source types: `chart` and `table`.
- Keep alignment with project architecture: `handler -> usecase -> repository/infra`.

## API Contract
- Endpoint: `POST /api/v1/exports`
- Request body fields:
  - `format`: `csv | jpeg | pdf`
  - `source_type`: `chart | table`
  - `mode`: `sync | async` (default `sync`)
  - `file_name`, `locale`, `timezone`, `delimiter` (optional metadata)
  - `payload`: data required by the selected strategy

### Sync mode
- Returns generated file as binary response.
- Required headers:
  - `Content-Type`
  - `Content-Disposition: attachment`

### Async mode
- Returns job envelope with status `queued`.
- Current state: only strategy contract is defined; job execution pipeline will be implemented in the next steps.

## Layered Design
- `internal/handlers/export_handler.go`
  - Input validation and HTTP mapping.
  - Maps domain/usecase errors to status codes.
- `internal/usecase/export_usecase.go`
  - Central orchestration and request validation.
  - Selects exporter strategy by `(format, source_type)`.
- `internal/repository/export_repository.go`
  - Contract for async job persistence.
- `models/export.go`
  - Export domain contracts.

## Strategy Registry
- Usecase stores strategies in map key `(format, source_type)`.
- Unsupported combinations return `ErrUnsupportedExportStrategy`.
- Implemented:
  - `pkg/export/csv`: `csv + chart` strategy.
  - `pkg/export/csv`: `csv + table` strategy.
  - `pkg/export/jpeg`: `jpeg + chart` strategy.
- Pending:
  - `pkg/export/pdf`

## Error Policy
- `400`: invalid request.
- `501`: strategy not implemented yet.
- `500`: unexpected internal failures.

## Next Deliverables
- MIDEV-6484: PDF strategy for table.
- MIDEV-6483: unit/integration/golden tests.
