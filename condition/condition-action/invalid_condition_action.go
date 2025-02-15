package condition_action

type InvalidConditionAction struct{}

func (a InvalidConditionAction) CheckCondition(input string, args []string) bool {
	return false
}
