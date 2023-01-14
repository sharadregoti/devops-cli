package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes/scheme"
)

func valiations() {
	// consolescanner := bufio.NewScanner(os.Stdin)

	// Check if any arguments were passed
	if len(os.Args) == 1 {
		fmt.Println("No arguments were passed.")
		return
	}

	// Print the command-line arguments
	// for i, arg := range os.Args[1:] {
	// 	fmt.Printf("Argument %d: %s\n", i+1, arg)
	// }

	f, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(consolescanner.Text())

	obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(f, nil, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
		return
	}

	// now use switch over the type of the object
	// and match each type-case
	switch o := obj.(type) {
	case *v1.Pod:
		if o.Namespace == "" {
			o.Namespace = "default"
		}
		findngs, err := podCheck(o.Namespace, &o.Spec)
		if err != nil {
			log.Fatal(err)
		}

		color.Green("Findings")
		// Print the list of strings
		for i, item := range findngs {
			fmt.Printf("%d: %s\n", i+1, item)
		}

	case *appsv1.Deployment:
		podCheck(o.Namespace, &o.Spec.Template.Spec)
	case *appsv1.StatefulSet:
		podCheck(o.Namespace, &o.Spec.Template.Spec)
	case *appsv1.DaemonSet:
		podCheck(o.Namespace, &o.Spec.Template.Spec)
	default:
		fmt.Printf("Type %v is unknown", o)
	}
}
