package validation

import (
	"context"
	"fmt"

	"github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/types"
	"github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// GetResources uses the dynamic client to fetch resources based on label selectors and resource rules.
func GetResources(resources *types.MatchResources, client *dynamic.DynamicClient) ([]unstructured.Unstructured, error) {
	namespaceList := make([]string, 0)
	if resources.NamespaceSelector != nil {
		list, err := getNamespaces(resources.NamespaceSelector, client)
		if err != nil {
			return nil, err
		}
		if len(list) == 0 {
			return nil, fmt.Errorf("no namespaces found for the given selector")
		}
		namespaceList = append(namespaceList, list...)
	} else {
		// If no namespace selector is provided, then we will search in all namespaces
		namespaceList = append(namespaceList, "")
	}

	resourceList := make([]unstructured.Unstructured, 0)
	objectListOptions := metav1.ListOptions{}
	if resources.ObjectSelector != nil {
		objectListOptions.LabelSelector = metav1.FormatLabelSelector(resources.ObjectSelector)
	}

	for _, gvr := range resources.ResourceRules {
		for _, namespace := range namespaceList {
			parentResource, subresource := util.ParseResource(gvr.Resource)
			resourceInterface := client.Resource(schema.GroupVersionResource{
				Group:    gvr.Group,
				Version:  gvr.Version,
				Resource: parentResource,
			})
			unstructuredList, err := resourceInterface.
				Namespace(namespace).
				List(context.Background(), objectListOptions)
			if err != nil {
				return nil, err
			}
			if subresource != "" {
				for _, object := range unstructuredList.Items {
					var unstructuredSubresource *unstructured.Unstructured
					unstructuredSubresource, err = getSubresource(object, resourceInterface, namespace, subresource)
					if err != nil {
						return nil, err
					}
					resourceList = append(resourceList, *unstructuredSubresource)
				}
			} else {
				resourceList = append(resourceList, unstructuredList.Items...)
			}
		}
	}
	return resourceList, nil
}

// getSubresource returns the subresource (eg- scale, status) associated with the parent resource.
func getSubresource(object unstructured.Unstructured, resourceInterface dynamic.NamespaceableResourceInterface, namespace string, subresource string) (*unstructured.Unstructured, error) {
	var unstructuredSubresource *unstructured.Unstructured
	var err error
	if object.GetNamespace() != "" {
		unstructuredSubresource, err = resourceInterface.
			Namespace(namespace).
			Get(context.Background(), object.GetName(), metav1.GetOptions{}, subresource)
	} else {
		// Non namespaced resources like nodes, namespaces, etc. do not have a namespace
		unstructuredSubresource, err = resourceInterface.
			Get(context.Background(), object.GetName(), metav1.GetOptions{}, subresource)
	}
	if err != nil {
		return nil, err
	}
	return unstructuredSubresource, nil
}

// getNamespaces returns namespaces based on namespace label selector
func getNamespaces(namespaceSelector *metav1.LabelSelector, client *dynamic.DynamicClient) ([]string, error) {
	namespaces, err := client.Resource(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "namespaces",
	}).List(context.Background(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(namespaceSelector),
	})
	if err != nil {
		return nil, err
	}
	namespaceList := make([]string, 0)
	for _, namespace := range namespaces.Items {
		namespaceList = append(namespaceList, namespace.GetName())
	}
	return namespaceList, nil
}
