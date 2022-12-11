package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func checkSecret(secretName, namespace string, keys []string) ([]string, error) {
	findings := []string{}
	// Create a new Kubernetes client
	// Load the kubeconfig file
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		// Get the user's home directory
		home, _ := os.UserHomeDir()

		// Append the default kubeconfig file path to the home directory
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return findings, nil
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return findings, nil
	}

	// Check if the secret exists
	secretRes, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return append(findings, fmt.Sprintf("secret refered %v does not exists", secretName)), nil
	} else if !errors.IsNotFound(err) {
		// Check if the keys exist in the secret
		for _, key := range keys {
			_, ok := secretRes.Data[key]
			if !ok {
				return append(findings, fmt.Sprintf("key '%s' does not exist in secret '%s'", key, secretName)), nil
			}
		}
	}
	if err != nil {
		return findings, nil
	}

	return findings, nil
}

func checkConfigMap(configMapName, namespace string, keys []string) ([]string, error) {
	findings := []string{}
	// Create a new Kubernetes client
	// Load the kubeconfig file
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		// Get the user's home directory
		home, _ := os.UserHomeDir()

		// Append the default kubeconfig file path to the home directory
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return findings, nil
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return findings, nil
	}

	// Check if the ConfigMap exists
	configMapData, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return append(findings, fmt.Sprintf("configmap refered %v does not exists", configMapName)), nil
	} else if !errors.IsNotFound(err) {
		// Check if the keys exist in the secret
		for _, key := range keys {
			_, ok := configMapData.Data[key]
			if !ok {
				return append(findings, fmt.Sprintf("key '%s' does not exist in configmap '%s'", key, configMapName)), nil
			}
		}
	}
	if err != nil {
		return findings, nil
	}

	return findings, nil
}
