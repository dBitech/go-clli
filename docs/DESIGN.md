# CLLI Go Implementation Design Document

## Overview

This document outlines the design and architecture for implementing the CLLI (Common Language Location Identifier) Go package, based on the Ruby gem `steventwheeler/clli`.

## Architecture Overview

```text
pkg/clli/
├── clli.go              # Core CLLI struct and parsing logic
├── patterns.go          # Regex patterns for CLLI validation
├── validation.go        # Validation functions
├── geographic.go        # Geographic information resolution
├── entity_types.go      # Entity type resolution
├── location_types.go    # Location type resolution  
├── errors.go            # Error definitions
├── data/                # Embedded data files
│   ├── regions.yaml     # Region code mappings
│   ├── cities.yaml      # City name mappings
│   ├── entities.yaml    # Entity type descriptions
│   └── locations.yaml   # Location type descriptions
└── internal/
    ├── patterns/        # Pattern matching utilities
    ├── data/           # Data loading and caching
    └── utils/          # Utility functions
```

## Core Components

### 1. CLLI Struct

The main `CLLI` struct represents a parsed Common Language Location Identifier:

```go
type CLLI struct {
    Original    string // Original input string
    Place       string // 4-character place abbreviation
    Region      string // 2-character region code
    NetworkSite string // 2-character network site (optional)
    EntityCode  string // 3-character entity code (optional)
    
    // Non-building location fields
    LocationCode string // 1-character location code (optional)
    LocationID   string // 4-character location ID (optional)
    
    // Customer location fields  
    CustomerCode string // 1-character customer code (optional)
    CustomerID   string // 4-character customer ID (optional)
    
    // Parsed metadata
    cliType CLLIType // Determined CLLI type
    valid   bool     // Validation status
}
```

### 2. Pattern Matching System

The pattern matching system uses compiled regular expressions for efficient parsing:

```go
type PatternMatcher struct {
    entityPattern     *regexp.Regexp
    nonBuildingPattern *regexp.Regexp
    customerPattern   *regexp.Regexp
    
    // Component patterns
    placePattern      *regexp.Regexp
    regionPattern     *regexp.Regexp
    networkSitePattern *regexp.Regexp
    entityCodePattern *regexp.Regexp
}
```

### 3. Data Management

Data files are embedded in the binary using Go's `embed` package and loaded lazily:

```go
//go:embed data/*.yaml
var dataFS embed.FS

type DataManager struct {
    regions   map[string]RegionData
    cities    map[string]map[string]CityData  
    entities  map[string]EntityTypeData
    locations map[string]LocationTypeData
    
    once sync.Once // Ensure single initialization
}
```

## Implementation Strategy

### Phase 1: Core Parsing (Priority 1)

1. **Basic Pattern Matching**
   - Implement character class patterns (a, a1, a2, n, x, x1, x2)
   - Create entity code patterns for Tables B-E
   - Build place, region, and network site patterns
   - Implement non-building and customer location patterns

2. **CLLI Struct and Parsing**
   - Define core CLLI struct with all fields
   - Implement `Parse()` function with strict validation
   - Add `ParseWithOptions()` for flexible parsing
   - Create `MustParse()` for panic-based parsing

3. **Basic Validation**
   - Implement component validation functions
   - Add format validation for each CLLI part
   - Create detailed error types and messages

### Phase 2: Geographic Resolution (Priority 2)

1. **Region Code Mapping**
   - Load region to ISO country/state mappings
   - Implement country code resolution
   - Add state/province code resolution
   - Create country/state name resolution

2. **City Name Resolution**
   - Load place + region to city mappings
   - Implement city name lookup functionality
   - Handle unknown or unmapped locations gracefully

3. **Data Loading**
   - Embed YAML data files in binary
   - Implement lazy loading with sync.Once
   - Add data validation and error handling
   - Cache parsed data structures

### Phase 3: Entity/Location Types (Priority 2)

1. **Entity Type Resolution**
   - Load entity code to description mappings
   - Implement pattern-based entity type lookup
   - Support switching/non-switching categorization
   - Handle entity code variations

