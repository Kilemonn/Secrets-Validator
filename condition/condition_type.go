package condition

import (
	"slices"
	"strings"

	condition_action "github.com/Kilemonn/Secrets-Validator/condition/condition-action"
)

type ConditionType uint

const (
	ConditionTypeInvalid   ConditionType = iota
	ConditionTypeUnique    ConditionType = iota
	ConditionTypeHasPrefix ConditionType = iota
	ConditionTypeHasSuffix ConditionType = iota
	ConditionTypeIsNumeric ConditionType = iota
	ConditionTypeIsBoolean ConditionType = iota
	ConditionTypeMatches   ConditionType = iota
)

func ConditionTypeFromString(input string) ConditionType {
	index := slices.Index(conditionTypeStrings(), strings.ToLower(input))
	if index == -1 {
		return ConditionTypeInvalid
	} else {
		return ConditionType(index)
	}
}

func conditionTypeStrings() []string {
	return []string{"invalid", "unique", "hasprefix", "hassuffix", "isnumeric", "isboolean", "matches"}
}

func (c ConditionType) IsValid() bool {
	return c != ConditionTypeInvalid
}

func (c ConditionType) expectedArgsCount() uint {
	args := []uint{0, 0, 1, 1, 0, 0, 1}
	return args[uint(c)]
}

func (c ConditionType) conditionActions(id string, pattern string) []condition_action.ConditionAction {
	return []condition_action.ConditionAction{
		condition_action.InvalidConditionAction{},
		condition_action.GetUniqueInstance(id),
		condition_action.HasPrefixConditionAction{},
		condition_action.HasSuffixConditionAction{},
		condition_action.IsNumericConditionAction{},
		condition_action.IsBooleanConditionAction{},
		condition_action.GetMatchesInstance(id, pattern),
	}
}

func (c ConditionType) getConditionAction(id string, pattern string) condition_action.ConditionAction {
	return c.conditionActions(id, pattern)[uint(c)]
}
