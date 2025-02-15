package condition_action

type UniqueConditionAction struct {
	seenValues map[string]bool
}

func NewUniqueConditionAction() (a UniqueConditionAction) {
	a.seenValues = make(map[string]bool)
	return
}

func (a UniqueConditionAction) CheckCondition(input string, args []string) bool {
	if _, exists := a.seenValues[input]; exists {
		return false
	} else {
		// TODO: Set in the credential name as the value so we can log it if there is a collision
		a.seenValues[input] = true
		return true
	}
}
