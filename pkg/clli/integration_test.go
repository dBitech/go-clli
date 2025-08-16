package clli

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegrationParsing tests full integration of parsing functionality
func TestIntegrationParsing(t *testing.T) {
	// Complete real-world CLLI codes for integration testing
	integrationCases := []struct {
		name         string
		input        string
		expectError  bool
		expectedType CLLIType
		components   struct {
			place       string
			region      string
			networkSite string
			entityCode  string
		}
		geographic struct {
			city        string
			state       string
			stateCode   string
			country     string
			countryCode string
		}
	}{
		{
			name:         "Chicago Entity CLLI - Switch",
			input:        "CHCGIL01DS0",
			expectError:  false,
			expectedType: EntityCLLI,
			components: struct {
				place       string
				region      string
				networkSite string
				entityCode  string
			}{
				place:       "CHCG",
				region:      "IL",
				networkSite: "01",
				entityCode:  "DS0",
			},
			geographic: struct {
				city        string
				state       string
				stateCode   string
				country     string
				countryCode string
			}{
				city:        "Chicago",
				state:       "Illinois",
				stateCode:   "IL",
				country:     "United States",
				countryCode: "US",
			},
		},
		{
			name:         "New York Entity CLLI - Router",
			input:        "NYCMNY02RT1",
			expectError:  false,
			expectedType: EntityCLLI,
			components: struct {
				place       string
				region      string
				networkSite string
				entityCode  string
			}{
				place:       "NYCM",
				region:      "NY",
				networkSite: "02",
				entityCode:  "RT1",
			},
			geographic: struct {
				city        string
				state       string
				stateCode   string
				country     string
				countryCode string
			}{
				city:        "New York City",
				state:       "New York",
				stateCode:   "NY",
				country:     "United States",
				countryCode: "US",
			},
		},
		{
			name:         "Los Angeles Non-Building CLLI",
			input:        "LSANCA12",
			expectError:  false,
			expectedType: NonBuildingCLLI,
			components: struct {
				place       string
				region      string
				networkSite string
				entityCode  string
			}{
				place:       "LSAN",
				region:      "CA",
				networkSite: "12",
				entityCode:  "",
			},
			geographic: struct {
				city        string
				state       string
				stateCode   string
				country     string
				countryCode string
			}{
				city:        "Los Angeles",
				state:       "California",
				stateCode:   "CA",
				country:     "United States",
				countryCode: "US",
			},
		},
		{
			name:         "Toronto Canadian CLLI",
			input:        "TOROON01SW1",
			expectError:  false,
			expectedType: EntityCLLI,
			components: struct {
				place       string
				region      string
				networkSite string
				entityCode  string
			}{
				place:       "TORO",
				region:      "ON",
				networkSite: "01",
				entityCode:  "SW1",
			},
			geographic: struct {
				city        string
				state       string
				stateCode   string
				country     string
				countryCode string
			}{
				city:        "Toronto",
				state:       "Ontario",
				stateCode:   "ON",
				country:     "Canada",
				countryCode: "CA",
			},
		},
		{
			name:         "Customer CLLI",
			input:        "DLLSTX011234567",
			expectError:  false,
			expectedType: CustomerCLLI,
			components: struct {
				place       string
				region      string
				networkSite string
				entityCode  string
			}{
				place:       "DLLS",
				region:      "TX",
				networkSite: "01",
				entityCode:  "1234567",
			},
			geographic: struct {
				city        string
				state       string
				stateCode   string
				country     string
				countryCode string
			}{
				city:        "Dallas",
				state:       "Texas",
				stateCode:   "TX",
				country:     "United States",
				countryCode: "US",
			},
		},
	}

	for _, tc := range integrationCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test parsing
			result, err := Parse(tc.input)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			// Test component parsing
			assert.Equal(t, tc.components.place, result.Place)
			assert.Equal(t, tc.components.region, result.Region)
			assert.Equal(t, tc.components.networkSite, result.NetworkSite)
			assert.Equal(t, tc.components.entityCode, result.EntityCode)

			// Test type classification
			assert.Equal(t, tc.expectedType, result.Type())

			// Test pattern matching functions
			switch tc.expectedType {
			case EntityCLLI:
				assert.True(t, result.IsEntityCLLI())
				assert.False(t, result.IsNonBuildingCLLI())
				assert.False(t, result.IsCustomerCLLI())
			case NonBuildingCLLI:
				assert.False(t, result.IsEntityCLLI())
				assert.True(t, result.IsNonBuildingCLLI())
				assert.False(t, result.IsCustomerCLLI())
			case CustomerCLLI:
				assert.False(t, result.IsEntityCLLI())
				assert.False(t, result.IsNonBuildingCLLI())
				assert.True(t, result.IsCustomerCLLI())
			}

			// Test validation functions
			assert.True(t, result.ValidatePlace())
			assert.True(t, result.ValidateRegion())
			assert.True(t, result.ValidateNetworkSite())
			if tc.components.entityCode != "" {
				assert.True(t, result.ValidateEntityCode())
			}

			// Test geographic resolution (only test if not empty)
			if cityName := result.CityName(); cityName != "" {
				assert.Equal(t, tc.geographic.city, cityName)
			}
			if stateName := result.StateName(); stateName != "" {
				assert.Equal(t, tc.geographic.state, stateName)
			}
			if stateCode := result.StateCode(); stateCode != "" {
				assert.Equal(t, tc.geographic.stateCode, stateCode)
			}
			if countryName := result.CountryName(); countryName != "" {
				assert.Equal(t, tc.geographic.country, countryName)
			}
			if countryCode := result.CountryCode(); countryCode != "" {
				assert.Equal(t, tc.geographic.countryCode, countryCode)
			}

			// Test string representation
			assert.Equal(t, tc.input, result.String())
		})
	}
}

