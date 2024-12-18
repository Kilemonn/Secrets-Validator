package condition

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetArguments(t *testing.T) {
	var cases = []struct {
		argString     string
		expectedCount uint
		expectsError  bool
		expected      []string
	}{
		{"", 0, false, []string{}},                               // Empty args, non requested
		{"", 1, true, []string{}},                                // Empty args, 1 requested
		{"arg", 0, false, []string{}},                            // 1 arg, with 0 requested
		{"Sometest-1237-", 1, false, []string{"Sometest-1237-"}}, // One arg with no "," all should be returned
		{"point, test", 2, false, []string{"point", "test"}},     // 2 Args requested, with spaces
	}

	for _, c := range cases {
		result, err := getArguments(c.argString, c.expectedCount)
		if c.expectsError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, len(c.expected), len(result))
			for i := range c.expected {
				require.Equal(t, c.expected[i], result[i])
			}
		}
	}
}

func TestNewCondition(t *testing.T) {
	var cases = []struct {
		conditionString string
		expectError     bool
		expectedType    ConditionType
		expectedArgs    []string
	}{
		{"Unique", false, ConditionTypeUnique, []string(nil)},        // 0 arg condition
		{"Unique(check)", false, ConditionTypeUnique, []string(nil)}, // 0 arg condition with args
		{"NotUnique", true, ConditionTypeUnique, []string(nil)},      // 0 arg but not matching any type

		{"HasPrefix", true, ConditionTypeHasPrefix, []string(nil)},                 // 1 arg required but not provided
		{"HasPrefix(test)", false, ConditionTypeHasPrefix, []string{"test"}},       // 1 arg required and provided
		{"HasPrefix(test, text)", false, ConditionTypeHasPrefix, []string{"test"}}, // 1 arg required and 2 provided

		{"HasPrefix( , text)", false, ConditionTypeHasPrefix, []string{""}}, // 1 arg required and 2 provided but first is a space

		{"HasPrefix(, text)", false, ConditionTypeHasPrefix, []string{"text"}}, // 1 arg required and 2 provided but first is empty string

		{"HasSuffix(test", true, ConditionTypeHasPrefix, []string(nil)}, // 1 arg required but no closing bracket
		{"HasSuffixtest)", true, ConditionTypeHasPrefix, []string(nil)}, // 1 arg required but no open bracket
	}

	for _, c := range cases {
		condition, err := NewCondition(c.conditionString)
		if c.expectError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, c.expectedType, condition.Type)
			require.Equal(t, c.expectedType.expectedArgsCount(), uint(len(condition.Args)))
			require.Equal(t, c.expectedArgs, condition.Args)
		}
	}
}

func TestApplyCondition(t *testing.T) {
	var cases = []struct {
		conditionString string
		id              string
		expectError     bool
		successInputs   []string
		failInputs      []string
	}{
		{"IsNumeric", "is-numeric", false, []string{"-19", "40", "9999999999999"}, []string{"wow", "-", "!@!#$!"}},
		{"IsBoolean", "is-boolean", false, []string{"true", "false"}, []string{"yes", "-", "222159", "!@!#$!"}},
		{"Unique", "is-unique1", false, []string{"unique1", "unique2", "unique3"}, []string{"unique1", "unique2", "unique3"}},

		// Duplicate the unique entry with the same ID to make sure that the same unique instance is used
		{"Unique", "is-unique1", false, []string{"unique4", "unique5", "unique6"}, []string{"unique1", "unique2", "unique3", "unique4", "unique5", "unique6"}},
		// Using a different ID should give us a different validator and not collide with the other unique instances
		{"Unique", "is-unique2", false, []string{"unique1", "unique2", "unique3"}, []string{"unique1", "unique2", "unique3"}},
		{"HasPrefix(test-)", "has-prefix", false, []string{"test-scenario1", "test-case1", "test-please"}, []string{"not-test-", "ttest-failed"}},
		{"HasSuffix(-test)", "has-suffix", false, []string{"some-test", "another-test"}, []string{"not-test-", "failing-tests"}},

		// Make sure an invalid constraint is forced to validate false to everything
		{"SomethingInvalid(arg1, arg2)", "invalid", true, []string{}, []string{"test test test", "1237532123", "true", "$!&@#($)"}},

		{"Matches(\\d+)", "matches-regex", false, []string{"1234", "1", "testwith number 1"}, []string{"test", "no numbers++--"}},
		{"Matches([)", "invalid-regex", false, []string{"test"}, []string{"?"}},
	}

	for _, c := range cases {
		condition, err := NewCondition(c.conditionString)
		if c.expectError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		for _, successTest := range c.successInputs {
			require.True(t, condition.ApplyCondition(c.id, successTest), c.id+" - "+successTest)
		}
		for _, failTest := range c.failInputs {
			require.False(t, condition.ApplyCondition(c.id, failTest), c.id+" - "+failTest)
		}
	}
}
