package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	shared "github.com/sharadregoti/devops-plugin-sdk"
	"github.com/sharadregoti/devops-plugin-sdk/proto"
	"github.com/tidwall/gjson"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
)

func (d *Kubernetes) watchResourceDynamic(ctx context.Context, args *proto.GetResourcesArgs) (watch.Interface, error) {
	shared.LogTrace("client: watchResourceDynamic called for resource type %s", args.ResourceType)

	rt, ok := d.resourceTypes[args.ResourceType]
	if !ok {
		return nil, shared.LogError("client: could not find resource type %s in current kubernetes context", args.ResourceType)
	}

	resourceId := schema.GroupVersionResource{
		Group:    rt.group,
		Version:  rt.version,
		Resource: rt.resourceTypeName,
	}

	var list watch.Interface
	var err error

	if args.IsolatorId == "all" {
		args.IsolatorId = v1.NamespaceAll
	}

	if rt.isNamespaced {
		list, err = d.dynamicClient.Resource(resourceId).Namespace(args.IsolatorId).Watch(ctx, metav1.ListOptions{})
	} else {
		list, err = d.dynamicClient.Resource(resourceId).Watch(ctx, metav1.ListOptions{})
	}
	if err != nil {
		return nil, err
	}

	return list, nil
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
		return "", shared.LogError("failed to get describe output, got %v", err)
	}

	return string(output), nil
}

func (d *Kubernetes) createResource(ctx context.Context, args *proto.ActionCreateResourceArgs) error {
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

	objMap, ok := args.Data.AsInterface().(map[string]interface{})
	if !ok {
		return fmt.Errorf("could not convert data to map[string]interface{}")
	}

	data := &unstructured.Unstructured{
		Object: objMap,
	}

	var err error
	if rt.isNamespaced {
		_, err = d.dynamicClient.Resource(resourceId).Namespace(args.IsolatorName).Create(ctx, data, metav1.CreateOptions{})
	} else {
		_, err = d.dynamicClient.Resource(resourceId).Create(ctx, data, metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}

	return err
}

func (d *Kubernetes) updateResource(ctx context.Context, args *proto.ActionUpdateResourceArgs) error {
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

	objMap, ok := args.Data.AsInterface().(map[string]interface{})
	if !ok {
		return fmt.Errorf("could not convert data to map[string]interface{}")
	}

	data := &unstructured.Unstructured{
		Object: objMap,
	}

	var err error
	if rt.isNamespaced {
		_, err = d.dynamicClient.Resource(resourceId).Namespace(args.IsolatorName).Update(ctx, data, metav1.UpdateOptions{})
	} else {
		_, err = d.dynamicClient.Resource(resourceId).Update(ctx, data, metav1.UpdateOptions{})
	}
	if err != nil {
		return err
	}

	return err
}

func (d *Kubernetes) deleteResource(ctx context.Context, args *proto.ActionDeleteResourceArgs) error {
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
		return "", shared.LogError("failed to get pod %s in namespace %s, got error %v", resourceName, namespace, err)
	}

	for _, c := range pod.Spec.Containers {
		return c.Name, nil
	}
	return "", nil
}

