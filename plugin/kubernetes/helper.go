package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/sharadregoti/devops/shared"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func getResourcesDynamically(c chan shared.WatchResourceResult, dynamic dynamic.Interface, ctx context.Context,
	group string, version string, resource string, namespace string) error {

	resourceId := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	fmt.Println(group, ":", version, ":", resource, ":", namespace, ":", resourceId)
	list, err := dynamic.Resource(resourceId).Namespace(namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	result := make([]interface{}, 0)
	for e := range list.ResultChan() {
		obj2, ok := e.Object.(interface{})
		if !ok {
			continue
		}
		result = append(result, obj2)

		obj, ok := e.Object.(*unstructured.Unstructured)
		if !ok {
			continue
		}

		c <- shared.WatchResourceResult{Type: strings.ToLower(string(e.Type)), Result: obj.UnstructuredContent()}

		// switch e.Type {
		// case watch.Added:
		// 	fmt.Println("Added:")
		// 	data, _ := json.MarshalIndent(obj, " ", " ")
		// 	fmt.Println(string(data))
		// 	c <- shared.WatchResourceResult{Type: strings.ToLower(string(e.Type)), Result: obj.UnstructuredContent()}
		// case watch.Deleted:
		// 	fmt.Println("Deleted:")
		// 	data, _ := json.MarshalIndent(obj, " ", " ")
		// 	fmt.Println(string(data))
		// case watch.Modified:
		// 	fmt.Println("Updated:")
		// 	data, _ := json.MarshalIndent(obj, " ", " ")
		// 	fmt.Println(string(data))
		// }
	}

	// result := make([]interface{}, 0)
	// for _, item := range list.Items {
	// 	var rawJson interface{}
	// 	err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &rawJson)
	// 	if err != nil {
	// 		return  err
	// 	}
	// 	result = append(result, rawJson)
	// }

	return nil
}

func listResourcesDynamically(dynamic dynamic.Interface, ctx context.Context, group string, version string, resource string, namespace string) ([]interface{}, error) {

	resourceId := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	fmt.Println(group, ":", version, ":", resource, ":", namespace, ":", resourceId)
	list, err := dynamic.Resource(resourceId).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, 0)
	for _, item := range list.Items {
		var rawJson interface{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &rawJson)
		if err != nil {
			return nil, err
		}
		result = append(result, rawJson)
	}

	return result, nil
}
