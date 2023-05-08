package util

import (
	"fmt"
	"strings"

	"github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/types"
	apiservercel "k8s.io/apiserver/pkg/admission/plugin/cel"
)

// ParseResource divides the resource into parent resource and subresource
func ParseResource(resource string) (parentResource, subresource string) {
	if strings.Contains(resource, "/") {
		splitResource := strings.Split(resource, "/")
		return splitResource[0], splitResource[1]
	}
	return resource, ""
}

// GetCompilationResults returns the results of compiling the selection preconditions
func GetCompilationResults(validations []types.Validation, perCallLimit uint64) ([]types.CombinedCompilationResult, error) {
	combinedCompilationResults := make([]types.CombinedCompilationResult, 0)
	for _, validation := range validations {
		validationExpressionCompilationResult := getValidationExpressionCompilationResult(validation.Expression, perCallLimit)
		if validationExpressionCompilationResult.Error != nil {
			return nil, fmt.Errorf("error while compiling CEL expression: '%s', error: '%v'", validation.Expression, *validationExpressionCompilationResult.Error)
		}

		combinedCompilationResult := types.CombinedCompilationResult{
			Key:                         validation.Key,
			ValidationCompilationResult: &validationExpressionCompilationResult,
		}

		if validation.MessageExpression != "" {
			messageExpressionCompilationResult := getMessageExpressionCompilationResult(validation.MessageExpression, perCallLimit)
			if messageExpressionCompilationResult.Error != nil {
				return nil, fmt.Errorf("error while compiling CEL expression: '%s', error: '%v'", validation.MessageExpression, *messageExpressionCompilationResult.Error)
			}
			combinedCompilationResult.MessageCompilationResult = &messageExpressionCompilationResult
		}
		combinedCompilationResults = append(combinedCompilationResults, combinedCompilationResult)
	}
	return combinedCompilationResults, nil
}

func getValidationExpressionCompilationResult(validationExpression string, perCallLimit uint64) apiservercel.CompilationResult {
	validationExpressionCompilationResult := apiservercel.CompileCELExpression(
		&types.ValidationExpressionAccessor{
			Expression: validationExpression,
		},
		apiservercel.OptionalVariableDeclarations{},
		perCallLimit,
	)
	return validationExpressionCompilationResult
}

func getMessageExpressionCompilationResult(messageExpression string, perCallLimit uint64) apiservercel.CompilationResult {
	messageExpressionCompilationResult := apiservercel.CompileCELExpression(
		&types.MessageExpressionAccessor{
			MessageExpression: messageExpression,
		},
		apiservercel.OptionalVariableDeclarations{},
		perCallLimit,
	)
	return messageExpressionCompilationResult
}
