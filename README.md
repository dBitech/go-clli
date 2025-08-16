# CLLI

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.19-007d9c)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/dbitech/go-clli.svg)](https://pkg.go.dev/github.com/dbitech/go-clli)
[![Go Tests](https://github.com/dBitech/go-clli/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/dBitech/go-clli/actions/workflows/go.yml?query=branch%3Amain)
[![Markdown Lint](https://github.com/dBitech/go-clli/actions/workflows/markdown.yml/badge.svg?branch=main)](https://github.com/dBitech/go-clli/actions/workflows/markdown.yml?query=branch%3Amain)

A comprehensive Go package for parsing, validating, and extracting information from CLLI (Common Language Location Identifier) codes used in telecommunications.

## Features

- 🔍 **Parse CLLI codes** - Extract place, region, network site, and entity components
- 🌍 **Geographic resolution** - Convert to city names, state/province codes, and country information
- 📋 **Entity type identification** - Resolve switching equipment types from entity codes
- ✅ **Comprehensive validation** - Strict and relaxed parsing modes with detailed error reporting
- 🚀 **High performance** - Microsecond parsing times with minimal memory allocation
- 🔒 **Thread-safe** - All operations are safe for concurrent use
- 📚 **Complete documentation** - Extensive API documentation and examples

## What is CLLI?

Common Language Location Identifier (CLLI) codes are used in the telecommunications industry to identify specific locations, buildings, and equipment. Originally developed by Bell System, they follow the format defined in Bell System Practices Section 795-100-100.

A typical CLLI looks like: `MPLSMNMSDS1`

- **MPLS** - Place code (Minneapolis)
- **MN** - Region code (Minnesota)  
- **MS** - Network site code
- **DS1** - Entity code (switching equipment type)

## Quick Start

### Installation

```bash
go get github.com/dbitech/go-clli
```

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/dbitech/go-clli/pkg/clli"
)

func main() {
    // Parse a CLLI code
    c, err := clli.Parse("MPLSMNMSDS1")
    if err != nil {
        log.Fatal(err)
    }

    // Get geographic information
    fmt.Printf("Location: %s, %s, %s\n", 
        c.CityName(), c.StateName(), c.CountryName())
    // Output: Minneapolis, Minnesota, United States

    // Access CLLI components
    fmt.Printf("Place: %s\n", c.Place)           // MPLS
    fmt.Printf("Region: %s\n", c.Region)         // MN
    fmt.Printf("Network Site: %s\n", c.NetworkSite) // MS
    fmt.Printf("Entity Code: %s\n", c.EntityCode)   // DS1

    // Get entity type description
    fmt.Printf("Entity Type: %s\n", c.EntityType())
}
```

## Supported CLLI Types

### Entity CLLIs (8-11 characters)

Standard building locations with telecommunications equipment:

```go
c, _ := clli.Parse("MPLSMNMSDS1")  // Minneapolis, MN - DS1 switching equipment
```

### Non-Building Location CLLIs (8-12 characters)

Remote locations without traditional buildings:

```go
c, _ := clli.Parse("MPLSMNB1234")  // Minneapolis, MN - Location B, ID 1234
```

### Customer Location CLLIs (8-12 characters)

Customer-specific locations:

```go
c, _ := clli.Parse("MPLSMN1A234")  // Minneapolis, MN - Customer 1, ID A234
```

## Advanced Usage

### Parsing Options

```go
// Strict parsing (default) - validates all patterns
c, err := clli.Parse("MPLSMNMSDS1")

// Relaxed parsing - allows partial or non-standard formats
opts := clli.ParseOptions{Strict: false}
c, err := clli.ParseWithOptions("MPLS", opts)
```

### Error Handling

```go
c, err := clli.Parse("INVALID")
if err != nil {
    switch {
    case errors.Is(err, clli.ErrInvalidCLLI):
        fmt.Println("General CLLI format error")
    case errors.Is(err, clli.ErrInvalidPlace):
        fmt.Println("Invalid place code")
    case errors.Is(err, clli.ErrInvalidRegion):
        fmt.Println("Invalid region code")
    default:
        fmt.Printf("Parsing error: %v\n", err)
    }
}
```

### Pattern Validation

```go
// Check CLLI type before parsing
if clli.IsEntityCLLI("MPLSMNMSDS1") {
    fmt.Println("This is an entity CLLI")
}

if clli.IsNonBuildingCLLI("MPLSMNB1234") {
    fmt.Println("This is a non-building location CLLI")
}

// Validate individual components
if err := clli.ValidatePlace("MPLS", true); err != nil {
    fmt.Printf("Invalid place: %v\n", err)
}
```

### Geographic Information

```go
c, _ := clli.Parse("MPLSMNMSDS1")

// ISO codes
fmt.Printf("Country Code: %s\n", c.CountryCode()) // US
fmt.Printf("State Code: %s\n", c.StateCode())     // MN

// Full names  
fmt.Printf("Country: %s\n", c.CountryName())      // United States
fmt.Printf("State: %s\n", c.StateName())          // Minnesota
fmt.Printf("City: %s\n", c.CityName())            // Minneapolis
```

## Documentation

- **[API Documentation](docs/API.md)** - Complete API reference
- **[Specifications](docs/SPECIFICATIONS.md)** - Technical specifications and patterns
- **[Examples](examples/)** - Usage examples and code samples

## Test Framework Summary

This project is developed TDD-first. The test suite covers parsing/classification, validators (Bell tables B–E), pattern matchers, geographic helpers, integration flows, and concurrency/memory. See TEST_FRAMEWORK_SUMMARY.md for a concise overview, including public API and error-handling behavior.

### Running tests

Run package tests:

```bash
go test ./pkg/clli -v
```

Run all tests:

```bash
go test ./...
```

## CLLI Structure Reference

```text
Position: 1234567890A
Format:   PPPPRRSSXXX
```

| Component | Positions | Description | Example |
|-----------|-----------|-------------|---------|
| Place     | 1-4       | City/location abbreviation | MPLS |
| Region    | 5-6       | State/province code | MN |
| Network Site | 7-8    | Site identifier (optional) | MS |
| Entity/Location | 9-11 | Equipment or location code | DS1 |

## Entity Code Patterns

The package supports entity codes defined in Bell System Practices Tables B-E:

### Switching Entities (Table B)

- `MG1`, `SG2`, `DS1` - Various switching equipment types
- `2CB`, `3GT` - Trunk and gateway equipment  
- `RS5`, `CTX` - Remote switching and control

### Non-Switching Entities (Table E)

- `F23`, `AA3` - Facility and auxiliary equipment
- `Q45` - Special equipment designations

### Location Codes

- `B1234` - Building locations
- `1A234` - Customer locations

## Performance

The CLLI package is optimized for high-performance parsing:

- **Sub-millisecond parsing** - Typical CLLIs parse in microseconds
- **Low memory footprint** - Minimal allocations and string reuse
- **Concurrent safe** - All operations support concurrent access
- **Lazy loading** - Data files loaded only when needed

## Error Types

```go
var (
    ErrInvalidCLLI     = errors.New("invalid CLLI format")
    ErrInvalidPlace    = errors.New("invalid place code")  
    ErrInvalidRegion   = errors.New("invalid region code")
    ErrInvalidSite     = errors.New("invalid network site code")
    ErrInvalidEntity   = errors.New("invalid entity code")
    ErrInvalidLocation = errors.New("invalid location code")
)
```

## Data Sources

The package includes data for:

- **Region mappings** - CLLI region codes to ISO 3166 countries/states
- **City names** - Place code to city name mappings  
- **Entity types** - Entity code to equipment type descriptions
- **Location types** - Location code descriptions

## Requirements

- Go 1.19 or later
- No external runtime dependencies
- Optional: Enhanced geographic data packages

## Contributing

Contributions are welcome! Please see our contributing guidelines for details on:

- Reporting bugs
- Suggesting enhancements  
- Submitting pull requests
- Code style and testing requirements

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

This Go implementation is based on the Ruby gem [`steventwheeler/clli`](https://github.com/steventwheeler/clli) and follows the Bell System Practices Section 795-100-100 specification for CLLI codes.

## Support

- 📖 Check the [API Documentation](docs/API.md) for detailed usage information
- 🐛 Report bugs via [GitHub Issues](https://github.com/dbitech/go-clli/issues)  
- 💬 Ask questions in [GitHub Discussions](https://github.com/dbitech/go-clli/discussions)
- 📧 Contact the maintainers for security issues

---

**Example Projects Using CLLI:**

- Telecommunications network management systems
- Location-based service applications  
- Equipment inventory and tracking systems
- Geographic data processing pipelines
