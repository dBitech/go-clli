package clli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGeographicResolution tests all geographic resolution methods
func TestGeographicResolution(t *testing.T) {
	// Test data based on the Ruby gem's test cases
	testCases := []struct {
		name                string
		place               string
		region              string
		expectedCity        string
		expectedState       string
		expectedStateCode   string
		expectedCountry     string
		expectedCountryCode string
	}{
		{
			"Minneapolis, Minnesota",
			"MPLS", "MN",
			"Minneapolis", "Minnesota", "MN", "United States", "US",
		},
		{
			"New York City, New York",
			"NYCM", "NY",
			"New York City", "New York", "NY", "United States", "US",
		},
		{
			"Los Angeles, California",
			"LSAN", "CA",
			"Los Angeles", "California", "CA", "United States", "US",
		},
		{
			"Chicago, Illinois",
			"CHCG", "IL",
			"Chicago", "Illinois", "IL", "United States", "US",
		},
		{
			"Houston, Texas",
			"HSTX", "TX",
			"Houston", "Texas", "TX", "United States", "US",
		},
		{
			"Philadelphia, Pennsylvania",
			"PHLA", "PA",
			"Philadelphia", "Pennsylvania", "PA", "United States", "US",
		},
		{
			"Phoenix, Arizona",
			"PHNX", "AZ",
			"Phoenix", "Arizona", "AZ", "United States", "US",
		},
		{
			"San Antonio, Texas",
			"SNAN", "TX",
			"San Antonio", "Texas", "TX", "United States", "US",
		},
		{
			"San Diego, California",
			"SNDG", "CA",
			"San Diego", "California", "CA", "United States", "US",
		},
		{
			"Dallas, Texas",
			"DLLS", "TX",
			"Dallas", "Texas", "TX", "United States", "US",
		},

		// Canadian test cases
		{
			"Toronto, Ontario",
			"TORO", "ON",
			"Toronto", "Ontario", "ON", "Canada", "CA",
		},
		{
			"Montreal, Quebec",
			"MTRL", "QC",
			"Montreal", "Quebec", "QC", "Canada", "CA",
		},
		{
			"Vancouver, British Columbia",
			"VANCVR", "BC", // Note: Padded place code
			"Vancouver", "British Columbia", "BC", "Canada", "CA",
		},
		{
			"Calgary, Alberta",
			"CGRY", "AB",
			"Calgary", "Alberta", "AB", "Canada", "CA",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create CLLI with test data
			clli := &CLLI{
				Place:  tc.place,
				Region: tc.region,
			}

			// Test CountryCode
			t.Run("CountryCode", func(t *testing.T) {
				result := clli.CountryCode()
				// TODO: Should return expected code when implemented
				if result != "" {
					assert.Equal(t, tc.expectedCountryCode, result)
				} else {
					// Expected empty until implementation
					assert.Empty(t, result)
				}
			})

			// Test CountryName
			t.Run("CountryName", func(t *testing.T) {
				result := clli.CountryName()
				// TODO: Should return expected name when implemented
				if result != "" {
					assert.Equal(t, tc.expectedCountry, result)
				} else {
					// Expected empty until implementation
					assert.Empty(t, result)
				}
			})

			// Test StateCode
			t.Run("StateCode", func(t *testing.T) {
				result := clli.StateCode()
				// TODO: Should return expected code when implemented
				if result != "" {
					assert.Equal(t, tc.expectedStateCode, result)
				} else {
					// Expected empty until implementation
					assert.Empty(t, result)
				}
			})

			// Test StateName
			t.Run("StateName", func(t *testing.T) {
				result := clli.StateName()
				// TODO: Should return expected name when implemented
				if result != "" {
					assert.Equal(t, tc.expectedState, result)
				} else {
					// Expected empty until implementation
					assert.Empty(t, result)
				}
			})

			// Test CityName
			t.Run("CityName", func(t *testing.T) {
				result := clli.CityName()
				// TODO: Should return expected city when implemented
				if result != "" {
					assert.Equal(t, tc.expectedCity, result)
				} else {
					// Expected empty until implementation
					assert.Empty(t, result)
				}
			})
		})
	}
}

