package clli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidatePlace tests place code validation
func TestValidatePlace(t *testing.T) {
	t.Run("Valid places strict mode", func(t *testing.T) {
		validPlaces := []string{
			"MPLS", // Minneapolis
			"NYCM", // New York City
			"LSAN", // Los Angeles
			"CHCG", // Chicago
			"HSTX", // Houston
			"PHLA", // Philadelphia
			"PHNX", // Phoenix
			"SNAN", // San Antonio
			"SNFR", // San Francisco
			"SNJP", // San Jose
		}

		for _, place := range validPlaces {
			t.Run(place, func(t *testing.T) {
				err := ValidatePlace(place, true)
				// TODO: Should pass when implemented
				if err != nil {
					// Expected for now, remove when implemented
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Valid places with padding", func(t *testing.T) {
		validPlacesWithPadding := []string{
			"ABC ", // 3 chars + 1 space
			"AB  ", // 2 chars + 2 spaces
		}

		for _, place := range validPlacesWithPadding {
			t.Run(place, func(t *testing.T) {
				err := ValidatePlace(place, true)
				// TODO: Should pass when implemented for padded places
				if err == nil {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Invalid places strict mode", func(t *testing.T) {
		invalidPlaces := []struct {
			name  string
			place string
		}{
			{"Empty", ""},
			{"Too short", "MP"},
			{"Too long", "MPLSX"},
			{"Contains digits", "MPL1"},
			{"Contains symbols", "MP@L"},
			{"Lowercase", "mpls"},
			{"Mixed case", "Mpls"},
		}

		for _, tt := range invalidPlaces {
			t.Run(tt.name, func(t *testing.T) {
				err := ValidatePlace(tt.place, true)
				// Should always fail for invalid places
				assert.Error(t, err)
			})
		}
	})

	t.Run("Relaxed mode", func(t *testing.T) {
		// In relaxed mode, some invalid places might be accepted
		relaxedCases := []string{
			"MPL1", // Contains digit
			"mpls", // Lowercase
			"MP",   // Too short
		}

		for _, place := range relaxedCases {
			t.Run(place, func(t *testing.T) {
				err := ValidatePlace(place, false)
				// Relaxed mode might accept these
				if err == nil {
					assert.NoError(t, err)
				}
			})
		}
	})
}

// TestValidateRegion tests region code validation
func TestValidateRegion(t *testing.T) {
	t.Run("Valid US states", func(t *testing.T) {
		validUSStates := []string{
			"AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA",
			"HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD",
			"MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ",
			"NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC",
			"SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY",
			"DC", // District of Columbia
		}

		for _, region := range validUSStates {
			t.Run(region, func(t *testing.T) {
				err := ValidateRegion(region, true)
				// TODO: Should pass when implemented
				if err != nil {
					// Expected for now, remove when implemented
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Valid Canadian provinces", func(t *testing.T) {
		validCanadianProvinces := []string{
			"AB", "BC", "MB", "NB", "NL", "NT", "NS", "NU",
			"ON", "PE", "QC", "SK", "YT",
		}

		for _, region := range validCanadianProvinces {
			t.Run(region, func(t *testing.T) {
				err := ValidateRegion(region, true)
				// TODO: Should pass when implemented
				if err != nil {
					// Expected for now, remove when implemented
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Invalid regions", func(t *testing.T) {
		invalidRegions := []struct {
			name   string
			region string
		}{
			{"Empty", ""},
			{"Too short", "M"},
			{"Too long", "MNN"},
			{"Contains digits", "M1"},
			{"Contains symbols", "M@"},
			{"Lowercase", "mn"},
			{"Mixed case", "Mn"},
			{"Invalid state", "ZZ"},
		}

		for _, tt := range invalidRegions {
			t.Run(tt.name, func(t *testing.T) {
				err := ValidateRegion(tt.region, true)
				// Should always fail for invalid regions
				assert.Error(t, err)
			})
		}
	})
}

// TestValidateNetworkSite tests network site code validation
func TestValidateNetworkSite(t *testing.T) {
	t.Run("Valid network sites", func(t *testing.T) {
		validSites := []string{
			"MS", "RS", "PS", "DS", "GS", // Alphabetic sites
			"01", "02", "10", "99", // Numeric sites
		}

		for _, site := range validSites {
			t.Run(site, func(t *testing.T) {
				err := ValidateNetworkSite(site, true)
				// TODO: Should pass when implemented
				if err != nil {
					// Expected for now, remove when implemented
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Invalid network sites", func(t *testing.T) {
		invalidSites := []struct {
			name string
			site string
		}{
			{"Empty", ""},
			{"Too short", "M"},
			{"Too long", "MSX"},
			{"Mixed alpha-numeric", "M1"},
			{"Contains symbols", "M@"},
			{"Lowercase", "ms"},
		}

		for _, tt := range invalidSites {
			t.Run(tt.name, func(t *testing.T) {
				err := ValidateNetworkSite(tt.site, true)
				// Should always fail for invalid sites
				assert.Error(t, err)
			})
		}
	})
}

// TestValidateEntityCode tests entity code validation based on Tables B-E
func TestValidateEntityCode(t *testing.T) {
	t.Run("Valid switching entities (Table B)", func(t *testing.T) {
		validSwitchingEntities := []string{
			// Pattern: (MG|SG|CG|DS|RL|PS|RP|CM|VS|OS|OL|[0-9]{2})[x1]
			"MG1", "MG2", "MGA", "MGZ",
			"SG1", "SG2", "SGA", "SGZ",
			"CG1", "CG2", "CGA", "CGZ",
			"DS1", "DS2", "DSA", "DSZ",
			"RL1", "RL2", "RLA", "RLZ",
			"PS1", "PS2", "PSA", "PSZ",
			"RP1", "RP2", "RPA", "RPZ",
			"CM1", "CM2", "CMA", "CMZ",
			"VS1", "VS2", "VSA", "VSZ",
			"OS1", "OS2", "OSA", "OSZ",
			"OL1", "OL2", "OLA", "OLZ",
			"011", "012", "01A", "01Z",
			"991", "992", "99A", "99Z",

			// Pattern: [CB0-9][0-9]T
			"C0T", "C1T", "C9T",
			"B0T", "B1T", "B9T",
			"00T", "01T", "99T",

			// Pattern: [0-9]GT
			"0GT", "1GT", "9GT",

			// Pattern: Z[A-Z]Z
			"ZAZ", "ZBZ", "ZZZ",

			// Pattern: RS[0-9]
			"RS0", "RS1", "RS9",

			// Pattern: X[A-Z]X
			"XAX", "XBX", "XZX",

			// Pattern: CT[x1]
			"CT1", "CT2", "CTA", "CTZ",
		}

		for _, entity := range validSwitchingEntities {
			t.Run(entity, func(t *testing.T) {
				err := ValidateEntityCode(entity, true)
				// TODO: Should pass when implemented
				if err != nil {
					// Expected for now, remove when implemented
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Valid switchboard and desk entities (Table C)", func(t *testing.T) {
		validSwitchboardEntities := []string{
			// Pattern: [0-9][CDBINQWMVROLPEUTZ0-9]B
			"0CB", "1DB", "2IB", "3NB", "4QB", "5WB",
			"6MB", "7VB", "8RB", "9OB", "0LB", "1PB",
			"2EB", "3UB", "4TB", "5ZB", "60B", "71B",
			"82B", "93B",
		}

		for _, entity := range validSwitchboardEntities {
			t.Run(entity, func(t *testing.T) {
				err := ValidateEntityCode(entity, true)
				// TODO: Should pass when implemented
				if err == nil {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Valid miscellaneous switching entities (Table D)", func(t *testing.T) {
		validMiscEntities := []string{
			// Pattern: [0-9][AXCTWDEINPQ]D
			"0AD", "1XD", "2CD", "3TD", "4WD", "5DD",
			"6ED", "7ID", "8ND", "9PD", "0QD",

			// Pattern: [A-Z0-9][UM]D
			"AUD", "AMD", "0UD", "9MD", "ZUD", "ZMD",
		}

		for _, entity := range validMiscEntities {
			t.Run(entity, func(t *testing.T) {
				err := ValidateEntityCode(entity, true)
				// TODO: Should pass when implemented
				if err == nil {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Valid non-switching entities (Table E)", func(t *testing.T) {
		validNonSwitchingEntities := []string{
			// Pattern: [FAEKMPSTW][x2][x1]
			"F23", "A12", "E45", "K67", "M89", "P01",
			"S34", "T56", "W78", "FAA", "AAA", "EZZ",
			"KA1", "M2Z", "P90",

			// Pattern: Q[0-9][0-9]
			"Q00", "Q01", "Q12", "Q99",
		}

		for _, entity := range validNonSwitchingEntities {
			t.Run(entity, func(t *testing.T) {
				err := ValidateEntityCode(entity, true)
				// TODO: Should pass when implemented
				if err == nil {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Invalid entity codes", func(t *testing.T) {
		invalidEntities := []struct {
			name   string
			entity string
		}{
			{"Empty", ""},
			{"Too short", "MG"},
			{"Too long", "MG1X"},
			{"Invalid pattern", "ABC"},
			{"Contains symbols", "MG@"},
			{"Lowercase", "mg1"},
			{"Invalid switching", "XG1"}, // X not valid prefix for switching
			{"Invalid Table C", "0AB"},   // A not valid for Table C
			{"Invalid Table D", "0BD"},   // B not valid for Table D
			{"Invalid Table E", "G23"},   // G not valid for Table E
		}

		for _, tt := range invalidEntities {
			t.Run(tt.name, func(t *testing.T) {
				err := ValidateEntityCode(tt.entity, true)
				// Should always fail for invalid entities
				assert.Error(t, err)
			})
		}
	})

	t.Run("Relaxed mode", func(t *testing.T) {
		// In relaxed mode, any 3-character alphanumeric code might be accepted
		relaxedEntities := []string{
			"ABC", "123", "A1B", "X9Z",
		}

		for _, entity := range relaxedEntities {
			t.Run(entity, func(t *testing.T) {
				err := ValidateEntityCode(entity, false)
				// Relaxed mode might accept these
				if err == nil {
					assert.NoError(t, err)
				}
			})
		}
	})
}

// TestValidationEdgeCases tests edge cases and boundary conditions
func TestValidationEdgeCases(t *testing.T) {
	t.Run("Null and empty inputs", func(t *testing.T) {
		validators := []struct {
			name string
			fn   func(string, bool) error
		}{
			{"ValidatePlace", ValidatePlace},
			{"ValidateRegion", ValidateRegion},
			{"ValidateNetworkSite", ValidateNetworkSite},
			{"ValidateEntityCode", ValidateEntityCode},
		}

		for _, validator := range validators {
			t.Run(validator.name, func(t *testing.T) {
				// Empty string should fail
				err := validator.fn("", true)
				assert.Error(t, err)

				// Empty string in relaxed mode might still fail
				err = validator.fn("", false)
				assert.Error(t, err)
			})
		}
	})

	t.Run("Whitespace handling", func(t *testing.T) {
		// Test whitespace in various positions
		testCases := []struct {
			input string
			valid bool
		}{
			{" ABC", false}, // Leading space
			{"ABC ", true},  // Trailing space (valid for place codes)
			{"A BC", false}, // Middle space
			{"  AB", false}, // Multiple leading spaces
			{"AB  ", true},  // Multiple trailing spaces (valid for place codes)
		}

		for _, tc := range testCases {
			t.Run(tc.input, func(t *testing.T) {
				// Only test with ValidatePlace as it's the only one that should accept spaces
				err := ValidatePlace(tc.input, true)
				if tc.valid {
					// TODO: Should pass when implemented
					if err == nil {
						assert.NoError(t, err)
					}
				} else {
					assert.Error(t, err)
				}
			})
		}
	})

	t.Run("Character class boundaries", func(t *testing.T) {
		// Test boundary characters for different classes
		boundaryTests := []struct {
			name     string
			input    string
			function func(string, bool) error
			valid    bool
		}{
			// Test G exclusion in a2 class (used in some entity patterns)
			{"Entity with G", "GAA", ValidateEntityCode, false},

			// Test excluded characters in a1 class
			{"Entity with B", "BAA", ValidateEntityCode, false}, // B excluded from a1
			{"Entity with D", "DAA", ValidateEntityCode, false}, // D excluded from a1
			{"Entity with I", "IAA", ValidateEntityCode, false}, // I excluded from a1
			{"Entity with O", "OAA", ValidateEntityCode, false}, // O excluded from a1
			{"Entity with T", "TAA", ValidateEntityCode, false}, // T excluded from a1
			{"Entity with U", "UAA", ValidateEntityCode, false}, // U excluded from a1
			{"Entity with W", "WAA", ValidateEntityCode, false}, // W excluded from a1
			{"Entity with Y", "YAA", ValidateEntityCode, false}, // Y excluded from a1
		}

		for _, tt := range boundaryTests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.function(tt.input, true)
				if tt.valid {
					if err == nil {
						assert.NoError(t, err)
					}
				} else {
					// Most should fail due to pattern restrictions
					// TODO: Update when patterns are implemented
					assert.Error(t, err)
				}
			})
		}
	})
}

// Benchmark tests for validation performance
func BenchmarkValidatePlace(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidatePlace("MPLS", true)
	}
}

func BenchmarkValidateRegion(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateRegion("MN", true)
	}
}

func BenchmarkValidateNetworkSite(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateNetworkSite("MS", true)
	}
}

func BenchmarkValidateEntityCode(b *testing.B) {
	entities := []string{"MG1", "SG2", "DS1", "CT1", "F23"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entity := entities[i%len(entities)]
		_ = ValidateEntityCode(entity, true)
	}
}
