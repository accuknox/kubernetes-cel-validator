package validation

import (
	"fmt"

	"github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/types"
	"github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/util"
	goceltypes "github.com/google/cel-go/common/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// FilterResourcesOnSelectionPreconditions filters resources based on selectionPreconditions.
func FilterResourcesOnSelectionPreconditions(resources []unstructured.Unstructured, selectionPreconditions []types.Validation) ([]unstructured.Unstructured, []types.ValidationFailure, error) {
	filteredResources := make([]unstructured.Unstructured, 0)
	selectionPreconditionValidationFailures := make([]types.ValidationFailure, 0)
	compilationResults, err := util.GetCompilationResults(selectionPreconditions, PerCallLimit)
	if err != nil {
		return nil, nil, err
	}

	for _, resource := range resources {
		evaluationPassed := true
		for _, compilationResult := range compilationResults {
			validationCompilationResult := compilationResult.ValidationCompilationResult
			validationExpressionProgram := validationCompilationResult.Program
			val, _, err := validationExpressionProgram.Eval(
				map[string]any{
					"object": resource.Object,
				},
			)
			if err != nil {
				return nil, nil, fmt.Errorf("error while evaluating CEL expression: '%s', error: '%v'",
					validationCompilationResult.ExpressionAccessor.GetExpression(), err)
			}
			if val != goceltypes.True {
				selectionPreconditionValidationFailure := types.ValidationFailure{
					Key:      compilationResult.Key,
					Resource: resource,
				}

				messageCompilationResult := compilationResult.MessageCompilationResult
				if messageCompilationResult != nil {
					messageExpressionProgram := messageCompilationResult.Program
					messageVal, _, err := messageExpressionProgram.Eval(
						map[string]any{
							"object": resource.Object,
						},
					)
					if err != nil {
						return nil, nil, fmt.Errorf("error while evaluating CEL expression: '%s', error: '%v'",
							messageCompilationResult.ExpressionAccessor.GetExpression(), err)
					}
					selectionPreconditionValidationFailure.Message = messageVal.Value().(string)
				}

				selectionPreconditionValidationFailures = append(selectionPreconditionValidationFailures, selectionPreconditionValidationFailure)
				evaluationPassed = false
				break
			}
		}
		if evaluationPassed {
			filteredResources = append(filteredResources, resource)
		}
	}
	return filteredResources, selectionPreconditionValidationFailures, nil
}
