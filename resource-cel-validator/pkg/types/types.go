package types

import (
	celgo "github.com/google/cel-go/cel"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apiservercel "k8s.io/apiserver/pkg/admission/plugin/cel"
)

type KubernetesResourcePrecondition struct {
	Name           string          `json:"name,omitempty" yaml:"name,omitempty"`
	MatchResources *MatchResources `json:"matchResources,omitempty" yaml:"matchResources,omitempty"`
	Validations    []Validation    `json:"validations,omitempty" yaml:"validations,omitempty"`
}

type MatchResources struct {
	NamespaceSelector      *metav1.LabelSelector         `json:"namespaceSelector,omitempty" yaml:"namespaceSelector,omitempty"`
	ObjectSelector         *metav1.LabelSelector         `json:"objectSelector,omitempty" yaml:"objectSelector,omitempty"`
	ResourceRules          []schema.GroupVersionResource `json:"resourceRules,omitempty" yaml:"resourceRules,omitempty"`
	SelectionPreconditions []Validation                  `json:"selectionPreconditions,omitempty" yaml:"selectionPreconditions,omitempty"`
}

type Validation struct {
	Key               string `json:"key,omitempty" yaml:"key,omitempty"`
	Expression        string `json:"expression,omitempty" yaml:"expression,omitempty"`
	MessageExpression string `json:"messageExpression,omitempty" yaml:"messageExpression,omitempty"`
}

type ValidationExpressionAccessor struct {
	Expression string
}

func (v *ValidationExpressionAccessor) GetExpression() string {
	return v.Expression
}

func (v *ValidationExpressionAccessor) ReturnTypes() []*celgo.Type {
	return []*celgo.Type{celgo.BoolType}
}

type MessageExpressionAccessor struct {
	MessageExpression string
}

func (m *MessageExpressionAccessor) GetExpression() string {
	return m.MessageExpression
}

func (m *MessageExpressionAccessor) ReturnTypes() []*celgo.Type {
	return []*celgo.Type{celgo.StringType}
}

type CombinedCompilationResult struct {
	// Key is the key of the validation, can be empty
	Key string

	// ValidationCompilationResult is not optional, cannot be nil
	ValidationCompilationResult *apiservercel.CompilationResult

	// MessageCompilationResult is optional, can be nil
	MessageCompilationResult *apiservercel.CompilationResult
}

type ValidationFailure struct {
	Key      string
	Message  string
	Resource unstructured.Unstructured
}
