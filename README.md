# Kubernetes CEL Validator

## Description

This is a simple tool to validate a Kubernetes cluster resources against a set of rules defined using the [CEL](https://github.com/google/cel-spec) language.

The CEL used uses both [CEL community libraries](https://kubernetes.io/docs/reference/using-api/cel/#cel-community-libraries) and [Kubernetes CEL libraries](https://kubernetes.io/docs/reference/using-api/cel/#kubernetes-cel-libraries).

Please read the CEL Introduction at https://kubernetes.io/docs/reference/using-api/cel/ to know about working with CEL.

## Usage

This is a library that can be used in your own Go code. You can use it in your own code by importing it:

```go
package testvalidate

import (
	celvalidator "github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg"
	celvalidatortypes "github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func ValidateResource(config *rest.Config) (bool, []celvalidatortypes.ValidationFailure, celvalidatortypes.ValidationFailure, error) {
	kubernetesResourcePrecondition := celvalidatortypes.KubernetesResourcePrecondition{
		Name: "test-precondition",
		MatchResources: &celvalidatortypes.MatchResources{
			NamespaceSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"test": "test",
				},
			},
			ObjectSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"test": "test",
				},
			},
			ResourceRules: []schema.GroupVersionResource{
				{
					Group:    "apps",
					Version:  "v1",
					Resource: "deployments",
				},
			},
			SelectionPreconditions: []celvalidatortypes.Validation{
				{
					Key:               "test-label-using-cel",
					Expression:        "object.metadata.labels.foo == 'bar'",
					MessageExpression: "resource.metadata.name + ' does not have label foo=bar'",
				},
			},
		},
		Validations: []celvalidatortypes.Validation{
			{
				Key:               "test-replicas",
				Expression:        "object.spec.replicas > 1",
				MessageExpression: "resource.metadata.name + ' has more than 1 replica'",
			},
		},
	}

	validationResult, selectionPreconditionValidationFailures, validationFailure, err := celvalidator.GetKubernetesResourcePreconditionResult(&kubernetesResourcePrecondition, config)
	return validationResult, selectionPreconditionValidationFailures, validationFailure, err
}
```

The method `celvalidator.GetKubernetesResourcePreconditionResult(..)` returns:
1. A boolean value indicating the result of the validation.
2. A list of `celvalidatortypes.ValidationFailure` objects which contain details of object that failed the `**SelectionPreconditions**`.
3. A `celvalidatortypes.ValidationFailure` object which contains details of the object that failed the `**Validations**`.
4. An error that might be encountered during validation phase, resource selection phase or during the creation of the Kubernetes client.

## Field Specification

It is really common to use YAML to specify the rules and then unmarshal it into the `celvalidatortypes.KubernetesResourcePrecondition` object. The following is an example of the YAML specification:

```yaml
name: replicas-precondition # Mandatory, name of the precondition
matchResources:             # Mandatory, resource selection criteria
  namespaceSelector:        # Optional, namespace selection criteria
    matchExpressions:       # Optional, namespace selection criteria
      - key: hoo
        operator: DoesNotExist
    matchLabels:            # Optional, namespace selection criteria
      kubernetes.io/metadata.name: validation-test-ns
  objectSelector:           # Optional, object selection criteria
    matchExpressions:       # Optional, object selection criteria
      - key: hoo
        operator: DoesNotExist
    matchLabels:            # Optional, object selection criteria
      app: nginx
  resourceRules:            # Mandatory, resource selection criteria, at least one resource rule is required
    - Group: apps
      Version: v1
      Resource: deployments
  selectionPreconditions:   # Optional, selection preconditions
    - key: "nginx-name"     # Optional, key of the precondition
      messageExpression: "'resource: ' + object.metadata.name + ' is not nginx'" # Optional, message to be displayed when the precondition fails
      expression: "object.metadata.name == 'nginx'"                              # Mandatory, CEL expression to be evaluated
validations:                                # Mandatory, validations, at least one validation is required
  - expression: "object.spec.replicas == 1" # Mandatory, CEL expression to be evaluated
    key: "single-replica"                   # Optional, key of the validation
    messageExpression: "'resource: ' + object.metadata.name + ' does not have 1 replica, it has ' + string(object.spec.replicas)" 
                                            # ⬆ Optional, message to be displayed when the validation fails
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -s -m'Add some feature'`) (Commits need to be signed)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
