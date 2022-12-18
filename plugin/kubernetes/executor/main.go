package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	k "github.com/sharadregoti/devops/plugin/kubernetes"
	"github.com/sharadregoti/devops/shared"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	d, err := k.New(logger)
	fmt.Println(err)

	// k := New(logger)
	// fmt.Println(d.GetGeneralInfo())
	// fmt.Println(d.Name())
	// fmt.Println(d.GetResourceTypeList())
	// d.DescribeResource()
	res, _ := d.GetResources(shared.GetResourcesArgs{ResourceType: "namespaces", IsolatorID: ""})
	fmt.Println(len(res))
}

// package main

// import (
// 	"context"
// 	"fmt"

// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	"k8s.io/apimachinery/pkg/runtime/schema"
// 	"k8s.io/client-go/dynamic"
// 	ctrl "sigs.k8s.io/controller-runtime"
// )

// func main() {
// 	ctx := context.Background()
// 	config := ctrl.GetConfigOrDie()
// 	dynamic := dynamic.NewForConfigOrDie(config)

// 	namespace := "default"
// 	items, err := GetResourcesDynamically(dynamic, ctx,
// 		"apps", "v1", "deployments", namespace)
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		for _, item := range items {
// 			fmt.Printf("%+v\n", item)
// 		}
// 	}
// }

// func GetResourcesDynamically(dynamic dynamic.Interface, ctx context.Context,
// 	group string, version string, resource string, namespace string) (
// 	[]unstructured.Unstructured, error) {

// 	resourceId := schema.GroupVersionResource{
// 		Group:    group,
// 		Version:  version,
// 		Resource: resource,
// 	}
// 	list, err := dynamic.Resource(resourceId).Namespace(namespace).
// 		List(ctx, metav1.ListOptions{})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return list.Items, nil
// }
