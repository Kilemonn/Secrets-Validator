package condition_action

var (
	uniqueActionsMap map[string]UniqueConditionAction = make(map[string]UniqueConditionAction)

	matchesActionsMap map[string]MatchesConditionAction = make(map[string]MatchesConditionAction)
)

type ConditionAction interface {
	CheckCondition(input string, args []string) bool
}

func GetUniqueInstance(id string) UniqueConditionAction {
	if entry, exists := uniqueActionsMap[id]; exists {
		return entry
	} else {
		uniqueAction := NewUniqueConditionAction()
		uniqueActionsMap[id] = uniqueAction
		return uniqueAction
	}
}

func GetMatchesInstance(id string, pattern string) MatchesConditionAction {
	if entry, exists := matchesActionsMap[id]; exists {
		return entry
	} else {
		matchesAction := NewMatchesConditionAction(id, pattern)
		matchesActionsMap[id] = matchesAction
		return matchesAction
	}
}
