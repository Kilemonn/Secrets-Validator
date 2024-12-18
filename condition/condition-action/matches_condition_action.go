package condition_action

import (
	"fmt"
	"regexp"
)

type MatchesConditionAction struct {
	regex *regexp.Regexp
}

func NewMatchesConditionAction(id string, pattern string) MatchesConditionAction {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Failed to initialise [Matches] condition's pattern with ID [%s], pattern [%s] with error [%s]. All condition checks will fail.", id, pattern, err.Error())
	}

	return MatchesConditionAction{
		regex: regex,
	}
}

func (a MatchesConditionAction) CheckCondition(input string, args []string) bool {
	if a.regex == nil {
		return false
	}
	return a.regex.Match([]byte(input))
}