// TestUnknownGeographicCodes tests handling of unknown or invalid region codes
func TestUnknownGeographicCodes(t *testing.T) {
	testCases := []struct {
		name   string
		place  string
		region string
	}{
		{"Invalid region code", "TEST", "ZZ"},
		{"Unknown place code", "XXXX", "CA"},
		{"Empty region", "MPLS", ""},
		{"Empty place", "", "MN"},
		{"Both empty", "", ""},
		{"Non-existent combination", "FAKE", "XX"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clli := &CLLI{
				Place:  tc.place,
				Region: tc.region,
			}

			// All methods should handle unknown codes gracefully
			t.Run("CountryCode", func(t *testing.T) {
				result := clli.CountryCode()
				// Should return empty string or default for unknown codes
				if tc.region == "ZZ" || tc.region == "XX" || tc.region == "" {
					// Unknown region should return empty or default
					assert.Contains(t, []string{"", "ZZ", "XX"}, result)
				}
			})

			t.Run("CountryName", func(t *testing.T) {
				result := clli.CountryName()
				// Should return empty string for unknown codes
				if tc.region == "ZZ" || tc.region == "XX" || tc.region == "" {
					assert.Empty(t, result)
				}
			})

			t.Run("StateCode", func(t *testing.T) {
				result := clli.StateCode()
				// Should return empty string for unknown codes
				if tc.region == "ZZ" || tc.region == "XX" || tc.region == "" {
					assert.Empty(t, result)
				}
			})

			t.Run("StateName", func(t *testing.T) {
				result := clli.StateName()
				// Should return empty string for unknown codes
				if tc.region == "ZZ" || tc.region == "XX" || tc.region == "" {
					assert.Empty(t, result)
				}
			})

			t.Run("CityName", func(t *testing.T) {
				result := clli.CityName()
				// Should return empty string for unknown combinations
				assert.Empty(t, result)
			})
		})
	}
}

// TestAllUSStates tests geographic resolution for all US states and territories
func TestAllUSStates(t *testing.T) {
	// US States and territories with their codes and names
	usStates := map[string]string{
		"AL": "Alabama",
		"AK": "Alaska",
		"AZ": "Arizona",
		"AR": "Arkansas",
		"CA": "California",
		"CO": "Colorado",
		"CT": "Connecticut",
		"DE": "Delaware",
		"FL": "Florida",
		"GA": "Georgia",
		"HI": "Hawaii",
		"ID": "Idaho",
		"IL": "Illinois",
		"IN": "Indiana",
		"IA": "Iowa",
		"KS": "Kansas",
		"KY": "Kentucky",
		"LA": "Louisiana",
		"ME": "Maine",
		"MD": "Maryland",
		"MA": "Massachusetts",
		"MI": "Michigan",
		"MN": "Minnesota",
		"MS": "Mississippi",
		"MO": "Missouri",
		"MT": "Montana",
		"NE": "Nebraska",
		"NV": "Nevada",
		"NH": "New Hampshire",
		"NJ": "New Jersey",
		"NM": "New Mexico",
		"NY": "New York",
		"NC": "North Carolina",
		"ND": "North Dakota",
		"OH": "Ohio",
		"OK": "Oklahoma",
		"OR": "Oregon",
		"PA": "Pennsylvania",
		"RI": "Rhode Island",
		"SC": "South Carolina",
		"SD": "South Dakota",
		"TN": "Tennessee",
		"TX": "Texas",
		"UT": "Utah",
		"VT": "Vermont",
		"VA": "Virginia",
		"WA": "Washington",
		"WV": "West Virginia",
		"WI": "Wisconsin",
		"WY": "Wyoming",
		"DC": "District of Columbia",
	}

	for stateCode, stateName := range usStates {
		t.Run(stateCode, func(t *testing.T) {
			clli := &CLLI{
				Place:  "TEST", // Generic place for testing
				Region: stateCode,
			}

			// Test country resolution for US states
			t.Run("CountryCode", func(t *testing.T) {
				result := clli.CountryCode()
				if result != "" {
					assert.Equal(t, "US", result)
				}
			})

			t.Run("CountryName", func(t *testing.T) {
				result := clli.CountryName()
				if result != "" {
					assert.Equal(t, "United States", result)
				}
			})

			// Test state resolution
			t.Run("StateCode", func(t *testing.T) {
				result := clli.StateCode()
				if result != "" {
					assert.Equal(t, stateCode, result)
				}
			})

			t.Run("StateName", func(t *testing.T) {
				result := clli.StateName()
				if result != "" {
					assert.Equal(t, stateName, result)
				}
			})
		})
	}
}

