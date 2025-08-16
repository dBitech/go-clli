# CLLI-API Test Framework Summary

## Overview

This repository follows a test-driven approach for the Go CLLI-API. As of 2025-08-16, the library is fully implemented and all tests pass.

## Test Suite Layout

- Module: go.mod / go.sum
- Source: `pkg/clli/clli.go`
- Tests (under `pkg/clli/`):

  - `clli_test.go` — core types, parse errors, basic parsing, and benchmarks
  - `validation_test.go` — place/region/network-site/entity validation and edge cases
  - `patterns_test.go` — IsEntityCLLI / IsNonBuildingCLLI / IsCustomerCLLI and boundaries
  - `geographic_test.go` — US/CA resolution, edge cases, and consistency checks
  - `integration_test.go` — end-to-end parsing, options, error handling, concurrency, memory, workflow

Note: `debug_test.go` at the repo root is for ad‑hoc checks and not part of library tests.

## Coverage Areas

### Core parsing and classification

- Entity, Non‑Building, and Customer CLLIs
- Special cases: 8‑char Non‑Building (PPPPRRNN) and 15‑char Customer (PPPPRRNNXXXXXXX)
- Strict by default; case normalization and whitespace trimming enabled by default

### Validation

- Place: 4 uppercase letters
- Region: US states and Canadian provinces/territories (2‑letter codes)
- Network site: digits‑only for Entity; alphanumeric allowed for Non‑Building/Customer
- Entity code: strict patterns aligned to Bell tables B–E (DS/RT/SW/MS/XC, numeric/T/GT/RS/X?X, Table C/D/E variants)

### Pattern matchers

- `IsEntityCLLI`, `IsNonBuildingCLLI`, `IsCustomerCLLI` with mutual exclusivity and boundary tests

### Geographic helpers

- Country/State/City lookups for common US/CA codes; case/normalization and partial‑data handling

### Integration, performance, and robustness

- ParseWithOptions permutations, MustParse panic behavior
- Concurrency and memory tests across large input sets
- String() and type representation consistency

## Public API (current)

Types and constants

```go
type CLLIType int

const (
  CLLITypeUnknown CLLIType = iota
  CLLITypeEntity
  CLLITypeNonBuilding
  CLLITypeCustomer
)

// Back‑compat aliases used by tests
const (
  EntityCLLI      = CLLITypeEntity
  NonBuildingCLLI = CLLITypeNonBuilding
  CustomerCLLI    = CLLITypeCustomer
)

type CLLI struct {
  Original     string
  Place        string
  Region       string
  NetworkSite  string
  EntityCode   string
  LocationCode string
  LocationID   string
  CustomerCode string
  CustomerID   string
}

type ParseOptions struct {
  Strict           bool // default true; alias StrictValidation supported
  StrictValidation bool // alias of Strict
  NormalizeCase    bool // default true
  TrimWhitespace   bool // default true
}
```

Core functions

- `Parse(input string) (*CLLI, error)`
- `ParseWithOptions(input string, opts *ParseOptions) (*CLLI, error)`
- `MustParse(input string) *CLLI`

Instance methods

- `String() string`, `Type() CLLIType`, `IsValid() bool`
- `IsEntityCLLI() bool`, `IsNonBuildingCLLI() bool`, `IsCustomerCLLI() bool`
- `ValidatePlace() bool`, `ValidateRegion() bool`, `ValidateNetworkSite() bool`, `ValidateEntityCode() bool`
- `CountryCode() string`, `CountryName() string`, `StateCode() string`, `StateName() string`, `CityName() string`
- Convenience: `EntityType() string`, `LocationType() string`

Package‑level validation helpers

- `ValidatePlace(place string, strict bool) error`
- `ValidateRegion(region string, strict bool) error`
- `ValidateNetworkSite(site string, strict bool) error`
- `ValidateEntityCode(code string, strict bool) error`

## Error Handling

- Internal parse failures use a `ParseError` containing Input, Position, Field, and Err.
- `ParseError.Error()` formats as: "parse error at position N in field FIELD: REASON" (no input string) for unit assertions.
- Public errors from `Parse`/`ParseWithOptions` wrap the `ParseError` with the original input using Go error wrapping; callers see "INPUT: parse error ...".
- `MustParse` panics with: `MustParse failed for input "INPUT": ERROR`.

## Status (2025‑08‑16)

- All tests in `pkg/clli` pass, including integration, concurrency, and memory tests.
- Defaults: `Strict=true`, `NormalizeCase=true`, `TrimWhitespace=true`.

## Try It

### Package tests

```powershell
go test ./pkg/clli -v
```

### Full repo tests

```powershell
go test ./...
```

## Next Steps

- Broaden geographic mappings and document data sources
- Expand entity code tables with additional references
- Add CI workflow and badges; publish tagged releases
- Enhance README with examples and usage recipes
