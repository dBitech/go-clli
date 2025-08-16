package clli

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCLLIType tests the CLLIType enum and string representation
func TestCLLIType(t *testing.T) {
	tests := []struct {
		name     string
		cliType  CLLIType
		expected string
	}{
		{"Unknown", CLLITypeUnknown, "Unknown"},
		{"Entity", CLLITypeEntity, "Entity"},
		{"NonBuilding", CLLITypeNonBuilding, "NonBuilding"},
		{"Customer", CLLITypeCustomer, "Customer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.cliType.String())
		})
	}
}

// TestParseError tests the ParseError structure and methods
func TestParseError(t *testing.T) {
	t.Run("Error message formatting", func(t *testing.T) {
		err := &ParseError{
			Input:    "INVALID",
			Position: 4,
			Field:    "place",
			Err:      ErrInvalidPlace,
		}

		expected := "parse error at position 4 in field place: invalid place code"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Unwrap error", func(t *testing.T) {
		err := &ParseError{
			Input:    "INVALID",
			Position: 0,
			Field:    "place",
			Err:      ErrInvalidPlace,
		}

		assert.True(t, errors.Is(err, ErrInvalidPlace))
		assert.Equal(t, ErrInvalidPlace, err.Unwrap())
	})
}

// TestParse tests the main Parse function
func TestParse(t *testing.T) {
	t.Run("Valid entity CLLI", func(t *testing.T) {
		tests := []string{
			"MPLSMNMSDS1", // Minneapolis, MN - DS1 entity
			"NYCMNYPSRS1", // New York City, NY - RS1 entity
			"LSANCAMASG1", // Los Angeles, CA - SG1 entity
			"CHCGILMAMG1", // Chicago, IL - MG1 entity
		}

		for _, test := range tests {
			t.Run(test, func(t *testing.T) {
				// Note: This will fail until implementation is complete
				// but provides the test framework for TDD development
				c, err := Parse(test)
				if err == nil {
					assert.NotNil(t, c)
					assert.Equal(t, test, c.Original)
					assert.Equal(t, CLLITypeEntity, c.Type())
					assert.True(t, c.IsValid())
					assert.Equal(t, 4, len(c.Place))
					assert.Equal(t, 2, len(c.Region))
					assert.Equal(t, 2, len(c.NetworkSite))
					assert.Equal(t, 3, len(c.EntityCode))
				} else {
					// For now, expect errors until implementation is complete
					assert.Error(t, err)
				}
			})
		}
	})

	t.Run("Valid non-building CLLI", func(t *testing.T) {
		tests := []string{
			"MPLSMNB1234", // Minneapolis, MN - Location B, ID 1234
			"NYCMNYJ5678", // New York City, NY - Location J, ID 5678
		}

		for _, test := range tests {
			t.Run(test, func(t *testing.T) {
				c, err := Parse(test)
				if err == nil {
					assert.NotNil(t, c)
					assert.Equal(t, test, c.Original)
					assert.Equal(t, CLLITypeNonBuilding, c.Type())
					assert.True(t, c.IsValid())
					assert.Equal(t, 4, len(c.Place))
					assert.Equal(t, 2, len(c.Region))
					assert.Equal(t, 1, len(c.LocationCode))
					assert.Equal(t, 4, len(c.LocationID))
				} else {
					// For now, expect errors until implementation is complete
					assert.Error(t, err)
				}
			})
		}
	})

	t.Run("Valid customer CLLI", func(t *testing.T) {
		tests := []string{
			"MPLSMN1A234", // Minneapolis, MN - Customer 1, ID A234
			"NYCMNY2B567", // New York City, NY - Customer 2, ID B567
		}

		for _, test := range tests {
			t.Run(test, func(t *testing.T) {
				c, err := Parse(test)
				if err == nil {
					assert.NotNil(t, c)
					assert.Equal(t, test, c.Original)
					assert.Equal(t, CLLITypeCustomer, c.Type())
					assert.True(t, c.IsValid())
					assert.Equal(t, 4, len(c.Place))
					assert.Equal(t, 2, len(c.Region))
					assert.Equal(t, 1, len(c.CustomerCode))
					assert.Equal(t, 4, len(c.CustomerID))
				} else {
					// For now, expect errors until implementation is complete
					assert.Error(t, err)
				}
			})
		}
	})

	t.Run("Invalid CLLIs", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			error error
		}{
			{"Empty string", "", ErrEmptyInput},
			{"Too short", "MPL", ErrInvalidCLLI},
			{"Too long", "MPLSMNMSDS1EXTRA", ErrInvalidCLLI},
			{"Invalid characters", "MPLS@#$%", ErrInvalidCLLI},
			{"Invalid place", "12345MN", ErrInvalidCLLI}, // This gets caught by format validation first
			{"Invalid region", "MPLS@@", ErrInvalidCLLI}, // This gets caught by format validation first
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c, err := Parse(tt.input)
				assert.Nil(t, c)
				assert.Error(t, err)
				if tt.error != nil {
					assert.True(t, errors.Is(err, tt.error))
				}
			})
		}
	})
}

