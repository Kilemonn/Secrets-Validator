package condition_action

import "strings"

type HasPrefixConditionAction struct{}

func (a HasPrefixConditionAction) CheckCondition(input string, args []string) bool {
	return strings.HasPrefix(input, args[0])
}
