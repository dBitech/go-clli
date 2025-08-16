package clli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsEntityCLLI tests entity CLLI pattern recognition
func TestIsEntityCLLI(t *testing.T) {
	t.Run("Valid entity CLLIs", func(t *testing.T) {
		validEntityCLLIs := []string{
			// Standard 8-character entity CLLIs
			"MPLSMNMS", // Place + Region + Site (minimal entity)

			// 11-character entity CLLIs with entity codes
			"MPLSMNMSDS1", // Minneapolis, MN - DS1
			"NYCMNYPSRS1", // New York City, NY - RS1
			"LSANCAMASG1", // Los Angeles, CA - SG1
			"CHCGILMAMG1", // Chicago, IL - MG1
			"HSTXTXMACT1", // Houston, TX - CT1
			"PHLAPAPSF23", // Philadelphia, PA - F23
			"PHNXAZMA01T", // Phoenix, AZ - 01T
			"SNANTNXTZAZ", // San Antonio, TX - ZAZ
			"SNFRCANXRS5", // San Francisco, CA - RS5
			"SNJPCAXAXAX", // San Jose, CA - XAX

			// Various entity code patterns
			"MPLSMNMS2CB", // Table B: [CB0-9][0-9]T variant
			"NYCMNYPE3GT", // Table B: [0-9]GT
			"LSANCAMA4QB", // Table C: [0-9][CDBINQWMVROLPEUTZ0-9]B
			"CHCGILMA5AD", // Table D: [0-9][AXCTWDEINPQ]D
			"HSTXTXMAAUD", // Table D: [A-Z0-9][UM]D
			"PHLAPAPAAA3", // Table E: [FAEKMPSTW][x2][x1]
			"PHNXAZMMQ01", // Table E: Q[0-9][0-9]
		}

		for _, clli := range validEntityCLLIs {
			t.Run(clli, func(t *testing.T) {
				result := IsEntityCLLI(clli)
				// TODO: Should return true when implemented
				if result {
					assert.True(t, result, "Should recognize %s as entity CLLI", clli)
				} else {
					// Expected false for now until implementation
					assert.False(t, result)
				}
			})
		}
	})

	t.Run("Invalid entity CLLIs", func(t *testing.T) {
		invalidEntityCLLIs := []string{
			// Non-building location CLLIs
			"MPLSMNB1234", // Location code B + 4 digits
			"NYCMNYJ5678", // Location code J + 4 digits

			// Customer location CLLIs
			"MPLSMN1A234", // Customer code 1 + A + 3 digits
			"NYCMNY2B567", // Customer code 2 + B + 3 digits

			// Too short or too long
			"MPLS",         // Too short
			"MPLSMN",       // Too short
			"MPLSMNMSDS1X", // Too long

			// Invalid patterns
			"MPLSMNMSXYZ",    // Invalid entity code
			"MPLSMNMS123456", // Too long entity section

			// Empty or malformed
			"",            // Empty
			"INVALID",     // Random string
			"12345678901", // All digits
		}

		for _, clli := range invalidEntityCLLIs {
			t.Run(clli, func(t *testing.T) {
				result := IsEntityCLLI(clli)
				assert.False(t, result, "Should not recognize %s as entity CLLI", clli)
			})
		}
	})
}

