# CLLI (Common Language Location Identifier) Go Package Specifications

## Overview

This Go package provides functionality to parse, validate, and extract information from CLLI (Common Language Location Identifier) codes used in telecommunications to identify locations. It is based on the Ruby gem `steventwheeler/clli` and follows the Bell System Practices Section 795-100-100 specification.

## Key Features

### 1. CLLI Parsing

- Parse CLLI strings into their component parts
- Extract place, region, network site, entity code, and location identifiers
- Support for both strict and relaxed parsing modes
- Comprehensive validation of CLLI format and structure

### 2. Geographic Information

- Convert CLLI region codes to ISO 3166 country codes
- Convert CLLI region codes to state/province codes
- Resolve city names from place codes and region codes
- Support for US and Canadian locations

### 3. Entity Type Resolution

- Identify switching entity types from entity codes
- Resolve entity code descriptions and types
- Support for multiple entity code patterns (Tables B-E)

### 4. Location Type Resolution

- Identify location types for non-building locations
- Resolve customer location codes
- Support for specialized location identifiers

## CLLI Structure

A CLLI code consists of up to 11 characters with the following structure:

```text
Position: 1234567890A
Format:   PPPPRRSSXXX
```

Where:

- **PPPP** (1-4): Place abbreviation (city identifier)
- **RR** (5-6): Region code (state/province)
- **SS** (7-8): Network site code
- **XXX** (9-11): Entity code, location code, or customer location identifier

## Supported CLLI Types

### 1. Entity CLLIs (8-11 characters)

- Standard building locations with equipment
- Network site + entity code combination
- Example: `MPLSMNMSDS1` (Minneapolis, MN, MS site, DS1 entity)

### 2. Non-Building Location CLLIs (8-12 characters)

- Remote locations without buildings
- Place + region + location code + 4-digit ID
- Example: `MPLSMNB1234` (Minneapolis, MN, location B, ID 1234)

### 3. Customer Location CLLIs (8-12 characters)

- Customer-specific locations
- Place + region + customer code + alphanumeric ID
- Example: `MPLSMN1A234` (Minneapolis, MN, customer 1, ID A234)

## Pattern Matching

### Character Classes

- **a**: All alphabetic characters (A-Z)
- **a1**: Alphabetic excluding B,D,I,O,T,U,W,Y
- **a2**: Alphabetic excluding G
- **n**: All numeric characters (0-9)
- **x**: All alphanumeric characters (A-Z, 0-9)
- **x1**: a1 + numeric characters
- **x2**: a2 + numeric characters

### Entity Code Patterns (Table B-E)

#### Switching Entities (Table B)

- Pattern: `(MG|SG|CG|DS|RL|PS|RP|CM|VS|OS|OL|[0-9]{2})[x1]`
- Pattern: `[CB0-9][0-9]T`
- Pattern: `[0-9]GT`
- Pattern: `Z[A-Z]Z`
- Pattern: `RS[0-9]`
- Pattern: `X[A-Z]X`
- Pattern: `CT[x1]`

#### Switchboard and Desk Entities (Table C)

- Pattern: `[0-9][CDBINQWMVROLPEUTZ0-9]B`

#### Miscellaneous Switching Entities (Table D)

- Pattern: `[0-9][AXCTWDEINPQ]D`
- Pattern: `[A-Z0-9][UM]D`

#### Non-Switching Entities (Table E)

- Pattern: `[FAEKMPSTW][x2][x1]`
- Pattern: `Q[0-9][0-9]`

## API Specification

### Core Types

```go
type CLLI struct {
    // Raw CLLI string
    Original string
    
    // Parsed components
    Place       string
    Region      string
    NetworkSite string
    EntityCode  string
    
    // Non-building location fields
    LocationCode string
    LocationID   string
    
    // Customer location fields
    CustomerCode string
    CustomerID   string
}
```

### Constructor Functions

```go
// Parse creates a new CLLI instance from a string with strict validation
func Parse(clli string) (*CLLI, error)

// ParseWithOptions creates a new CLLI with parsing options
func ParseWithOptions(clli string, opts ParseOptions) (*CLLI, error)

// MustParse creates a new CLLI instance, panicking on error
func MustParse(clli string) *CLLI
```

### Parse Options

```go
type ParseOptions struct {
    Strict bool // Enable strict pattern validation
}
```

### Core Methods

```go
// IsValid returns true if the CLLI is properly formatted
func (c *CLLI) IsValid() bool

// Type returns the type of CLLI (entity, non-building, customer)
func (c *CLLI) Type() CLLIType

// String returns the original CLLI string
func (c *CLLI) String() string
```

### Geographic Methods

```go
// CountryCode returns the ISO 3166 country code
func (c *CLLI) CountryCode() string

// CountryName returns the full country name
func (c *CLLI) CountryName() string

// StateCode returns the ISO 3166 state/province code
func (c *CLLI) StateCode() string

// StateName returns the full state/province name
func (c *CLLI) StateName() string

// CityName returns the city name if known
func (c *CLLI) CityName() string
```

### Entity/Location Type Methods

```go
// EntityType returns the entity type description
func (c *CLLI) EntityType() string

// LocationType returns the location type description
func (c *CLLI) LocationType() string
```

### Validation Functions

```go
// ValidatePlace validates the place component
func ValidatePlace(place string, strict bool) error

// ValidateRegion validates the region component
func ValidateRegion(region string, strict bool) error

// ValidateNetworkSite validates the network site component  
func ValidateNetworkSite(site string, strict bool) error

// ValidateEntityCode validates the entity code component
func ValidateEntityCode(code string, strict bool) error
```

### Pattern Matching Functions

```go
// IsEntityCLLI returns true if the CLLI matches entity patterns
func IsEntityCLLI(clli string) bool

// IsNonBuildingCLLI returns true if the CLLI matches non-building patterns
func IsNonBuildingCLLI(clli string) bool

// IsCustomerCLLI returns true if the CLLI matches customer patterns
func IsCustomerCLLI(clli string) bool
```

## Error Handling

The package defines specific error types for different validation failures:

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

The package requires data files for:

1. **Region Conversions**: Mapping CLLI region codes to ISO 3166 codes
2. **City Names**: Mapping place+region combinations to city names
3. **Entity Types**: Descriptions for entity code patterns
4. **Location Types**: Descriptions for location code patterns

## Performance Requirements

- Parse operations should complete in < 1ms for typical CLLIs
- Memory usage should be minimal with shared data structures
- Support concurrent access to parsing functions
- Lazy loading of data files

## Testing Requirements

- Unit tests for all parsing functions with valid/invalid inputs
- Integration tests with real CLLI data sets
- Benchmark tests for performance validation
- Property-based testing for pattern validation
- Test coverage > 90%

## Compatibility

- Go 1.19 or later
- No external dependencies for core functionality
- Optional dependencies for enhanced geographic data
