package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sharadregoti/devops/common"
	"github.com/sharadregoti/devops/shared"
	"github.com/tidwall/gjson"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/dynamic"
)

func getResourcesDynamically(c chan shared.WatchResourceResult, dynamic dynamic.Interface, ctx context.Context, group string, version string, resource string, namespace string) error {
	resourceId := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

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

func (d *Kubernetes) getPodLogs(resourceName, namespace, containerName string) (string, error) {
	// Set the command to execute
	command := "kubectl"

	if containerName == "" {
		cont, err := d.getContainers(context.Background(), namespace, resourceName)
		if err != nil {
			return "", err
		}
		containerName = cont
	}

	arguments := []string{command, "logs", resourceName, "-n", namespace, "-f", containerName}

	d.logger.Debug(fmt.Sprintf("Fetching logs for %s %v", command, arguments))

	// cc := exec.Command(command, arguments...)

	// cc.Stdin = os.Stdin
	// cc.Stderr = os.Stderr
	// cc.Stdout = os.Stdout

	// if err := cc.Start(); err != nil {
	// 	common.Error(d.logger, fmt.Sprintf("failed to get logs, got %v", err))
	// 	return "", err
	// }

	// go func() {
	// 	for range d.activeChans {
	// 		d.logger.Debug("Closing log resource")
	// 		if err := cc.Process.Signal(os.Interrupt); err != nil {
	// 			common.Error(d.logger, fmt.Sprintf("failed to close log stream, got %v", err))
	// 		}
	// 		return
	// 	}
	// }()

	d.logger.Debug("Log fetching started")
	return strings.Join(arguments, " "), nil
}

func (d *Kubernetes) execPod(resourceName, namespace, containerName string) (string, error) {
	// Set the command to execute
	command := "kubectl"

	if containerName == "" {
		cont, err := d.getContainers(context.Background(), namespace, resourceName)
		if err != nil {
			return "", err
		}
		containerName = cont
	}

	arguments := []string{command, "exec", resourceName, "-n", namespace, "-it", "-c", containerName, "--", "sh"}

	return strings.Join(arguments, " "), nil
}

func (d *Kubernetes) portForward(resourceName, namespace string, args map[string]interface{}) (string, error) {
	// Set the command to execute
	command := "kubectl"

	cp := args["containerPort"].(string)
	lp := args["localPort"].(string)
	addr := args["address"].(string)

	if cp == "" {
		return "", fmt.Errorf("container port not provided")
	}
	if lp == "" {
		return "", fmt.Errorf("container local port not provided")
	}
	if addr == "" {
		return "", fmt.Errorf("address not provided")
	}

	arguments := []string{"port-forward", "-n", namespace, "--address", addr, resourceName, fmt.Sprintf("%s:%s", lp, cp)}

	cmd := exec.Command(command, arguments...)
	if err := cmd.Start(); err != nil {
		return "", err
	}

	d.logger.Debug("Port forward started")
	return "", nil
}

func (d *Kubernetes) DescribeResource(resourceType, resourceName, namespace string) (string, error) {
	// Set the command to execute
	command := "kubectl"
	arguments := []string{"describe", resourceType, resourceName, "-n", namespace}

	// Execute the command
	output, err := exec.Command(command, arguments...).Output()
	if err != nil {
		common.Error(d.logger, fmt.Sprintf("failed to get describe output, got %v", err))
		return "", err
	}

	return string(output), nil
}

func (d *Kubernetes) deleteResource(ctx context.Context, args shared.ActionDeleteResourceArgs) error {
	rt, ok := d.resourceTypes[args.ResourceType]
	if !ok {
		d.logger.Debug(fmt.Sprintf("Could not find resource type %s in current kubernetes context", args.ResourceType))
		return fmt.Errorf("Delete: could not find resource type %s in current kubernetes context", args.ResourceType)
	}

	resourceId := schema.GroupVersionResource{
		Group:    rt.group,
		Version:  rt.version,
		Resource: rt.resourceTypeName,
	}

	var err error
	if rt.isNamespaced {
		err = d.dynamicClient.Resource(resourceId).Namespace(args.IsolatorName).Delete(ctx, args.ResourceName, metav1.DeleteOptions{})
	} else {
		err = d.dynamicClient.Resource(resourceId).Delete(ctx, args.ResourceName, metav1.DeleteOptions{})
	}
	if err != nil {
		return err
	}

	return err
}

func (d *Kubernetes) getContainers(ctx context.Context, namespace string, resourceName string) (string, error) {
	pod, err := d.normalClient.CoreV1().Pods(namespace).Get(ctx, resourceName, metav1.GetOptions{})
	if err != nil {
		common.Error(d.logger, fmt.Sprintf("failed to get pod %s in namespace %s, got error %v", resourceName, namespace, err))
		return "", err
	}

	for _, c := range pod.Spec.Containers {
		return c.Name, nil
	}
	return "", nil
}

func (d *Kubernetes) getPods(ctx context.Context, namespace string, resourceName, resourceType string) ([]interface{}, error) {
	arr, err := d.listResources(ctx, shared.GetResourcesArgs{
		ResourceName: resourceName,
		ResourceType: resourceType,
		IsolatorID:   namespace,
	}, "")
	if err != nil {
		return nil, err
	}

	if len(arr) == 0 {
		common.Error(d.logger, "length of resource is zero")
		return make([]interface{}, 0), err
	}

	res := arr[0].(map[string]interface{})
	strData, err := json.Marshal(res)
	if err != nil {
		common.Error(d.logger, fmt.Sprintf("failed to json marshal, got error %v", err))
		return nil, err
	}

	path := "spec.selector.matchLabels"
	if resourceType == "services" {
		path = "spec.selector"
	}
	value := gjson.Get(string(strData), path)
	if !value.IsObject() {
		return nil, nil
	}

	selector := labels.NewSelector()
	for key, v := range value.Value().(map[string]interface{}) {
		l2, _ := labels.NewRequirement(key, selection.Equals, []string{v.(string)})
		selector = selector.Add(*l2)
	}

	resultArr, err := d.listResources(ctx, shared.GetResourcesArgs{
		ResourceName: "",
		ResourceType: "pods",
		IsolatorID:   namespace,
	}, selector.String())
	// list, err := d.normalClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
	// 	LabelSelector: selector.String(),
	// })
	if err != nil {
		common.Error(d.logger, fmt.Sprintf("failed to get pod %s in namespace %s, got error %v", resourceName, namespace, err))
		return nil, err
	}

	// resultArr := make([]interface{}, 0)
	// for _, p := range list.Items {
	// 	resultArr = append(resultArr, p)
	// }

	return resultArr, nil
}

func (d *Kubernetes) listResources(ctx context.Context, args shared.GetResourcesArgs, label string) ([]interface{}, error) {
	rt, ok := d.resourceTypes[args.ResourceType]
	if !ok {
		d.logger.Debug(fmt.Sprintf("Could not find resource type %s in current kubernetes context", args.ResourceType))
		return []interface{}{}, fmt.Errorf("List: could not find resource type %s in current kubernetes context", args.ResourceType)
	}

	resourceId := schema.GroupVersionResource{
		Group:    rt.group,
		Version:  rt.version,
		Resource: rt.resourceTypeName,
	}

	var list *unstructured.UnstructuredList
	var err error

	if args.ResourceName != "" {
		// Single get
		uData, err := d.dynamicClient.Resource(resourceId).Namespace(args.IsolatorID).Get(ctx, args.ResourceName, metav1.GetOptions{})
		if err != nil {
			common.Error(d.logger, fmt.Sprintf("failed to get dynamic resource: %v", err))
			return nil, err
		}
		list = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{*uData}}
	} else {
		// List
		if rt.isNamespaced {
			list, err = d.dynamicClient.Resource(resourceId).Namespace(args.IsolatorID).List(ctx, metav1.ListOptions{LabelSelector: label})
		} else {
			list, err = d.dynamicClient.Resource(resourceId).List(ctx, metav1.ListOptions{LabelSelector: label})
		}
		if err != nil {
			return nil, err
		}
	}

	result := make([]interface{}, 0)
	for _, item := range list.Items {
		var rawJson map[string]interface{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &rawJson)
		if err != nil {
			common.Error(d.logger, fmt.Sprintf("failed to unstructure resource: %v", err))
			return nil, err
		}
		meta := rawJson["metadata"].(map[string]interface{})
		delete(meta, "managedFields")
		rawJson["metadata"] = meta
		result = append(result, rawJson)
	}

	return result, nil
}