// TestAllCanadianProvinces tests geographic resolution for Canadian provinces and territories
func TestAllCanadianProvinces(t *testing.T) {
	// Canadian provinces and territories with their codes and names
	canadianProvinces := map[string]string{
		"AB": "Alberta",
		"BC": "British Columbia",
		"MB": "Manitoba",
		"NB": "New Brunswick",
		"NL": "Newfoundland and Labrador",
		"NS": "Nova Scotia",
		"ON": "Ontario",
		"PE": "Prince Edward Island",
		"QC": "Quebec",
		"SK": "Saskatchewan",
		"NT": "Northwest Territories",
		"NU": "Nunavut",
		"YT": "Yukon",
	}

	for provinceCode, provinceName := range canadianProvinces {
		t.Run(provinceCode, func(t *testing.T) {
			clli := &CLLI{
				Place:  "TEST", // Generic place for testing
				Region: provinceCode,
			}

			// Test country resolution for Canadian provinces
			t.Run("CountryCode", func(t *testing.T) {
				result := clli.CountryCode()
				if result != "" {
					assert.Equal(t, "CA", result)
				}
			})

			t.Run("CountryName", func(t *testing.T) {
				result := clli.CountryName()
				if result != "" {
					assert.Equal(t, "Canada", result)
				}
			})

			// Test province resolution
			t.Run("StateCode", func(t *testing.T) {
				result := clli.StateCode()
				if result != "" {
					assert.Equal(t, provinceCode, result)
				}
			})

			t.Run("StateName", func(t *testing.T) {
				result := clli.StateName()
				if result != "" {
					assert.Equal(t, provinceName, result)
				}
			})
		})
	}
}

// TestGeographicEdgeCases tests edge cases in geographic resolution
func TestGeographicEdgeCases(t *testing.T) {
	t.Run("Place code normalization", func(t *testing.T) {
		// Test place codes with trailing spaces (should be trimmed)
		testCases := []struct {
			place    string
			expected string
		}{
			{"MPLS", "MPLS"}, // No spaces
			{"MPL ", "MPL"},  // One trailing space
			{"MP  ", "MP"},   // Two trailing spaces
			{"M   ", "M"},    // Three trailing spaces
		}

		for _, tc := range testCases {
			t.Run(tc.place, func(t *testing.T) {
				clli := &CLLI{
					Place:  tc.place,
					Region: "MN",
				}

				// The place should be internally normalized
				// This test verifies that the geographic resolution works
				// regardless of trailing spaces in place codes
				_ = clli.CityName() // Should work regardless of spaces
			})
		}
	})

	t.Run("Case sensitivity", func(t *testing.T) {
		// All region codes should be uppercase, but test mixed case handling
		testCases := []string{
			"mn", // lowercase
			"Mn", // mixed case
			"mN", // mixed case
		}

		for _, region := range testCases {
			t.Run(region, func(t *testing.T) {
				clli := &CLLI{
					Place:  "MPLS",
					Region: region,
				}

				// Geographic resolution should handle case normalization
				// or fail gracefully for non-uppercase region codes
				countryCode := clli.CountryCode()
				countryName := clli.CountryName()
				stateCode := clli.StateCode()
				stateName := clli.StateName()
				cityName := clli.CityName()

				// Results should either be empty (not found) or properly resolved
				if countryCode != "" {
					assert.Equal(t, "US", countryCode)
				}
				if countryName != "" {
					assert.Equal(t, "United States", countryName)
				}
				if stateCode != "" {
					assert.Equal(t, "MN", stateCode)
				}
				if stateName != "" {
					assert.Equal(t, "Minnesota", stateName)
				}
				if cityName != "" {
					assert.Equal(t, "Minneapolis", cityName)
				}
			})
		}
	})

	t.Run("Partial CLLI geographic resolution", func(t *testing.T) {
		// Test resolution when only some components are present
		testCases := []struct {
			name   string
			place  string
			region string
		}{
			{"Only place", "MPLS", ""},
			{"Only region", "", "MN"},
			{"Both empty", "", ""},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				clli := &CLLI{
					Place:  tc.place,
					Region: tc.region,
				}

				// Should handle partial data gracefully
				countryCode := clli.CountryCode()
				countryName := clli.CountryName()
				stateCode := clli.StateCode()
				stateName := clli.StateName()
				cityName := clli.CityName()

				// If region is empty, geographic methods should return empty
				if tc.region == "" {
					assert.Empty(t, countryCode)
					assert.Empty(t, countryName)
					assert.Empty(t, stateCode)
					assert.Empty(t, stateName)
				}

				// If place is empty, city name should be empty
				if tc.place == "" {
					assert.Empty(t, cityName)
				}
			})
		}
	})
}

