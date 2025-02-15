package constraint

import (
	"github.com/Kilemonn/Secrets-Validator/condition"
	"github.com/Kilemonn/Secrets-Validator/pattern"
)

type Constraint struct {
	Name      string
	Pattern   pattern.Pattern
	Condition condition.Condition
}

func NewConstraint(name string, patternString string, conditionString string) (constraint Constraint, err error) {
	conditionObj, err := condition.NewCondition(conditionString)
	if err != nil {
		return
	}
	patternObj, err := pattern.NewPattern(patternString)
	if err != nil {
		return
	}
	return Constraint{
		Name:      name,
		Pattern:   patternObj,
		Condition: conditionObj,
	}, nil
}