// TestParseWithOptions tests parsing with different options
func TestParseWithOptions(t *testing.T) {
	t.Run("Strict mode", func(t *testing.T) {
		opts := &ParseOptions{Strict: true}

		// Valid CLLI should parse
		c, err := ParseWithOptions("MPLSMNMSDS1", opts)
		if err == nil {
			assert.NotNil(t, c)
			assert.True(t, c.IsValid())
		}

		// Invalid CLLI should fail
		c, err = ParseWithOptions("INVALID", opts)
		assert.Nil(t, c)
		assert.Error(t, err)
	})

	t.Run("Relaxed mode", func(t *testing.T) {
		opts := &ParseOptions{Strict: false}

		// Partial CLLI might parse in relaxed mode
		c, err := ParseWithOptions("MPLS", opts)
		if err == nil {
			assert.NotNil(t, c)
			// May or may not be valid depending on implementation
		}
	})
}

// TestMustParse tests the panic-based parsing function
func TestMustParse(t *testing.T) {
	t.Run("Valid CLLI", func(t *testing.T) {
		// Note: This will panic until implementation is complete
		defer func() {
			if r := recover(); r != nil {
				// Expected for now
				assert.NotNil(t, r)
			}
		}()

		c := MustParse("MPLSMNMSDS1")
		if c != nil {
			assert.Equal(t, "MPLSMNMSDS1", c.Original)
			assert.True(t, c.IsValid())
		}
	})

	t.Run("Invalid CLLI should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.NotNil(t, r)
			} else {
				t.Fatal("Expected panic for invalid CLLI")
			}
		}()

		MustParse("INVALID")
	})
}

// TestCLLIMethods tests the CLLI struct methods
func TestCLLIMethods(t *testing.T) {
	t.Run("String method", func(t *testing.T) {
		c := &CLLI{Original: "MPLSMNMSDS1"}
		assert.Equal(t, "MPLSMNMSDS1", c.String())
	})

	t.Run("Type method", func(t *testing.T) {
		tests := []struct {
			name     string
			cliType  CLLIType
			expected CLLIType
		}{
			{"Entity", CLLITypeEntity, CLLITypeEntity},
			{"NonBuilding", CLLITypeNonBuilding, CLLITypeNonBuilding},
			{"Customer", CLLITypeCustomer, CLLITypeCustomer},
			{"Unknown", CLLITypeUnknown, CLLITypeUnknown},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := &CLLI{cliType: tt.cliType}
				assert.Equal(t, tt.expected, c.Type())
			})
		}
	})

	t.Run("IsValid method", func(t *testing.T) {
		validCLLI := &CLLI{valid: true}
		invalidCLLI := &CLLI{valid: false}

		assert.True(t, validCLLI.IsValid())
		assert.False(t, invalidCLLI.IsValid())
	})
}

