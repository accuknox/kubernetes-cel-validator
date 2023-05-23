package validation

import (
	"fmt"

	"github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/types"
)

func ValidatePreconditionData(precondition *types.KubernetesResourcePrecondition) error {
	if precondition.Name == "" {
		return fmt.Errorf("no name specified, name must be specified")
	}
	if precondition.MatchResources == nil {
		return fmt.Errorf("no match resources specified, match resources must be specified")
	}

	resourceRules := precondition.MatchResources.ResourceRules
	if len(resourceRules) == 0 {
		return fmt.Errorf("no resource rules specified, at least one resource rule must be specified")
	}

	validations := precondition.Validations
	if len(validations) == 0 {
		return fmt.Errorf("no validations specified, at least one validation must be specified")
	}

	validationMap := make(map[string]bool, len(validations))

	for _, validation := range validations {
		// Key is optional
		if validation.Key != "" {
			if _, ok := validationMap[validation.Key]; !ok {
				validationMap[validation.Key] = true
			} else {
				return fmt.Errorf("duplicate validation key: %s", validation.Key)
			}
		}
	}

	selectionPreconditions := precondition.MatchResources.SelectionPreconditions
	selectionPreconditionsMap := make(map[string]bool, len(selectionPreconditions))

	for _, selectionPrecondition := range selectionPreconditions {
		// Key is optional
		if selectionPrecondition.Key != "" {
			if _, ok := selectionPreconditionsMap[selectionPrecondition.Key]; !ok {
				selectionPreconditionsMap[selectionPrecondition.Key] = true
			} else {
				return fmt.Errorf("duplicate selection precondition key: %s", selectionPrecondition.Key)
			}
		}
	}

	return nil
}
