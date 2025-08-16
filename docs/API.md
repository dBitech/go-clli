# CLLI Go Package API Documentation

## Overview

The CLLI package provides a complete Go implementation for parsing and working with Common Language Location Identifier (CLLI) codes used in telecommunications.

## Installation

```bash
go get github.com/dbitech/go-clli
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/dbitech/go-clli/pkg/clli"
)

func main() {
    // Parse a CLLI code
    c, err := clli.Parse("MPLSMNMSDS1")
    if err != nil {
        panic(err)
    }

    // Access geographic information
    location := fmt.Sprintf("%s, %s, %s", 
        c.CityName(), c.StateName(), c.CountryName())
    fmt.Println(location) // Output: Minneapolis, Minnesota, United States

    // Access CLLI components
    fmt.Printf("Place: %s\n", c.Place)           // MPLS
    fmt.Printf("Region: %s\n", c.Region)         // MN
    fmt.Printf("Site: %s\n", c.NetworkSite)     // MS
    fmt.Printf("Entity: %s\n", c.EntityCode)    // DS1
    fmt.Printf("Type: %s\n", c.EntityType())    // Entity description
}
```

## Core Types

### CLLI

The main CLLI structure represents a parsed Common Language Location Identifier.

```go
type CLLI struct {
    Original    string // The original CLLI string
    Place       string // 4-character place abbreviation
    Region      string // 2-character region code
    NetworkSite string // 2-character network site code (optional)
    EntityCode  string // 3-character entity code (optional)
    
    // Non-building location fields (mutually exclusive with entity)
    LocationCode string // 1-character location code (optional)
    LocationID   string // 4-character location ID (optional)
    
    // Customer location fields (mutually exclusive with entity)
    CustomerCode string // 1-character customer code (optional)
    CustomerID   string // 4-character customer ID (optional)
}
```

### CLLIType

Enumeration representing the type of CLLI.

```go
type CLLIType int

const (
    CLLITypeUnknown CLLIType = iota
    CLLITypeEntity           // Standard entity CLLI (8-11 chars)
    CLLITypeNonBuilding     // Non-building location (8-12 chars)
    CLLITypeCustomer        // Customer location (8-12 chars)
)

func (t CLLIType) String() string
```

### ParseOptions

Configuration options for parsing CLLIs.

```go
type ParseOptions struct {
    Strict bool // Enable strict pattern validation (default: true)
}
```

## Constructor Functions

### Parse

Creates a new CLLI instance from a string with default strict validation.

```go
func Parse(clli string) (*CLLI, error)
```

**Parameters:**

- `clli` - The CLLI string to parse

**Returns:**

- `*CLLI` - Parsed CLLI instance
- `error` - Error if parsing fails

**Example:**

```go
c, err := clli.Parse("MPLSMNMSDS1")
if err != nil {
    log.Fatal(err)
}
```

### ParseWithOptions

Creates a new CLLI instance with custom parsing options.

```go
func ParseWithOptions(clli string, opts ParseOptions) (*CLLI, error)
```

**Parameters:**

- `clli` - The CLLI string to parse
- `opts` - Parsing options

**Returns:**

- `*CLLI` - Parsed CLLI instance
- `error` - Error if parsing fails

**Example:**

```go
c, err := clli.ParseWithOptions("MPLSMN", ParseOptions{Strict: false})
if err != nil {
    log.Fatal(err)
}
```

### MustParse

Creates a new CLLI instance, panicking on parse error. Use only when input is guaranteed to be valid.

```go
func MustParse(clli string) *CLLI
```

**Parameters:**

- `clli` - The CLLI string to parse

**Returns:**

- `*CLLI` - Parsed CLLI instance

**Panics:** If parsing fails

**Example:**

```go
c := clli.MustParse("MPLSMNMSDS1")
```

## Core Methods

### IsValid

Returns true if the CLLI is properly formatted according to specification.

```go
func (c *CLLI) IsValid() bool
```

**Returns:**

- `bool` - True if CLLI is valid

### Type

Returns the type of CLLI (entity, non-building, or customer).

```go
func (c *CLLI) Type() CLLIType
```

**Returns:**

- `CLLIType` - The CLLI type

### String

Returns the original CLLI string.

```go
func (c *CLLI) String() string
```

**Returns:**

- `string` - Original CLLI string

## Geographic Methods

### CountryCode

Returns the ISO 3166 two-character country code for the CLLI location.

```go
func (c *CLLI) CountryCode() string
```

**Returns:**

- `string` - ISO 3166 country code (e.g., "US", "CA")

### CountryName

Returns the full country name for the CLLI location.

```go
func (c *CLLI) CountryName() string
```

**Returns:**

- `string` - Full country name (e.g., "United States", "Canada")

### StateCode

Returns the ISO 3166 state or province code for the CLLI location.

```go
func (c *CLLI) StateCode() string
```

**Returns:**

- `string` - State/province code (e.g., "MN", "QC") or empty if unknown

### StateName

Returns the full state or province name for the CLLI location.

```go
func (c *CLLI) StateName() string
```

**Returns:**

- `string` - Full state/province name (e.g., "Minnesota", "Quebec") or empty if unknown

### CityName

Returns the city name for the CLLI location if known.

```go
func (c *CLLI) CityName() string
```

**Returns:**

- `string` - City name (e.g., "Minneapolis") or empty if unknown

## Entity and Location Type Methods