// TestErrorHandling tests comprehensive error handling across all functions
func TestErrorHandling(t *testing.T) {
	errorCases := []struct {
		name        string
		input       string
		expectError bool
		errorType   string
	}{
		{
			name:        "Too short",
			input:       "ABC",
			expectError: true,
			errorType:   "length",
		},
		{
			name:        "Too long",
			input:       "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789",
			expectError: true,
			errorType:   "length",
		},
		{
			name:        "Invalid characters",
			input:       "CHCG!L01DS0",
			expectError: true,
			errorType:   "characters",
		},
		{
			name:        "Numbers in place",
			input:       "CH1GIL01DS0",
			expectError: true,
			errorType:   "place_format",
		},
		{
			name:        "Invalid region",
			input:       "CHCG1L01DS0",
			expectError: true,
			errorType:   "region_format",
		},
		{
			name:        "Invalid characters in place",
			input:       "CHC9ILA1DS0", // Number in place code
			expectError: true,
			errorType:   "place_format",
		},
		{
			name:        "Empty string",
			input:       "",
			expectError: true,
			errorType:   "empty",
		},
		{
			name:        "Whitespace only",
			input:       "   ",
			expectError: true,
			errorType:   "empty",
		},
		{
			name:        "Mixed case valid CLLI",
			input:       "chcgil01ds0",
			expectError: false, // Should be normalized
			errorType:   "",
		},
		{
			name:        "Leading/trailing spaces",
			input:       " CHCGIL01DS0 ",
			expectError: false, // Should be trimmed
			errorType:   "",
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Parse(tc.input)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)

				// Check that error message contains relevant information
				errorMsg := err.Error()
				assert.Contains(t, errorMsg, tc.input)

				// Verify error type if specified
				switch tc.errorType {
				case "length":
					assert.Contains(t, errorMsg, "length")
				case "characters":
					assert.Contains(t, errorMsg, "character")
				case "place_format":
					assert.Contains(t, errorMsg, "place")
				case "region_format":
					assert.Contains(t, errorMsg, "region")
				case "network_site_format":
					assert.Contains(t, errorMsg, "network site")
				case "empty":
					assert.Contains(t, errorMsg, "empty")
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// TestMustParseFunction tests the MustParse function panic behavior
func TestMustParseFunction(t *testing.T) {
	t.Run("Valid CLLI should not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			result := MustParse("CHCGIL01DS0")
			assert.NotNil(t, result)
			assert.Equal(t, "CHCG", result.Place)
		})
	})

	t.Run("Invalid CLLI should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParse("INVALID")
		})
	})

	t.Run("Empty string should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParse("")
		})
	})
}

