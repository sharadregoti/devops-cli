package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	shared "github.com/sharadregoti/devops-plugin-sdk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/sharadregoti/devops-plugin-sdk/proto"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"
)

func (h *Helm) listRepos() ([]interface{}, error) {
	// Create a new repository file object
	settings := cli.New()
	f, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load repository file: %v\n", err)
		os.Exit(1)
	}

	// List all repositories
	result := make([]interface{}, 0)
	for _, r := range f.Repositories {
		result = append(result, map[string]interface{}{
			"name": r.Name,
			"url":  r.URL,
		})
	}

	return result, nil
}

func (h *Helm) watchSwitcher(args *proto.GetResourcesArgs, ch chan interface{}) error {
	result := make([]interface{}, 0)

	switch args.ResourceType {
	case "releases":
		releases, err := h.listReleases(args)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Failed to list releases: %v", err))
			return err
		}
		result = releases

	case "namespaces":
		releases, err := h.listNamespaces(args)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Failed to list releases: %v", err))
			return err
		}
		result = releases

	case "repos":
		releases, err := h.listRepos()
		if err != nil {
			h.logger.Error(fmt.Sprintf("Failed to list repos: %v", err))
			return err
		}
		result = releases
	}

	for _, r := range result {
		ch <- r
	}

	return nil
}

func (h *Helm) watchReleases(args *proto.GetResourcesArgs, watcherDone chan struct{}) chan interface{} {
	ch := make(chan interface{}, 1)

	go func() {
		shared.LogDebug("plugin routine: resource watcher has been started for resource type (%s)", args.ResourceType)
		defer shared.LogDebug("plugin routine: resource watcher has been stopped for resource type (%s)", args.ResourceType)

		if err := h.watchSwitcher(args, ch); err != nil {
			return
		}

		for {
			select {
			case <-watcherDone:
				shared.LogTrace("plugin routine: Done received for resource type (%s)", args.ResourceType)
				close(ch)
				return

			case <-time.After(1 * time.Minute):
				if err := h.watchSwitcher(args, ch); err != nil {
					return
				}
			}
		}
	}()

	return ch
}

func (h *Helm) listReleases(args *proto.GetResourcesArgs) ([]interface{}, error) {
	actionConfig := new(action.Configuration)

	// Create a ConfigFlags struct instance with initialized values from rest.Config
	gh := genericclioptions.NewConfigFlags(true)
	gh.APIServer = &h.restConfig.Host
	gh.BearerToken = &h.restConfig.BearerToken
	gh.CAFile = &h.restConfig.CAFile
	gh.Namespace = &args.IsolatorId

	newIsolator := args.IsolatorId
	if args.IsolatorId == "all" {
		newIsolator = ""
	}

	if err := actionConfig.Init(gh, newIsolator, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}

	// if err := actionConfig.Init(NewRESTClientGetter("dev-xlr8s", h.currentKubeConfigPath), "dev-xlr8s", os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
	// 	return nil, err
	// }

	lCli := action.NewList(actionConfig)
	list, err := lCli.Run()
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, 0)
	for _, r := range list {
		result = append(result, map[string]interface{}{
			"app_version": r.Chart.AppVersion(),
			"chart":       r.Chart.Name(),
			"name":        r.Name,
			"namespace":   r.Namespace,
			"revision":    r.Version,
			"status":      r.Info.Status.String(),
			"updated":     r.Info.LastDeployed.String(),
		})
	}

	return result, nil
}

func (d *Helm) listNamespaces(args *proto.GetResourcesArgs) ([]interface{}, error) {

	nsList, err := d.normalClient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, shared.LogError("failed to list namespaces: %v", err)
	}

	result := make([]interface{}, 0)
	for _, item := range nsList.Items {
		var rawJson map[string]interface{}

		b, err := json.Marshal(item)
		if err != nil {
			return nil, shared.LogError("failed to marshal event object: %v", err)
		}

		err = json.Unmarshal(b, &rawJson)
		if err != nil {
			return nil, shared.LogError("failed to unmarshal event object: %v", err)
		}

		meta := rawJson["metadata"].(map[string]interface{})
		delete(meta, "managedFields")
		rawJson["metadata"] = meta

		result = append(result, rawJson)
	}

	return result, nil
}
