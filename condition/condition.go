package condition

import (
	"fmt"
	"slices"
	"strings"
)

type Condition struct {
	Type ConditionType
	Args []string
}

func NewCondition(conditionString string) (condition Condition, err error) {
	index := strings.Index(conditionString, "(")
	conditionTypeString := conditionString
	if index != -1 {
		conditionTypeString = conditionString[0:index]
	}
	conditionType := ConditionTypeFromString(conditionTypeString)
	if !conditionType.IsValid() {
		err = fmt.Errorf("provided condition type [%s] is invalid", conditionTypeString)
		return
	}
	argumentsString := ""
	if index != -1 {
		argumentsString = conditionString[index+1:]

		closingBracketIndex := strings.LastIndex(argumentsString, ")")
		if closingBracketIndex != -1 {
			argumentsString = argumentsString[:len(argumentsString)-1]
		} else {
			err = fmt.Errorf("no closing bracket provided in condition [%s]", conditionString)
			return
		}
	}

	args, err := getArguments(argumentsString, conditionType.expectedArgsCount())
	if err != nil {
		fmt.Printf("Failed to create condition. Error: [%s]\n", err.Error())
		return
	}

	condition.Type = conditionType
	condition.Args = args

	return
}

func getArguments(arg string, expectedArgsCount uint) (args []string, err error) {
	split := strings.Split(arg, ",")
	// Remove empty
	split = slices.DeleteFunc(split, func(e string) bool { return e == "" })
	if len(split) > int(expectedArgsCount) {
		fmt.Printf("Too many arguments provided, expected [%d] but found [%d]. Only the first [%d] will be used.\n", expectedArgsCount, len(split), expectedArgsCount)
	} else if len(split) < int(expectedArgsCount) {
		err = fmt.Errorf("not enough arguments supplied for condition type")
		return
	}
	for i := range expectedArgsCount {
		args = append(args, strings.TrimSpace(split[i]))
	}
	return
}

func (c Condition) ApplyCondition(id string, input string) bool {
	return c.Type.getConditionAction(id).CheckCondition(input, c.Args)
}
