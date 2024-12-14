package pattern

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatternMatches(t *testing.T) {
	var cases = []struct {
		regex        string
		expectsError bool
		successTests []string
		failTests    []string
	}{
		{"ALL", false, []string{"some random !@#$!@$12471o31 text"}, []string{}},                             // Using "ALL"
		{"db-connection-string-.+", false, []string{"db-connection-string-998234"}, []string{"fail please"}}, // Valid regex with success and failure mathcing checks
		{"invalid-pattern*+", true, []string{}, []string{}},                                                  // "Failing to initialise regex"
	}

	for _, c := range cases {
		pattern, err := NewPattern(c.regex)
		if c.expectsError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			for _, test := range c.successTests {
				require.True(t, pattern.Matches(test))
			}
			for _, test := range c.failTests {
				require.False(t, pattern.Matches(test))
			}
		}
	}
}