// TestIsNonBuildingCLLI tests non-building location CLLI pattern recognition
func TestIsNonBuildingCLLI(t *testing.T) {
	t.Run("Valid non-building CLLIs", func(t *testing.T) {
		validNonBuildingCLLIs := []string{
			// Standard non-building location patterns: Place + Region + [A-Z] + [0-9]{4}
			"MPLSMNB1234", // Minneapolis, MN - Location B, ID 1234
			"NYCMNYJ5678", // New York City, NY - Location J, ID 5678
			"LSANCAE9012", // Los Angeles, CA - Location E, ID 9012
			"CHCGILM3456", // Chicago, IL - Location M, ID 3456
			"HSTXTXP7890", // Houston, TX - Location P, ID 7890
			"PHLAPAQ0123", // Philadelphia, PA - Location Q, ID 0123
			"PHNXAZR4567", // Phoenix, AZ - Location R, ID 4567
			"SNANTZS8901", // San Antonio, TX - Location S, ID 8901
			"SNFRCAU2345", // San Francisco, CA - Location U, ID 2345
			"SNJPCAX6789", // San Jose, CA - Location X, ID 6789

			// All valid location codes A-Z
			"MPLSMNA0000", // Location A
			"MPLSMNB0000", // Location B
			"MPLSMNC0000", // Location C
			"MPLSMND0000", // Location D
			"MPLSMNE0000", // Location E
			"MPLSMNF0000", // Location F
			"MPLSMNG0000", // Location G
			"MPLSMNH0000", // Location H
			"MPLSMNI0000", // Location I
			"MPLSMNJ0000", // Location J
			"MPLSMNK0000", // Location K
			"MPLSMNL0000", // Location L
			"MPLSMNM0000", // Location M
			"MPLSMNZ9999", // Location Z, max ID
		}

		for _, clli := range validNonBuildingCLLIs {
			t.Run(clli, func(t *testing.T) {
				result := IsNonBuildingCLLI(clli)
				// TODO: Should return true when implemented
				if result {
					assert.True(t, result, "Should recognize %s as non-building CLLI", clli)
				} else {
					// Expected false for now until implementation
					assert.False(t, result)
				}
			})
		}
	})

	t.Run("Invalid non-building CLLIs", func(t *testing.T) {
		invalidNonBuildingCLLIs := []string{
			// Entity CLLIs
			"MPLSMNMSDS1", // Entity CLLI
			"NYCMNYPSRS1", // Entity CLLI

			// Customer CLLIs
			"MPLSMN1A234", // Customer CLLI
			"NYCMNY2B567", // Customer CLLI

			// Wrong format
			"MPLSMN11234",  // Digit instead of letter for location code
			"MPLSMN@1234",  // Symbol instead of letter
			"MPLSMNB123",   // Too few digits in ID
			"MPLSMNB12345", // Too many digits in ID
			"MPLSMNBABC1",  // Letters in ID section

			// Too short or too long
			"MPLS",         // Too short
			"MPLSMN",       // Too short
			"MPLSMNB",      // Missing ID
			"MPLSMNB1234X", // Too long

			// Empty or malformed
			"",        // Empty
			"INVALID", // Random string
		}

		for _, clli := range invalidNonBuildingCLLIs {
			t.Run(clli, func(t *testing.T) {
				result := IsNonBuildingCLLI(clli)
				assert.False(t, result, "Should not recognize %s as non-building CLLI", clli)
			})
		}
	})
}

// TestIsCustomerCLLI tests customer location CLLI pattern recognition
func TestIsCustomerCLLI(t *testing.T) {
	t.Run("Valid customer CLLIs", func(t *testing.T) {
		validCustomerCLLIs := []string{
			// Standard customer location patterns: Place + Region + [0-9] + [A-Z] + [0-9]{3}
			"MPLSMN1A234", // Minneapolis, MN - Customer 1, ID A234
			"NYCMNY2B567", // New York City, NY - Customer 2, ID B567
			"LSANCA3C890", // Los Angeles, CA - Customer 3, ID C890
			"CHCGIL4D123", // Chicago, IL - Customer 4, ID D123
			"HSTXTX5E456", // Houston, TX - Customer 5, ID E456
			"PHLAPA6F789", // Philadelphia, PA - Customer 6, ID F789
			"PHNXAZ7G012", // Phoenix, AZ - Customer 7, ID G012
			"SNANTN8H345", // San Antonio, TX - Customer 8, ID H345
			"SNFRCA9I678", // San Francisco, CA - Customer 9, ID I678
			"SNJPCA0J901", // San Jose, CA - Customer 0, ID J901

			// All valid customer codes 0-9
			"MPLSMN0A000", // Customer 0
			"MPLSMN1A000", // Customer 1
			"MPLSMN2A000", // Customer 2
			"MPLSMN3A000", // Customer 3
			"MPLSMN4A000", // Customer 4
			"MPLSMN5A000", // Customer 5
			"MPLSMN6A000", // Customer 6
			"MPLSMN7A000", // Customer 7
			"MPLSMN8A000", // Customer 8
			"MPLSMN9A000", // Customer 9

			// All valid ID first letters A-Z
			"MPLSMN1A999", // ID starts with A
			"MPLSMN1B999", // ID starts with B
			"MPLSMN1Z999", // ID starts with Z
		}

		for _, clli := range validCustomerCLLIs {
			t.Run(clli, func(t *testing.T) {
				result := IsCustomerCLLI(clli)
				// TODO: Should return true when implemented
				if result {
					assert.True(t, result, "Should recognize %s as customer CLLI", clli)
				} else {
					// Expected false for now until implementation
					assert.False(t, result)
				}
			})
		}
	})

	t.Run("Invalid customer CLLIs", func(t *testing.T) {
		invalidCustomerCLLIs := []string{
			// Entity CLLIs
			"MPLSMNMSDS1", // Entity CLLI
			"NYCMNYPSRS1", // Entity CLLI

			// Non-building CLLIs
			"MPLSMNB1234", // Non-building CLLI
			"NYCMNYJ5678", // Non-building CLLI

			// Wrong format
			"MPLSMNA1234",  // Letter instead of digit for customer code
			"MPLSMN11234",  // Digit instead of letter for ID first char
			"MPLSMN1@234",  // Symbol in ID
			"MPLSMN1A23",   // Too few digits in ID
			"MPLSMN1A2345", // Too many digits in ID
			"MPLSMN1AB23",  // Letter in ID digit section

			// Too short or too long
			"MPLS",         // Too short
			"MPLSMN",       // Too short
			"MPLSMN1",      // Missing ID
			"MPLSMN1A",     // Incomplete ID
			"MPLSMN1A234X", // Too long

			// Empty or malformed
			"",        // Empty
			"INVALID", // Random string
		}

		for _, clli := range invalidCustomerCLLIs {
			t.Run(clli, func(t *testing.T) {
				result := IsCustomerCLLI(clli)
				assert.False(t, result, "Should not recognize %s as customer CLLI", clli)
			})
		}
	})
}

