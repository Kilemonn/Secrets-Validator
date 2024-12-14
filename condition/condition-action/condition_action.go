package condition_action

var (
	// TODO: Create a new instance for each different constraint definition that uses Unique
	// Need to create and hold this variable, since it's state needs to be retained
	uniqueActionsMap map[string]UniqueConditionAction = make(map[string]UniqueConditionAction)
)

func GetUniqueInstance(id string) UniqueConditionAction {
	if entry, exists := uniqueActionsMap[id]; exists {
		return entry
	} else {
		uniqueAction := NewUniqueConditionAction()
		uniqueActionsMap[id] = uniqueAction
		return uniqueAction
	}
}

type ConditionAction interface {
	CheckCondition(input string, args []string) bool
}