func (d *Kubernetes) getPods(ctx context.Context, namespace string, resourceName, resourceType string) ([]interface{}, error) {
	arr, err := d.listResources(ctx, &proto.GetResourcesArgs{
		ResourceName: resourceName,
		ResourceType: resourceType,
		IsolatorId:   namespace,
	}, "")
	if err != nil {
		return nil, err
	}

	if len(arr) == 0 {
		return make([]interface{}, 0), shared.LogError("length of resource is zero")
	}

	res := arr[0].(map[string]interface{})
	strData, err := json.Marshal(res)
	if err != nil {
		return nil, shared.LogError("failed to json marshal, got error %v", err)
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

	resultArr, err := d.listResources(ctx, &proto.GetResourcesArgs{
		ResourceName: "",
		ResourceType: "pods",
		IsolatorId:   namespace,
	}, selector.String())
	// list, err := d.normalClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
	// 	LabelSelector: selector.String(),
	// })
	if err != nil {
		return nil, shared.LogError("failed to get pod %s in namespace %s, got error %v", resourceName, namespace, err)
	}

	// resultArr := make([]interface{}, 0)
	// for _, p := range list.Items {
	// 	resultArr = append(resultArr, p)
	// }

	return resultArr, nil
}

func (d *Kubernetes) listResources(ctx context.Context, args *proto.GetResourcesArgs, label string) ([]interface{}, error) {
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

	if args.IsolatorId == "all" {
		args.IsolatorId = v1.NamespaceAll
	}

	if args.ResourceName != "" {
		// Single get
		var uData *unstructured.Unstructured
		if rt.isNamespaced {
			uData, err = d.dynamicClient.Resource(resourceId).Namespace(args.IsolatorId).Get(ctx, args.ResourceName, metav1.GetOptions{})
			if err != nil {
				return nil, shared.LogError("failed to get dynamic resource: %v", err)
			}
		} else {
			uData, err = d.dynamicClient.Resource(resourceId).Get(ctx, args.ResourceName, metav1.GetOptions{})
			if err != nil {
				return nil, shared.LogError("failed to get dynamic resource: %v", err)
			}
		}
		list = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{*uData}}
	} else {
		// List
		if rt.isNamespaced {
			list, err = d.dynamicClient.Resource(resourceId).Namespace(args.IsolatorId).List(ctx, metav1.ListOptions{LabelSelector: label})
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
			return nil, shared.LogError("failed to unstructure resource: %v", err)
		}
		meta := rawJson["metadata"].(map[string]interface{})
		delete(meta, "managedFields")
		rawJson["metadata"] = meta
		if args.ResourceType == "pods" {
			f, err := json.Marshal(rawJson)
			if err != nil {
				return nil, shared.LogError("failed to unstructure resource: %v", err)
			}

			obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(f, nil, nil)
			if err != nil {
				return nil, shared.LogError("failed to unstructure resource: %v", err)
			}
			// var po v1.Pod
			// err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &rawJson)
			// if err != nil {
			// 	return nil, shared.LogError(// 	common.Error(d.logger, fmt.Sprintf("failed to unstructure resource: %v", err)))
			// }
			// rawJson["customCalculatedStatus"] = Phase(&po)
			rawJson["devops"] = map[string]interface{}{
				"customCalculatedStatus": Phase(obj.(*v1.Pod)),
			}
		}
		result = append(result, rawJson)
	}

	return result, nil
}

// Phase reports the given pod phase.
func Phase(po *v1.Pod) string {
	status := string(po.Status.Phase)
	if po.Status.Reason != "" {
		if po.DeletionTimestamp != nil && po.Status.Reason == "NodeLost" {
			return "Unknown"
		}
		status = po.Status.Reason
	}

	status, ok := initContainerPhase(po.Status, len(po.Spec.InitContainers), status)
	if ok {
		return status
	}

	status, ok = containerPhase(po.Status, status)
	if ok && status == "Completed" {
		status = string(v1.PodRunning)
	}
	if po.DeletionTimestamp == nil {
		return status
	}

	return "Terminating"
}

func containerPhase(st v1.PodStatus, status string) (string, bool) {
	var running bool
	for i := len(st.ContainerStatuses) - 1; i >= 0; i-- {
		cs := st.ContainerStatuses[i]
		switch {
		case cs.State.Waiting != nil && cs.State.Waiting.Reason != "":
			status = cs.State.Waiting.Reason
		case cs.State.Terminated != nil && cs.State.Terminated.Reason != "":
			status = cs.State.Terminated.Reason
		case cs.State.Terminated != nil:
			if cs.State.Terminated.Signal != 0 {
				status = "Signal:" + strconv.Itoa(int(cs.State.Terminated.Signal))
			} else {
				status = "ExitCode:" + strconv.Itoa(int(cs.State.Terminated.ExitCode))
			}
		case cs.Ready && cs.State.Running != nil:
			running = true
		}
	}

	return status, running
}

func initContainerPhase(st v1.PodStatus, initCount int, status string) (string, bool) {
	for i, cs := range st.InitContainerStatuses {
		s := checkContainerStatus(cs, i, initCount)
		if s == "" {
			continue
		}
		return s, true
	}

	return status, false
}

// ----------------------------------------------------------------------------
// Helpers..

func checkContainerStatus(cs v1.ContainerStatus, i, initCount int) string {
	switch {
	case cs.State.Terminated != nil:
		if cs.State.Terminated.ExitCode == 0 {
			return ""
		}
		if cs.State.Terminated.Reason != "" {
			return "Init:" + cs.State.Terminated.Reason
		}
		if cs.State.Terminated.Signal != 0 {
			return "Init:Signal:" + strconv.Itoa(int(cs.State.Terminated.Signal))
		}
		return "Init:ExitCode:" + strconv.Itoa(int(cs.State.Terminated.ExitCode))
	case cs.State.Waiting != nil && cs.State.Waiting.Reason != "" && cs.State.Waiting.Reason != "PodInitializing":
		return "Init:" + cs.State.Waiting.Reason
	default:
		return "Init:" + strconv.Itoa(i) + "/" + strconv.Itoa(initCount)
	}
}
