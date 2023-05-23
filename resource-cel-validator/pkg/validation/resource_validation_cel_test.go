package validation

import (
	"testing"

	"github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/types"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestGetValidationResults(t *testing.T) {
	validations := []types.Validation{
		{
			Expression: "object.metadata.namespace == 'validation-test-ns'",
		},
		{
			Key:               "nginx-name",
			Expression:        "object.metadata.name == 'nginx'",
			MessageExpression: "object.metadata.name + ' is not nginx'",
		},
	}

	result, validationFailure, err := GetValidationResults([]unstructured.Unstructured{deploymentResource}, validations)
	assert.Equal(t, result, true)
	assert.Equal(t, validationFailure, types.ValidationFailure{})
	assert.Nil(t, err)

	result, validationFailure, err = GetValidationResults([]unstructured.Unstructured{serviceResource}, validations)
	assert.Equal(t, result, false)
	expectedValidationFailure := types.ValidationFailure{
		Key:      "",
		Message:  "",
		Resource: serviceResource,
	}
	assert.Equal(t, validationFailure, expectedValidationFailure)
	assert.Nil(t, err)

	validations = []types.Validation{
		{
			Expression:        "object.metadata.namespace == 'validation-test-ns'",
			MessageExpression: "object.metadata.namespace + ' is not validation-test-ns'",
		},
		{
			Key:               "nginx-name",
			Expression:        "object.metadata.name == 'nginx'",
			MessageExpression: "object.metadata.name + ' is not nginx'",
		},
	}

	result, validationFailure, err = GetValidationResults([]unstructured.Unstructured{deploymentResource, serviceResource}, validations)
	assert.Equal(t, result, false)
	expectedValidationFailure = types.ValidationFailure{
		Key:      "",
		Message:  "validation-test-ns-2 is not validation-test-ns",
		Resource: serviceResource,
	}
	assert.Equal(t, validationFailure, expectedValidationFailure)
	assert.Nil(t, err)

	validations = []types.Validation{
		{
			Expression: "object.metadata.NAMESPACE == 'validation-test-ns'",
		},
	}

	result, validationFailure, err = GetValidationResults([]unstructured.Unstructured{deploymentResource}, validations)
	assert.Equal(t, result, false)
	assert.Equal(t, validationFailure, types.ValidationFailure{})
	expectedErrorString := "error while evaluating CEL expression: 'object.metadata.NAMESPACE == 'validation-test-ns'', error: 'no such key: NAMESPACE'"
	assert.EqualError(t, err, expectedErrorString)
}