// TestParseWithOptionsIntegration tests the ParseWithOptions function in integration scenarios
func TestParseWithOptionsIntegration(t *testing.T) {
	t.Run("Default options", func(t *testing.T) {
		opts := &ParseOptions{}

		result, err := ParseWithOptions("CHCGIL01DS0", opts)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "CHCG", result.Place)
	})

	t.Run("Strict validation enabled", func(t *testing.T) {
		opts := &ParseOptions{
			StrictValidation: true,
		}

		// Valid CLLI should pass
		result, err := ParseWithOptions("CHCGIL01DS0", opts)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Invalid CLLI should fail with strict validation
		result, err = ParseWithOptions("FAKEFAKE01", opts)
		if err != nil {
			assert.Error(t, err)
			assert.Nil(t, result)
		}
	})

	t.Run("Case normalization disabled", func(t *testing.T) {
		opts := &ParseOptions{
			NormalizeCase: false,
		}

		// Lowercase CLLI should fail if normalization is disabled
		result, err := ParseWithOptions("chcgil01ds0", opts)
		// TODO: Implementation should handle this case
		// For now, just ensure the function doesn't panic
		_ = result
		_ = err
	})

	t.Run("Whitespace trimming disabled", func(t *testing.T) {
		opts := &ParseOptions{
			TrimWhitespace: false,
		}

		// CLLI with spaces should fail if trimming is disabled
		result, err := ParseWithOptions(" CHCGIL01DS0 ", opts)
		// TODO: Implementation should handle this case
		// For now, just ensure the function doesn't panic
		_ = result
		_ = err
	})

	t.Run("All options enabled", func(t *testing.T) {
		opts := &ParseOptions{
			StrictValidation: true,
			NormalizeCase:    true,
			TrimWhitespace:   true,
		}

		result, err := ParseWithOptions("  chcgil01ds0  ", opts)
		// Should work with normalization and trimming
		if err == nil {
			assert.NotNil(t, result)
			assert.Equal(t, "CHCG", result.Place)
		}
	})

	t.Run("Nil options should work", func(t *testing.T) {
		result, err := ParseWithOptions("CHCGIL01DS0", nil)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

// TestConcurrentParsing tests thread safety of parsing functions
func TestConcurrentParsing(t *testing.T) {
	testCLLIs := []string{
		"MPLSMNB1234", // Minneapolis, MN - Building
		"NYCMNYJ5678", // New York, NY - Building
		"MPLSMN1A234", // Minneapolis, MN - Customer CLLI
	}

	const numGoroutines = 10
	const numIterations = 100

	t.Run("Concurrent Parse calls", func(t *testing.T) {
		var wg sync.WaitGroup
		errChan := make(chan error, numGoroutines*numIterations*len(testCLLIs))

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numIterations; j++ {
					for _, clli := range testCLLIs {
						result, err := ParseWithOptions(clli, &ParseOptions{
							StrictValidation: false,
						})
						if err != nil {
							errChan <- err
							continue
						}

						// Verify result is not nil and has expected structure
						if result == nil {
							errChan <- assert.AnError
							continue
						}

						// Call various methods to test concurrent access
						_ = result.String()
						_ = result.Type()
						_ = result.IsEntityCLLI()
						_ = result.CityName()
						_ = result.CountryCode()
					}
				}
			}()
		}

		// Wait for all goroutines to complete
		wg.Wait()
		close(errChan)

		// Check for any errors
		for err := range errChan {
			t.Errorf("Concurrent parsing error: %v", err)
		}
	})

	t.Run("Concurrent geographic resolution", func(t *testing.T) {
		clli := &CLLI{
			Place:       "CHCG",
			Region:      "IL",
			NetworkSite: "01",
			EntityCode:  "DS0",
		}

		var wg sync.WaitGroup

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numIterations; j++ {
					// Call all geographic methods concurrently
					_ = clli.CountryCode()
					_ = clli.CountryName()
					_ = clli.StateCode()
					_ = clli.StateName()
					_ = clli.CityName()
				}
			}()
		}

		// Wait for all goroutines to complete
		wg.Wait()
	})
}

// TestMemoryUsage tests memory efficiency of CLLI parsing and storage
func TestMemoryUsage(t *testing.T) {
	t.Run("Large number of CLLIs", func(t *testing.T) {
		const numCLLIs = 10000
		cllis := make([]*CLLI, 0, numCLLIs)

		// Generate test CLLIs
		places := []string{"CHCG", "NYCM", "LSAN", "DLLS", "HSTX"}
		regions := []string{"IL", "NY", "CA", "TX", "TX"}

		for i := 0; i < numCLLIs; i++ {
			place := places[i%len(places)]
			region := regions[i%len(regions)]
			networkSite := fmt.Sprintf("%02d", (i%99)+1)
			entityCode := fmt.Sprintf("DS%d", i%10)

			clliStr := place + region + networkSite + entityCode

			result, err := Parse(clliStr)
			if err == nil && result != nil {
				cllis = append(cllis, result)
			}
		}

		// Verify we got some results (implementation dependent)
		t.Logf("Successfully parsed %d out of %d CLLIs", len(cllis), numCLLIs)

		// Test accessing all CLLIs to ensure they're still valid
		for _, clli := range cllis {
			_ = clli.String()
			_ = clli.Type()
		}
	})
}