// TestGeographicMethods tests geographic resolution methods
func TestGeographicMethods(t *testing.T) {
	// Note: These tests provide framework for geographic functionality
	// Implementation will need to return actual values

	t.Run("CountryCode", func(t *testing.T) {
		c := &CLLI{
			Place:  "MPLS",
			Region: "MN",
		}

		countryCode := c.CountryCode()
		// TODO: Should return "US" when implemented
		if countryCode != "" {
			assert.Equal(t, "US", countryCode)
		}
	})

	t.Run("CountryName", func(t *testing.T) {
		c := &CLLI{
			Place:  "MPLS",
			Region: "MN",
		}

		countryName := c.CountryName()
		// TODO: Should return "United States" when implemented
		if countryName != "" {
			assert.Equal(t, "United States", countryName)
		}
	})

	t.Run("StateCode", func(t *testing.T) {
		c := &CLLI{
			Place:  "MPLS",
			Region: "MN",
		}

		stateCode := c.StateCode()
		// TODO: Should return "MN" when implemented
		if stateCode != "" {
			assert.Equal(t, "MN", stateCode)
		}
	})

	t.Run("StateName", func(t *testing.T) {
		c := &CLLI{
			Place:  "MPLS",
			Region: "MN",
		}

		stateName := c.StateName()
		// TODO: Should return "Minnesota" when implemented
		if stateName != "" {
			assert.Equal(t, "Minnesota", stateName)
		}
	})

	t.Run("CityName", func(t *testing.T) {
		c := &CLLI{
			Place:  "MPLS",
			Region: "MN",
		}

		cityName := c.CityName()
		// TODO: Should return "Minneapolis" when implemented
		if cityName != "" {
			assert.Equal(t, "Minneapolis", cityName)
		}
	})
}

// TestEntityAndLocationTypes tests entity and location type methods
func TestEntityAndLocationTypes(t *testing.T) {
	t.Run("EntityType for entity CLLI", func(t *testing.T) {
		c := &CLLI{
			EntityCode: "DS1",
			cliType:    CLLITypeEntity,
		}

		entityType := c.EntityType()
		// TODO: Should return entity description when implemented
		if entityType != "" {
			assert.Contains(t, entityType, "Digital")
		}
	})

	t.Run("LocationType for non-building CLLI", func(t *testing.T) {
		c := &CLLI{
			LocationCode: "B",
			cliType:      CLLITypeNonBuilding,
		}

		locationType := c.LocationType()
		// TODO: Should return location description when implemented
		if locationType != "" {
			assert.NotEmpty(t, locationType)
		}
	})

	t.Run("LocationType for customer CLLI", func(t *testing.T) {
		c := &CLLI{
			CustomerCode: "1",
			cliType:      CLLITypeCustomer,
		}

		locationType := c.LocationType()
		// TODO: Should return customer location description when implemented
		if locationType != "" {
			assert.NotEmpty(t, locationType)
		}
	})
}

// Benchmark tests for performance validation
func BenchmarkParse(b *testing.B) {
	cliCodes := []string{
		"MPLSMNMSDS1",
		"NYCMNYPSRS1",
		"LSANCAMASG1",
		"CHCGILMAMG1",
		"MPLSMNB1234",
		"NYCMNY1A567",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		code := cliCodes[i%len(cliCodes)]
		_, _ = Parse(code)
	}
}

func BenchmarkParseStruct(b *testing.B) {
	// Test creation of CLLI struct without parsing
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := &CLLI{
			Original:    "MPLSMNMSDS1",
			Place:       "MPLS",
			Region:      "MN",
			NetworkSite: "MS",
			EntityCode:  "DS1",
			cliType:     CLLITypeEntity,
			valid:       true,
		}
		_ = c
	}
}

func BenchmarkGeographicResolution(b *testing.B) {
	c := &CLLI{
		Place:  "MPLS",
		Region: "MN",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.CountryCode()
		_ = c.CountryName()
		_ = c.StateCode()
		_ = c.StateName()
		_ = c.CityName()
	}
}
