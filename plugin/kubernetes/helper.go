package kubernetes

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sharadregoti/devops/shared"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func getResourcesDynamically(c chan shared.WatchResourceResult, dynamic dynamic.Interface, ctx context.Context, group string, version string, resource string, namespace string) error {
	resourceId := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	// fmt.Println(group, ":", version, ":", resource, ":", namespace, ":", resourceId)
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

func (d *Kubernetes) DescribeResource(resourceType, resourceName, namespace string) (string, error) {
	// Set the command to execute
	command := "kubectl"
	arguments := []string{"describe", resourceType, resourceName, "-n", namespace}

	// Execute the command
	output, err := exec.Command(command, arguments...).Output()
	if err != nil {
		d.logger.Error("failed to execute describe command", arguments, err)
		return "", err
	}

	// Print the command output
	return string(output), nil
}

// func (d *Kubernetes) getLogs(dynamic dynamic.Interface, ctx context.Context, group string, version string, resource string, namespace string, resourceName string) error {
// 	r := d.normalClient.CoreV1().Pods("").GetLogs("", &v1.PodLogOptions{})
// 	d.normalClient.CoreV1()
// 	res := r.Do(ctx)
// 	// res.
// }

func deleteResourcesDynamically(dynamic dynamic.Interface, ctx context.Context, group string, version string, resource string, namespace string, resourceName string) error {
	resourceId := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	err := dynamic.Resource(resourceId).Namespace(namespace).Delete(ctx, resourceName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return err
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
