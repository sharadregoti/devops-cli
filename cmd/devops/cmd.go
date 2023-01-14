package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	core "github.com/sharadregoti/devops"
	"github.com/sharadregoti/devops/common"
	"github.com/spf13/cobra"
)

const VERSION = "0.2.0"

func main() {
	if err := NewCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// NewCommand return xlr8s sub commands
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devops",
		Short: "Your helping hand for DevOps",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			startTime := time.Now()

			log.Println("Checking for updates...")

			if common.Release {
				version, err := getLatestVersion()
				if err != nil {
					return err
				}

				if isGreater(version, VERSION) {
					log.Println("Update available to version:", version)
					log.Println("Upgrading devops CLI...")
					if err = runUpgrade(); err != nil {
						fmt.Println("Updgrade failed: try manual upgration")
					} else {
						return nil
					}
				}
			}

			go func() {
				common.ConnInit()
				common.IncrementAppStarts()
			}()

			defer func() {
				endTime := time.Now()

				// Report the usage time to some external service
				common.ReportUsageTime(startTime, endTime)
			}()

			fmt.Println("Documentation & Issues can be viewed at: https://github.com/sharadregoti/devops-cli")
			core.Init()
			return nil
		},
	}

	// Add all sub commands
	// cmd.AddCommand(NewInitCommand())
	cmd.AddCommand(NewVersionCommand())
	return cmd
}

// NewConsumerBindingCommand creates a consumer subcommand, which executes consumer bindings
func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "init",
		Short:  "Installs plugins as per configuration",
		PreRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(cmd *cobra.Command, args []string) error {
			// rootFS := osfs.New("/")
			// wd, _ := os.Getwd()
			// fs, _ := rootFS.Chroot(wd)
			// if err := utils.ExecuteConsumeBinding(rootFS, fs, viper.GetString("resource-id"), viper.GetString("component-name"), viper.GetString("env"), viper.GetString("kube-config-path"), viper.GetString("outputs-dir"), viper.GetString("xlr8s-url")); err != nil {
			// 	return errors.Wrap(err)
			// }
			return nil
		},
	}

	return cmd
}

func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Short:  "Prints the version of devops CLI",
		PreRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(VERSION)
			// rootFS := osfs.New("/")
			// wd, _ := os.Getwd()
			// fs, _ := rootFS.Chroot(wd)
			// if err := utils.ExecuteConsumeBinding(rootFS, fs, viper.GetString("resource-id"), viper.GetString("component-name"), viper.GetString("env"), viper.GetString("kube-config-path"), viper.GetString("outputs-dir"), viper.GetString("xlr8s-url")); err != nil {
			// 	return errors.Wrap(err)
			// }
			return nil
		},
	}

	return cmd
}

func getLatestVersion() (string, error) {
	// Replace with the URL of the public git repository you want to fetch the latest tag from
	url := "https://api.github.com/repos/sharadregoti/devops-cli/releases/latest"

	// Make an HTTP GET request to the GitHub API
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("failed to fetch latest version: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to fetch latest version: %v\n", err)
		return "", err
	}

	type Release struct {
		TagName string `json:"tag_name"`
	}

	// Unmarshal the JSON response into a Release struct
	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		log.Printf("failed to unmarshal json response: %v\n", err)
		return "", err
	}

	return release.TagName, nil
}