### EntityType

Returns a description of the entity type for entity CLLIs.

```go
func (c *CLLI) EntityType() string
```

**Returns:**

- `string` - Entity type description or empty if not an entity CLLI

### LocationType

Returns a description of the location type for non-building or customer CLLIs.

```go
func (c *CLLI) LocationType() string
```

**Returns:**

- `string` - Location type description or empty if not a location CLLI

## Validation Functions

### ValidatePlace

Validates the place component of a CLLI.

```go
func ValidatePlace(place string, strict bool) error
```

**Parameters:**

- `place` - Place code to validate (4 characters)
- `strict` - Enable strict pattern validation

**Returns:**

- `error` - Error if validation fails

### ValidateRegion

Validates the region component of a CLLI.

```go
func ValidateRegion(region string, strict bool) error
```

**Parameters:**

- `region` - Region code to validate (2 characters)
- `strict` - Enable strict pattern validation

**Returns:**

- `error` - Error if validation fails

### ValidateNetworkSite

Validates the network site component of a CLLI.

```go
func ValidateNetworkSite(site string, strict bool) error
```

**Parameters:**

- `site` - Network site code to validate (2 characters)
- `strict` - Enable strict pattern validation

**Returns:**

- `error` - Error if validation fails

### ValidateEntityCode

Validates the entity code component of a CLLI.

```go
func ValidateEntityCode(code string, strict bool) error
```

**Parameters:**

- `code` - Entity code to validate (3 characters)
- `strict` - Enable strict pattern validation

**Returns:**

- `error` - Error if validation fails

## Pattern Matching Functions

### IsEntityCLLI

Returns true if the string matches entity CLLI patterns.

```go
func IsEntityCLLI(clli string) bool
```

**Parameters:**

- `clli` - CLLI string to test

**Returns:**

- `bool` - True if matches entity pattern

### IsNonBuildingCLLI

Returns true if the string matches non-building location CLLI patterns.

```go
func IsNonBuildingCLLI(clli string) bool
```

**Parameters:**

- `clli` - CLLI string to test

**Returns:**

- `bool` - True if matches non-building pattern

### IsCustomerCLLI

Returns true if the string matches customer location CLLI patterns.

```go
func IsCustomerCLLI(clli string) bool
```

**Parameters:**

- `clli` - CLLI string to test

**Returns:**

- `bool` - True if matches customer pattern

## Error Types

The package defines specific error types for different validation scenarios:

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

## Usage Examples

### Basic Parsing

```go
// Parse different types of CLLIs
entity, _ := clli.Parse("MPLSMNMSDS1")     // Entity CLLI
nonBuilding, _ := clli.Parse("MPLSMNB1234") // Non-building CLLI
customer, _ := clli.Parse("MPLSMN1A234")    // Customer CLLI

fmt.Printf("Entity type: %s\n", entity.Type())
fmt.Printf("Non-building type: %s\n", nonBuilding.Type())
fmt.Printf("Customer type: %s\n", customer.Type())
```

### Geographic Information

```go
c, _ := clli.Parse("MPLSMNMSDS1")

// Get location details
fmt.Printf("City: %s\n", c.CityName())         // Minneapolis
fmt.Printf("State: %s\n", c.StateName())       // Minnesota
fmt.Printf("State Code: %s\n", c.StateCode())  // MN
fmt.Printf("Country: %s\n", c.CountryName())   // United States
fmt.Printf("Country Code: %s\n", c.CountryCode()) // US
```

### Entity Information

```go
c, _ := clli.Parse("MPLSMNMSDS1")

if c.Type() == clli.CLLITypeEntity {
    fmt.Printf("Entity Code: %s\n", c.EntityCode)
    fmt.Printf("Entity Type: %s\n", c.EntityType())
}
```

### Validation

```go
// Validate individual components
if err := clli.ValidatePlace("MPLS", true); err != nil {
    fmt.Printf("Invalid place: %v\n", err)
}

if err := clli.ValidateRegion("MN", true); err != nil {
    fmt.Printf("Invalid region: %v\n", err)
}

// Check pattern types
if clli.IsEntityCLLI("MPLSMNMSDS1") {
    fmt.Println("This is an entity CLLI")
}
```

### Error Handling

```go
c, err := clli.Parse("INVALID")
if err != nil {
    switch {
    case errors.Is(err, clli.ErrInvalidCLLI):
        fmt.Println("General CLLI format error")
    case errors.Is(err, clli.ErrInvalidPlace):
        fmt.Println("Place code is invalid")
    case errors.Is(err, clli.ErrInvalidRegion):
        fmt.Println("Region code is invalid")
    default:
        fmt.Printf("Other error: %v\n", err)
    }
}
```

### Relaxed Parsing

```go
// Parse with relaxed validation for partial or malformed CLLIs
opts := clli.ParseOptions{Strict: false}
c, err := clli.ParseWithOptions("MPLS", opts)
if err == nil {
    fmt.Printf("Place: %s\n", c.Place)
    fmt.Printf("Valid: %t\n", c.IsValid())
}
```

## Thread Safety

All parsing and validation functions are thread-safe and can be called concurrently from multiple goroutines. The CLLI struct is immutable after creation.

## Performance

- Parsing operations complete in microseconds for typical CLLIs
- Memory allocation is minimized through string reuse
- Data files are loaded lazily and cached for performance
- No external network calls required for basic parsing