2. **Location Type Resolution**
   - Load location code descriptions
   - Implement location type lookup
   - Support both non-building and customer locations

### Phase 4: Advanced Features (Priority 3)

1. **Pattern Recognition**
   - Add `IsEntityCLLI()`, `IsNonBuildingCLLI()`, `IsCustomerCLLI()`
   - Implement type detection without full parsing
   - Optimize pattern matching performance

2. **Enhanced Validation**
   - Add component-specific validation functions
   - Implement cross-component validation rules
   - Support custom validation modes

3. **Performance Optimization**
   - Benchmark parsing operations
   - Optimize regex compilation and caching
   - Minimize memory allocations
   - Add concurrent safety measures

### Phase 5: Testing and Documentation (Ongoing)

1. **Comprehensive Testing**
   - Unit tests for all parsing functions
   - Integration tests with real CLLI data
   - Property-based testing for pattern validation
   - Benchmark tests for performance validation
   - Test coverage > 90%

2. **Documentation**
   - Complete API documentation with examples
   - Performance characteristics documentation
   - Migration guide from Ruby gem
   - Usage patterns and best practices

## Pattern Implementation Details

### Character Classes

Based on the Ruby implementation, define character class constants:

```go
const (
    // Alphabetic characters
    charClassA  = "A-Z"                    // All letters
    charClassA1 = "ACE-HJ-NP-SVXZ"       // Excluding B,D,I,O,T,U,W,Y  
    charClassA2 = "A-FH-Z"               // Excluding G
    
    // Numeric characters
    charClassN = "0-9"                    // All digits
    
    // Combined classes
    charClassX  = charClassA + charClassN  // All alphanumeric
    charClassX1 = charClassA1 + charClassN // A1 + numeric
    charClassX2 = charClassA2 + charClassN // A2 + numeric
)
```

### Entity Code Patterns

Implement patterns for each entity table:

```go
// Table B - Switching Entities
var switchingEntityPatterns = []string{
    `(MG|SG|CG|DS|RL|PS|RP|CM|VS|OS|OL|[0-9]{2})[` + charClassX1 + `]`,
    `[CB0-9][0-9]T`,
    `[0-9]GT`, 
    `Z[A-Z]Z`,
    `RS[0-9]`,
    `X[A-Z]X`,
    `CT[` + charClassX1 + `]`,
}

// Table C - Switchboard and Desk Entities  
var switchboardEntityPattern = `[0-9][CDBINQWMVROLPEUTZ0-9]B`

// Table D - Miscellaneous Switching Entities
var miscEntityPatterns = []string{
    `[0-9][AXCTWDEINPQ]D`,
    `[A-Z0-9][UM]D`,
}

// Table E - Non-Switching Entities
var nonSwitchingPatterns = []string{
    `[FAEKMPSTW][` + charClassX2 + `][` + charClassX1 + `]`,
    `Q[0-9][0-9]`,
}
```

### Complete CLLI Pattern

Build the complete CLLI matching pattern:

```go
func buildCLLIPattern(strict bool) string {
    place := buildPlacePattern(strict)
    region := `[A-Z]{2}`
    networkSite := `([A-Z]{2}|[0-9]{2})`
    entityCode := buildEntityCodePattern(strict)
    nonBuildingLocation := `[A-Z][0-9]{4}`
    customerLocation := `[0-9][A-Z][0-9]{3}`
    
    return fmt.Sprintf(`^%s%s(?:%s(?:%s)?|%s|%s)$`,
        place, region, networkSite, entityCode, 
        nonBuildingLocation, customerLocation)
}
```

## Data Structure Design

### Region Data

```go
type RegionData struct {
    CountryCode string `yaml:"country_code"` // ISO 3166 country
    CountryName string `yaml:"country_name"`
    StateCode   string `yaml:"state_code"`   // ISO 3166 state  
    StateName   string `yaml:"state_name"`
}
```

### City Data

