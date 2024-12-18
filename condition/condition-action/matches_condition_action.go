package condition_action

import "regexp"

type MatchesConditionAction struct {
	regex *regexp.Regexp
}

func NewMatchesConditionAction(pattern string) MatchesConditionAction {
	regex := regexp.MustCompile(pattern)

	return MatchesConditionAction{
		regex: regex,
	}
}

func (a MatchesConditionAction) CheckCondition(input string, args []string) bool {
	return a.regex.Match([]byte(input))
}
