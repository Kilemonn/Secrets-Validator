package condition_action

import "strconv"

type IsNumericConditionAction struct{}

func (a IsNumericConditionAction) CheckCondition(input string, args []string) bool {
	_, err := strconv.ParseInt(input, 10, 64)
	return err == nil
}
