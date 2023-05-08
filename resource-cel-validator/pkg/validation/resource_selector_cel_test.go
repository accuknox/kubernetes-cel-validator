package validation

import (
	"testing"

	"github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/types"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var deploymentResource = unstructured.Unstructured{
	Object: map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"labels": map[string]interface{}{
				"app": "nginx",
			},
			"name":      "nginx",
			"namespace": "validation-test-ns",
		},
		"spec": map[string]interface{}{
			"replicas": 1,
			"selector": map[string]interface{}{
				"matchLabels": map[string]interface{}{
					"app": "nginx",
				},
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"app": "nginx",
					},
				},
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"image": "nginx",
							"name":  "nginx",
						},
					},
				},
			},
		},
	},
}

var serviceResource = unstructured.Unstructured{
	Object: map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata": map[string]interface{}{
			"labels": map[string]interface{}{
				"app": "nginx",
			},
			"name":      "nginx",
			"namespace": "validation-test-ns-2",
		},
		"spec": map[string]interface{}{
			"ports": []interface{}{
				map[string]interface{}{
					"port": 80,
				},
			},
			"selector": map[string]interface{}{
				"app": "nginx",
			},
		},
	},
}

func TestFilterResourcesOnSelectionPreconditions(t *testing.T) {

	selectionPreconditions := []types.Validation{
		{
			Expression: "object.metadata.namespace == 'validation-test-ns'",
		},
		{
			Key:               "nginx-name",
			Expression:        "object.metadata.name == 'nginx'",
			MessageExpression: "object.metadata.name + ' is not nginx'",
		},
	}

	filteredResource, selectionPreconditionValidationFailures, err := FilterResourcesOnSelectionPreconditions([]unstructured.Unstructured{deploymentResource}, selectionPreconditions)
	assert.Equal(t, deploymentResource, filteredResource[0])
	assert.Empty(t, selectionPreconditionValidationFailures)
	assert.Nil(t, err)

	selectionPreconditions = []types.Validation{
		{
			Expression: "object.metadata.foo == 'validation-test-ns'",
		},
	}
	filteredResource, selectionPreconditionValidationFailures, err = FilterResourcesOnSelectionPreconditions([]unstructured.Unstructured{deploymentResource}, selectionPreconditions)
	assert.Empty(t, filteredResource)
	assert.Empty(t, selectionPreconditionValidationFailures)

	expectedErrorString := "error while evaluating CEL expression: 'object.metadata.foo == 'validation-test-ns'', error: 'no such key: foo'"
	assert.EqualError(t, err, expectedErrorString)

	selectionPreconditions = []types.Validation{
		{
			Expression: "object.metadata.namespace == 'validation-test-ns'",
		},
		{
			Expression: "object.metadata.NAME == 'nginx'",
		},
	}
	filteredResource, selectionPreconditionValidationFailures, err = FilterResourcesOnSelectionPreconditions([]unstructured.Unstructured{deploymentResource}, selectionPreconditions)
	assert.Empty(t, filteredResource)
	assert.Empty(t, selectionPreconditionValidationFailures)

	expectedErrorString = "error while evaluating CEL expression: 'object.metadata.NAME == 'nginx'', error: 'no such key: NAME'"
	assert.EqualError(t, err, expectedErrorString)

	selectionPreconditions = []types.Validation{}
	filteredResource, selectionPreconditionValidationFailures, err = FilterResourcesOnSelectionPreconditions([]unstructured.Unstructured{deploymentResource}, selectionPreconditions)

	assert.Equal(t, deploymentResource, filteredResource[0])
	assert.Empty(t, selectionPreconditionValidationFailures)
	assert.Nil(t, err)
}

func TestSelectionPreconditionValidationFailuresInFilterResourcesOnSelectionPreconditions(t *testing.T) {
	selectionPreconditions := []types.Validation{
		{
			Key:               "nginx-2-name",
			Expression:        "object.metadata.name == 'nginx-2'",
			MessageExpression: "object.metadata.name + ' is not nginx-2'",
		},
	}
	expectedSelectionPreconditionValidationFailures := []types.ValidationFailure{
		{
			Key:      "nginx-2-name",
			Message:  "nginx is not nginx-2",
			Resource: deploymentResource,
		},
	}

	filteredResource, selectionPreconditionValidationFailures, err := FilterResourcesOnSelectionPreconditions([]unstructured.Unstructured{deploymentResource}, selectionPreconditions)
	assert.Empty(t, filteredResource)
	assert.Equal(t, expectedSelectionPreconditionValidationFailures, selectionPreconditionValidationFailures)
	assert.Nil(t, err)

	selectionPreconditions = []types.Validation{
		{
			Key:               "nginx-name",
			Expression:        "object.metadata.name == 'nginx'",
			MessageExpression: "object.metadata.name + ' is not nginx'",
		},
		{
			Key:               "nginx-2-name",
			Expression:        "object.metadata.name == 'nginx-2'",
			MessageExpression: "object.metadata.name + ' is not nginx-2'",
		},
	}

	filteredResource, selectionPreconditionValidationFailures, err = FilterResourcesOnSelectionPreconditions([]unstructured.Unstructured{deploymentResource}, selectionPreconditions)
	assert.Empty(t, filteredResource)
	assert.Equal(t, expectedSelectionPreconditionValidationFailures, selectionPreconditionValidationFailures)
	assert.Nil(t, err)

	selectionPreconditions = []types.Validation{
		{
			Key:               "nginx-name",
			Expression:        "object.metadata.name == 'nginx'",
			MessageExpression: "object.metadata.name + ' is not nginx'",
		},
		{
			Key:               "kind-service",
			Expression:        "object.kind == 'Service'",
			MessageExpression: "object.kind + ' is not Service'",
		},
	}
	expectedSelectionPreconditionValidationFailures = []types.ValidationFailure{
		{
			Key:      "kind-service",
			Message:  "Deployment is not Service",
			Resource: deploymentResource,
		},
	}

	filteredResource, selectionPreconditionValidationFailures, err = FilterResourcesOnSelectionPreconditions([]unstructured.Unstructured{deploymentResource, serviceResource}, selectionPreconditions)
	assert.Equal(t, serviceResource, filteredResource[0])
	assert.Equal(t, expectedSelectionPreconditionValidationFailures, selectionPreconditionValidationFailures)
	assert.Nil(t, err)
}