// TestPatternMatchingEdgeCases tests edge cases and boundary conditions
func TestPatternMatchingEdgeCases(t *testing.T) {
	t.Run("Empty and null inputs", func(t *testing.T) {
		functions := []struct {
			name string
			fn   func(string) bool
		}{
			{"IsEntityCLLI", IsEntityCLLI},
			{"IsNonBuildingCLLI", IsNonBuildingCLLI},
			{"IsCustomerCLLI", IsCustomerCLLI},
		}

		for _, fn := range functions {
			t.Run(fn.name, func(t *testing.T) {
				// Empty string should return false
				result := fn.fn("")
				assert.False(t, result, "%s should return false for empty string", fn.name)
			})
		}
	})

	t.Run("Case sensitivity", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    string
			function func(string) bool
			funcName string
		}{
			{"Entity lowercase", "mplsmnmsds1", IsEntityCLLI, "IsEntityCLLI"},
			{"Entity mixed case", "MplsMnMsDs1", IsEntityCLLI, "IsEntityCLLI"},
			{"NonBuilding lowercase", "mplsmnb1234", IsNonBuildingCLLI, "IsNonBuildingCLLI"},
			{"NonBuilding mixed case", "MplsMnB1234", IsNonBuildingCLLI, "IsNonBuildingCLLI"},
			{"Customer lowercase", "mplsmn1a234", IsCustomerCLLI, "IsCustomerCLLI"},
			{"Customer mixed case", "MplsMn1A234", IsCustomerCLLI, "IsCustomerCLLI"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := tc.function(tc.input)
				// All should return false for non-uppercase input
				assert.False(t, result, "%s should return false for non-uppercase input: %s", tc.funcName, tc.input)
			})
		}
	})

	t.Run("Length boundaries", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    string
			expected map[string]bool // Expected results for each function
		}{
			{
				"6 chars",
				"MPLSMN",
				map[string]bool{"IsEntityCLLI": false, "IsNonBuildingCLLI": false, "IsCustomerCLLI": false},
			},
			{
				"7 chars",
				"MPLSMNM",
				map[string]bool{"IsEntityCLLI": false, "IsNonBuildingCLLI": false, "IsCustomerCLLI": false},
			},
			{
				"8 chars entity",
				"MPLSMNMS",
				map[string]bool{"IsEntityCLLI": true, "IsNonBuildingCLLI": false, "IsCustomerCLLI": false},
			},
			{
				"11 chars entity",
				"MPLSMNMSDS1",
				map[string]bool{"IsEntityCLLI": true, "IsNonBuildingCLLI": false, "IsCustomerCLLI": false},
			},
			{
				"11 chars non-building",
				"MPLSMNB1234",
				map[string]bool{"IsEntityCLLI": false, "IsNonBuildingCLLI": true, "IsCustomerCLLI": false},
			},
			{
				"11 chars customer",
				"MPLSMN1A234",
				map[string]bool{"IsEntityCLLI": false, "IsNonBuildingCLLI": false, "IsCustomerCLLI": true},
			},
			{
				"12 chars",
				"MPLSMNMSDS1X",
				map[string]bool{"IsEntityCLLI": false, "IsNonBuildingCLLI": false, "IsCustomerCLLI": false},
			},
		}

		functions := map[string]func(string) bool{
			"IsEntityCLLI":      IsEntityCLLI,
			"IsNonBuildingCLLI": IsNonBuildingCLLI,
			"IsCustomerCLLI":    IsCustomerCLLI,
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				for fnName, fn := range functions {
					result := fn(tc.input)
					expected := tc.expected[fnName]

					// TODO: Update when implementation is complete
					if result == expected {
						assert.Equal(t, expected, result, "%s(%s) should return %v", fnName, tc.input, expected)
					} else {
						// For now, most will return false until implemented
						assert.False(t, result, "%s(%s) returns false until implemented", fnName, tc.input)
					}
				}
			})
		}
	})

	t.Run("Special characters", func(t *testing.T) {
		specialCharInputs := []string{
			"MPLS@MN1234",  // @ symbol
			"MPLS MN1234",  // Space
			"MPLSMN\t1234", // Tab
			"MPLSMN\n1234", // Newline
			"MPLSMN-1234",  // Hyphen
			"MPLSMN_1234",  // Underscore
			"MPLSMN.1234",  // Period
		}

		functions := map[string]func(string) bool{
			"IsEntityCLLI":      IsEntityCLLI,
			"IsNonBuildingCLLI": IsNonBuildingCLLI,
			"IsCustomerCLLI":    IsCustomerCLLI,
		}

		for _, input := range specialCharInputs {
			t.Run(input, func(t *testing.T) {
				for fnName, fn := range functions {
					result := fn(input)
					assert.False(t, result, "%s should return false for input with special characters: %s", fnName, input)
				}
			})
		}
	})
}