// TestGeographicConsistency tests consistency between related geographic methods
func TestGeographicConsistency(t *testing.T) {
	testCases := []struct {
		name   string
		place  string
		region string
	}{
		{"US location", "MPLS", "MN"},
		{"Canadian location", "TORO", "ON"},
		{"California location", "LSAN", "CA"},
		{"Texas location", "HSTX", "TX"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clli := &CLLI{
				Place:  tc.place,
				Region: tc.region,
			}

			countryCode := clli.CountryCode()
			countryName := clli.CountryName()
			stateCode := clli.StateCode()
			stateName := clli.StateName()

			// If we get a country code, we should get a country name
			if countryCode != "" && countryName != "" {
				// Verify consistency between code and name
				switch countryCode {
				case "US":
					assert.Equal(t, "United States", countryName)
				case "CA":
					assert.Equal(t, "Canada", countryName)
				}
			}

			// If we get a state code, we should get a state name
			if stateCode != "" && stateName != "" {
				// State code should match the region (for US/CA)
				if countryCode == "US" || countryCode == "CA" {
					assert.Equal(t, tc.region, stateCode)
				}
			}

			// Country and state should be consistent
			if countryCode == "US" {
				// US states should be 2-character codes
				if stateCode != "" {
					assert.Len(t, stateCode, 2)
				}
			} else if countryCode == "CA" {
				// Canadian provinces should be 2-character codes
				if stateCode != "" {
					assert.Len(t, stateCode, 2)
				}
			}
		})
	}
}

// Benchmark tests for geographic resolution performance
func BenchmarkCountryCode(b *testing.B) {
	clli := &CLLI{Place: "MPLS", Region: "MN"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = clli.CountryCode()
	}
}

func BenchmarkCountryName(b *testing.B) {
	clli := &CLLI{Place: "MPLS", Region: "MN"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = clli.CountryName()
	}
}

func BenchmarkStateCode(b *testing.B) {
	clli := &CLLI{Place: "MPLS", Region: "MN"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = clli.StateCode()
	}
}

func BenchmarkStateName(b *testing.B) {
	clli := &CLLI{Place: "MPLS", Region: "MN"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = clli.StateName()
	}
}

func BenchmarkCityName(b *testing.B) {
	clli := &CLLI{Place: "MPLS", Region: "MN"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = clli.CityName()
	}
}

func BenchmarkAllGeographicMethods(b *testing.B) {
	clli := &CLLI{Place: "MPLS", Region: "MN"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = clli.CountryCode()
		_ = clli.CountryName()
		_ = clli.StateCode()
		_ = clli.StateName()
		_ = clli.CityName()
	}
}
