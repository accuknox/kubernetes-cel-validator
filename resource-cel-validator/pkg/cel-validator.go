package pkg

import (
	"github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/types"
	"github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/validation"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

func GetKubernetesResourcePreconditionResult(precondition *types.KubernetesResourcePrecondition, config *rest.Config) (bool, []types.ValidationFailure, types.ValidationFailure, error) {
	err := validation.ValidatePreconditionData(precondition)
	if err != nil {
		return false, nil, types.ValidationFailure{}, err
	}

	var dynamicClientSet *dynamic.DynamicClient
	dynamicClientSet, err = dynamic.NewForConfig(config)
	if err != nil {
		return false, nil, types.ValidationFailure{}, err
	}

	resourceList, err := validation.GetResources(precondition.MatchResources, dynamicClientSet)
	if err != nil {
		return false, nil, types.ValidationFailure{}, err
	}

	filteredResources, selectionPreconditionValidationFailures, err := validation.FilterResourcesOnSelectionPreconditions(resourceList, precondition.MatchResources.SelectionPreconditions)
	if err != nil {
		return false, nil, types.ValidationFailure{}, err
	}

	validationResult, validationFailure, err := validation.GetValidationResults(filteredResources, precondition.Validations)
	return validationResult, selectionPreconditionValidationFailures, validationFailure, err
}