```go
type CityData struct {
    Name      string `yaml:"name"`
    AltNames  []string `yaml:"alt_names,omitempty"`
    TimeZone  string `yaml:"timezone,omitempty"`
    Latitude  float64 `yaml:"latitude,omitempty"`
    Longitude float64 `yaml:"longitude,omitempty"`  
}
```

### Entity Type Data

```go
type EntityTypeData struct {
    Description string   `yaml:"description"`
    Category    string   `yaml:"category"`    // switching, non-switching, etc.
    SubCategory string   `yaml:"subcategory"` // specific equipment type
    Patterns    []string `yaml:"patterns"`    // matching patterns
}
```

## Error Handling Strategy

Define specific error types for different failure modes:

```go
var (
    ErrInvalidCLLI     = errors.New("invalid CLLI format")
    ErrInvalidPlace    = errors.New("invalid place code")
    ErrInvalidRegion   = errors.New("invalid region code") 
    ErrInvalidSite     = errors.New("invalid network site code")
    ErrInvalidEntity   = errors.New("invalid entity code")
    ErrInvalidLocation = errors.New("invalid location code")
)

type ParseError struct {
    Input    string // Original input
    Position int    // Error position (0-based)
    Field    string // Field name that failed
    Err      error  // Underlying error
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("parse error at position %d in field %s: %v", 
        e.Position, e.Field, e.Err)
}
```

## Performance Considerations

### Optimization Strategies

1. **Compiled Patterns**: Pre-compile all regex patterns at package initialization
2. **String Reuse**: Avoid unnecessary string allocations during parsing  
3. **Lazy Loading**: Load data files only when geographic/entity functions are called
4. **Caching**: Cache parsed data structures in memory
5. **Concurrent Safety**: Use sync.Once and read-only data structures

### Memory Management

```go
// Singleton pattern for global data
var (
    globalPatterns *PatternMatcher
    globalData     *DataManager
    initOnce       sync.Once
)

func init() {
    initOnce.Do(func() {
        globalPatterns = compilePatterns()
        globalData = &DataManager{}
    })
}
```

### Benchmarking Targets

- Parse typical CLLI (11 chars): < 1 microsecond
- Geographic resolution: < 10 microseconds  
- Memory per CLLI instance: < 200 bytes
- Concurrent goroutines: No contention

## Testing Strategy

### Unit Test Coverage

1. **Pattern Matching**: Test all entity code patterns from Tables B-E
2. **Component Validation**: Test each CLLI component validation
3. **Geographic Resolution**: Test region/city lookups with known data
4. **Error Handling**: Test all error conditions and edge cases
5. **Performance**: Benchmark parsing and resolution operations

### Test Data

Use the same test cases from the Ruby implementation:

- Valid CLLI examples for each type
- Invalid CLLI examples with expected errors
- Real-world CLLI data from telecommunications datasets
- Edge cases and boundary conditions

### Integration Testing

Test against real CLLI datasets:

- NANPA thousands file (North American numbering plan)
- Telecommunications equipment inventories
- Network topology databases

## Migration from Ruby Implementation

### API Compatibility

Maintain similar method names and behavior where possible:

```go
// Ruby: clli.place
// Go:   clli.Place

// Ruby: clli.city_name  
// Go:   clli.CityName()

// Ruby: clli.state_code
// Go:   clli.StateCode()
```

### Data File Compatibility

Convert YAML data files from Ruby implementation:

- Use same data structure where possible
- Handle differences in YAML parsing
- Validate data integrity during conversion

## Deployment and Distribution

### Package Structure

```text
github.com/dbitech/go-clli/
├── pkg/clli/           # Main package
├── cmd/clli-cli/       # Command-line tool (optional)  
├── examples/           # Usage examples
├── docs/              # Documentation
├── data/              # Source data files
└── scripts/           # Build and test scripts
```

### Release Strategy

1. **Alpha Release**: Core parsing functionality
2. **Beta Release**: Geographic and entity resolution
3. **Stable Release**: Complete feature parity with Ruby gem
4. **Enhanced Release**: Go-specific optimizations and features

This design provides a solid foundation for implementing the CLLI Go package with high performance, comprehensive functionality, and maintainable code structure.
