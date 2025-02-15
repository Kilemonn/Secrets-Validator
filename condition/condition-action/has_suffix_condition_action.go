package condition_action

import "strings"

type HasSuffixConditionAction struct{}

func (a HasSuffixConditionAction) CheckCondition(input string, args []string) bool {
	return strings.HasSuffix(input, args[0])
}
