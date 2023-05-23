package validation

import (
	"testing"

	"github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/types"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestValidatePreconditionData(t *testing.T) {
	precondition := &types.KubernetesResourcePrecondition{
		Name: "duplicate-key-test",
		MatchResources: &types.MatchResources{
			ResourceRules: []schema.GroupVersionResource{
				{
					Group:    "apps",
					Version:  "v1",
					Resource: "deployments",
				},
			},
		},
	}
	err := ValidatePreconditionData(precondition)
	expectedErrString := "no validations specified, at least one validation must be specified"
	assert.EqualError(t, err, expectedErrString)

	precondition = &types.KubernetesResourcePrecondition{
		Name: "duplicate-key-test",
	}
	err = ValidatePreconditionData(precondition)
	expectedErrString = "no match resources specified, match resources must be specified"
	assert.EqualError(t, err, expectedErrString)

	precondition = &types.KubernetesResourcePrecondition{
		Name:           "duplicate-key-test",
		MatchResources: &types.MatchResources{},
	}
	err = ValidatePreconditionData(precondition)
	expectedErrString = "no resource rules specified, at least one resource rule must be specified"
	assert.EqualError(t, err, expectedErrString)

	precondition = &types.KubernetesResourcePrecondition{}
	err = ValidatePreconditionData(precondition)
	expectedErrString = "no name specified, name must be specified"
	assert.EqualError(t, err, expectedErrString)

	precondition = &types.KubernetesResourcePrecondition{
		Name: "duplicate-key-test",
		MatchResources: &types.MatchResources{
			ResourceRules: []schema.GroupVersionResource{
				{
					Group:    "apps",
					Version:  "v1",
					Resource: "deployments",
				},
			},
			SelectionPreconditions: []types.Validation{
				{
					Key:        "duplicate-key-test",
					Expression: "true",
				},
				{
					Key:        "duplicate-key-test",
					Expression: "true",
				},
			},
		},
		Validations: []types.Validation{
			{
				Expression: "true",
			},
		},
	}

	err = ValidatePreconditionData(precondition)
	expectedErrString = "duplicate selection precondition key: duplicate-key-test"
	assert.EqualError(t, err, expectedErrString)

	precondition = &types.KubernetesResourcePrecondition{
		Name: "duplicate-key-test",
		MatchResources: &types.MatchResources{
			ResourceRules: []schema.GroupVersionResource{
				{
					Group:    "apps",
					Version:  "v1",
					Resource: "deployments",
				},
			},
		},
		Validations: []types.Validation{
			{
				Key:        "duplicate-key-test",
				Expression: "true",
			},
			{
				Key:        "duplicate-key-test",
				Expression: "true",
			},
		},
	}

	err = ValidatePreconditionData(precondition)
	expectedErrString = "duplicate validation key: duplicate-key-test"
	assert.EqualError(t, err, expectedErrString)
}