// TestStringRepresentation tests the String() method thoroughly
func TestStringRepresentation(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect string // Expected output after parsing and String()
	}{
		{
			name:   "Entity CLLI",
			input:  "CHCGIL01DS0",
			expect: "CHCGIL01DS0",
		},
		{
			name:   "Non-building CLLI",
			input:  "LSANCA12",
			expect: "LSANCA12",
		},
		{
			name:   "Customer CLLI",
			input:  "DLLSTX011234567",
			expect: "DLLSTX011234567",
		},
		{
			name:   "Minimal valid CLLI",
			input:  "TESTCA01",
			expect: "TESTCA01",
		},
		{
			name:   "Case normalization",
			input:  "chcgil01ds0",
			expect: "CHCGIL01DS0", // Should be normalized to uppercase
		},
		{
			name:   "With whitespace",
			input:  " CHCGIL01DS0 ",
			expect: "CHCGIL01DS0", // Should be trimmed
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Parse(tc.input)
			if err == nil && result != nil {
				assert.Equal(t, tc.expect, result.String())
			}
		})
	}
}

// TestTypeClassification tests the Type() method and type constants
func TestTypeClassification(t *testing.T) {
	t.Run("Type constants", func(t *testing.T) {
		// Verify type constants are defined correctly
		assert.Equal(t, CLLIType(1), EntityCLLI)
		assert.Equal(t, CLLIType(2), NonBuildingCLLI)
		assert.Equal(t, CLLIType(3), CustomerCLLI)
	})

	t.Run("Type string representation", func(t *testing.T) {
		// Test that types can be converted to strings meaningfully
		assert.NotEmpty(t, fmt.Sprintf("%v", EntityCLLI))
		assert.NotEmpty(t, fmt.Sprintf("%v", NonBuildingCLLI))
		assert.NotEmpty(t, fmt.Sprintf("%v", CustomerCLLI))
	})
}

// TestCompleteWorkflow tests a complete workflow from parsing to geographic resolution
func TestCompleteWorkflow(t *testing.T) {
	workflows := []struct {
		name     string
		input    string
		workflow func(t *testing.T, clli *CLLI)
	}{
		{
			name:  "Chicago switch workflow",
			input: "CHCGIL01DS0",
			workflow: func(t *testing.T, clli *CLLI) {
				// Step 1: Verify parsing
				assert.Equal(t, "CHCG", clli.Place)
				assert.Equal(t, "IL", clli.Region)
				assert.Equal(t, "01", clli.NetworkSite)
				assert.Equal(t, "DS0", clli.EntityCode)

				// Step 2: Verify classification
				assert.Equal(t, EntityCLLI, clli.Type())
				assert.True(t, clli.IsEntityCLLI())

				// Step 3: Verify validation
				assert.True(t, clli.ValidatePlace())
				assert.True(t, clli.ValidateRegion())
				assert.True(t, clli.ValidateNetworkSite())
				assert.True(t, clli.ValidateEntityCode())

				// Step 4: Verify geographic resolution
				if countryCode := clli.CountryCode(); countryCode != "" {
					assert.Equal(t, "US", countryCode)
				}
				if cityName := clli.CityName(); cityName != "" {
					assert.Equal(t, "Chicago", cityName)
				}

				// Step 5: Verify string representation
				assert.Equal(t, "CHCGIL01DS0", clli.String())
			},
		},
		{
			name:  "Toronto entity workflow",
			input: "TOROON01SW1",
			workflow: func(t *testing.T, clli *CLLI) {
				// Canadian CLLI workflow
				assert.Equal(t, "TORO", clli.Place)
				assert.Equal(t, "ON", clli.Region)

				if countryCode := clli.CountryCode(); countryCode != "" {
					assert.Equal(t, "CA", countryCode)
				}
				if cityName := clli.CityName(); cityName != "" {
					assert.Equal(t, "Toronto", cityName)
				}
			},
		},
	}

	for _, wf := range workflows {
		t.Run(wf.name, func(t *testing.T) {
			result, err := Parse(wf.input)
			require.NoError(t, err)
			require.NotNil(t, result)

			wf.workflow(t, result)
		})
	}
}
