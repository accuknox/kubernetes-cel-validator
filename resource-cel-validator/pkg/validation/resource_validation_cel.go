package validation

import (
	"fmt"

	"github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/types"
	"github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/util"
	goceltypes "github.com/google/cel-go/common/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func GetValidationResults(resources []unstructured.Unstructured, validations []types.Validation) (bool, types.ValidationFailure, error) {
	compilationResults, err := util.GetCompilationResults(validations, PerCallLimit)
	if err != nil {
		return false, types.ValidationFailure{}, err
	}

	for _, resource := range resources {
		for _, compilationResult := range compilationResults {
			validationCompilationResult := compilationResult.ValidationCompilationResult
			validationExpressionProgram := validationCompilationResult.Program
			val, _, err := validationExpressionProgram.Eval(
				map[string]any{
					"object": resource.Object,
				},
			)
			if err != nil {
				return false, types.ValidationFailure{},
					fmt.Errorf("error while evaluating CEL expression: '%s', error: '%v'",
						validationCompilationResult.ExpressionAccessor.GetExpression(), err)
			}
			if val != goceltypes.True {
				validationFailure := types.ValidationFailure{
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
						return false, types.ValidationFailure{},
							fmt.Errorf("error while evaluating CEL expression: '%s', error: '%v'",
								messageCompilationResult.ExpressionAccessor.GetExpression(), err)
					}
					validationFailure.Message = messageVal.Value().(string)
				}

				return false, validationFailure, nil
			}
		}
	}
	return true, types.ValidationFailure{}, nil
}
