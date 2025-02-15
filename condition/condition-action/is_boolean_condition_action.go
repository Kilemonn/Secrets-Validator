package condition_action

import "strconv"

type IsBooleanConditionAction struct{}

func (a IsBooleanConditionAction) CheckCondition(input string, args []string) bool {
	_, err := strconv.ParseBool(input)
	return err == nil
}