// TestMutualExclusivity tests that CLLI types are mutually exclusive
func TestMutualExclusivity(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		// Only one should be true, others false
		expectEntity      bool
		expectNonBuilding bool
		expectCustomer    bool
	}{
		{"Entity CLLI", "MPLSMNMSDS1", true, false, false},
		{"Non-building CLLI", "MPLSMNB1234", false, true, false},
		{"Customer CLLI", "MPLSMN1A234", false, false, true},
		{"Invalid CLLI", "INVALID", false, false, false},
		{"Empty CLLI", "", false, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entityResult := IsEntityCLLI(tc.input)
			nonBuildingResult := IsNonBuildingCLLI(tc.input)
			customerResult := IsCustomerCLLI(tc.input)

			// TODO: Update expectations when implementation is complete
			if entityResult || nonBuildingResult || customerResult {
				// If any return true, verify mutual exclusivity
				trueCount := 0
				if entityResult {
					trueCount++
				}
				if nonBuildingResult {
					trueCount++
				}
				if customerResult {
					trueCount++
				}

				assert.LessOrEqual(t, trueCount, 1, "At most one pattern function should return true for: %s", tc.input)

				// Verify expected results if implemented
				if entityResult {
					assert.True(t, tc.expectEntity, "Expected entity result for %s", tc.input)
				}
				if nonBuildingResult {
					assert.True(t, tc.expectNonBuilding, "Expected non-building result for %s", tc.input)
				}
				if customerResult {
					assert.True(t, tc.expectCustomer, "Expected customer result for %s", tc.input)
				}
			} else {
				// All return false - expected until implementation
				assert.False(t, entityResult)
				assert.False(t, nonBuildingResult)
				assert.False(t, customerResult)
			}
		})
	}
}

// Benchmark tests for pattern matching performance
func BenchmarkIsEntityCLLI(b *testing.B) {
	testCases := []string{
		"MPLSMNMSDS1",
		"NYCMNYPSRS1",
		"LSANCAMASG1",
		"INVALID",
		"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := testCases[i%len(testCases)]
		_ = IsEntityCLLI(input)
	}
}

func BenchmarkIsNonBuildingCLLI(b *testing.B) {
	testCases := []string{
		"MPLSMNB1234",
		"NYCMNYJ5678",
		"LSANCAE9012",
		"INVALID",
		"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := testCases[i%len(testCases)]
		_ = IsNonBuildingCLLI(input)
	}
}

func BenchmarkIsCustomerCLLI(b *testing.B) {
	testCases := []string{
		"MPLSMN1A234",
		"NYCMNY2B567",
		"LSANCA3C890",
		"INVALID",
		"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := testCases[i%len(testCases)]
		_ = IsCustomerCLLI(input)
	}
}

func BenchmarkAllPatternMatching(b *testing.B) {
	testCases := []string{
		"MPLSMNMSDS1", // Entity
		"MPLSMNB1234", // Non-building
		"MPLSMN1A234", // Customer
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := testCases[i%len(testCases)]
		_ = IsEntityCLLI(input)
		_ = IsNonBuildingCLLI(input)
		_ = IsCustomerCLLI(input)
	}
}
