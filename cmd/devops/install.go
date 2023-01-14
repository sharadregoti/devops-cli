package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var installScript = `#!/usr/bin/env bash

# Execute this file using command "bash install.sh"

# Devops CLI location
: ${SPACE_CLI_INSTALL_DIR:="/usr/local/bin"}

# sudo is required to copy binary to SPACE_CLI_INSTALL_DIR for linux
: ${USE_SUDO:="false"}

# Http request CLI
SPACE_CLI_HTTP_REQUEST_CLI=curl

# GitHub Organization and repo name to download release
GITHUB_ORG=sharadregoti
GITHUB_REPO=devops-cli

# Dapr CLI filename
SPACE_CLI_FILENAME=devops

SPACE_CLI_FILE="${SPACE_CLI_INSTALL_DIR}/${SPACE_CLI_FILENAME}"

getSystemInfo() {
    ARCH=$(uname -m)
    case $ARCH in
        armv7*) ARCH="arm";;
        aarch64) ARCH="arm64";;
        x86_64) ARCH="x86_64";;
        386) ARCH="i386";;
    esac

    OS=$(echo ` + "`uname`" + `|tr '[:upper:]' '[:lower:]')

    # Most linux distro needs root permission to copy the file to /usr/local/bin
    if [ "$OS" == "linux" ] && [ "$SPACE_CLI_INSTALL_DIR" == "/usr/local/bin" ]; then
        USE_SUDO="true"
    fi
}

verifySupported() {

    local supported=(darwin-amd64 linux-amd64 darwin-x86_64 linux-x86_64 darwin-i386 linux-i386)

    local current_osarch="${OS}-${ARCH}"

    for osarch in "${supported[@]}"; do
        if [ "$osarch" == "$current_osarch" ]; then
            echo "Your system is ${OS}_${ARCH}"
            return
        fi
    done

    echo "No prebuilt binary for ${current_osarch}"
    exit 1
}

runAsRoot() {
    local CMD="$*"

    if [ $EUID -ne 0 -a $USE_SUDO = "true" ]; then
        CMD="sudo $CMD"
    fi

    $CMD
}

checkHttpRequestCLI() {
    if type "curl" > /dev/null; then
        SPACE_CLI_HTTP_REQUEST_CLI=curl
    elif type "wget" > /dev/null; then
        SPACE_CLI_HTTP_REQUEST_CLI=wget
    else
        echo "Either curl or wget is required"
        exit 1
    fi
}

checkExistingDapr() {
    if [ -f "$SPACE_CLI_FILE" ]; then
        echo -e "\nDevops CLI is detected:"
        $SPACE_CLI_FILE version
        echo -e "Reinstalling Devops CLI - ${SPACE_CLI_FILE}...\n"
    else
        echo -e "Installing Devops CLI...\n"
    fi
}

getLatestRelease() {
    local spaceCliReleaseUrl="https://api.github.com/repos/${GITHUB_ORG}/${GITHUB_REPO}/releases"
    local latest_release=""

    if [ "$SPACE_CLI_HTTP_REQUEST_CLI" == "curl" ]; then
        latest_release=$(curl -s $spaceCliReleaseUrl | grep \"tag_name\" | grep -v rc | awk 'NR==1{print $2}' |  sed -n 's/\"\(.*\)\",/\1/p')
    else
        latest_release=$(wget -q --header="Accept: application/json" -O - $spaceCliReleaseUrl | grep \"tag_name\" | grep -v rc | awk 'NR==1{print $2}' |  sed -n 's/\"\(.*\)\",/\1/p')
    fi

    ret_val=$latest_release
}

downloadFile() {
    LATEST_RELEASE_TAG=$1

    SPACE_CLI_ARTIFACT="devops_${LATEST_RELEASE_TAG}_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_BASE="https://storage.googleapis.com/devops-cli-artifacts/releases/devops/${LATEST_RELEASE_TAG}"
    DOWNLOAD_URL="${DOWNLOAD_BASE}/${SPACE_CLI_ARTIFACT}"

    # Create the temp directory
    SPACE_TMP_ROOT=$(mktemp -dt space-cli-install-XXXXXX)
    ARTIFACT_TMP_FILE="${SPACE_TMP_ROOT}/${SPACE_CLI_ARTIFACT}"

    echo "Downloading $DOWNLOAD_URL ..."
    if [ "$SPACE_CLI_HTTP_REQUEST_CLI" == "curl" ]; then
        curl "$DOWNLOAD_URL" -o "$ARTIFACT_TMP_FILE"
    else
        wget -q -O "$ARTIFACT_TMP_FILE" "$DOWNLOAD_URL"
    fi

    if [ ! -f "$ARTIFACT_TMP_FILE" ]; then
        echo "failed to download $DOWNLOAD_URL ..."
        exit 1
    fi

    
    # For plugins
    # Check if the ".devops" directory already exists in the home directory
    if [ ! -d "$HOME/.devops" ]; then
    # If it doesn't exist, create it
    mkdir "$HOME/.devops"
    fi

    # Check if the ".devops/plugins" directory already exists in the home directory
    if [ ! -d "$HOME/.devops/plugins" ]; then
    # If it doesn't exist, create it
    mkdir "$HOME/.devops/plugins"
    fi

    # Check if the ".devops/plugins" directory already exists in the home directory
    if [ ! -d "$HOME/.devops/plugins/kubernetes" ]; then
    # If it doesn't exist, create it
    mkdir "$HOME/.devops/plugins/kubernetes"
    fi

    SPACE_CLI_ARTIFACT="devops-kubernetes-plugin_${LATEST_RELEASE_TAG}_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="${DOWNLOAD_BASE}/${SPACE_CLI_ARTIFACT}"

    echo "Downloading kubernetes plugin"

    wget -O "$HOME/.devops/plugins/kubernetes/devops.tar.gz" "${DOWNLOAD_URL}"
    tar -xzf "$HOME/.devops/plugins/kubernetes/devops.tar.gz" -C "$HOME/.devops/plugins/kubernetes"

}

installFile() {
    tar xf "$ARTIFACT_TMP_FILE" -C "$SPACE_TMP_ROOT"
    local tmp_root_space_cli="$SPACE_TMP_ROOT/$SPACE_CLI_FILENAME"

    if [ ! -f "$tmp_root_space_cli" ]; then
        echo "Failed to unpack Devops CLI executable."
        exit 1
    fi

    chmod o+x $tmp_root_space_cli
    echo "Removing old installation"
    runAsRoot rm "$SPACE_CLI_INSTALL_DIR/devops"
    runAsRoot cp "$tmp_root_space_cli" "$SPACE_CLI_INSTALL_DIR"

    if [ -f "$SPACE_CLI_FILE" ]; then
        echo "$SPACE_CLI_FILENAME installed into $SPACE_CLI_INSTALL_DIR successfully."
        RED='\033[0;34m'
        NC='\033[0m' # No Color
        # echo -e "For enabling auto complete follow instructions provided by this command ${RED}space-cli completion --help${NC}"
        # $SPACE_CLI_FILE --version
    else
        echo "Failed to install $SPACE_CLI_FILENAME"
        exit 1
    fi
}

fail_trap() {
    result=$?
    if [ "$result" != "0" ]; then
        echo "Failed to install Devops CLI"
        echo "For support, go to https://github.com/sharadregoti/devops-cli"
    fi
    cleanup
    exit $result
}

cleanup() {
    if [[ -d "${SPACE_TMP_ROOT:-}" ]]; then
        rm -rf "SPACE_TMP_ROOT"
    fi
}

installCompleted() {
    echo -e "\nTo get started with Space Cloud, please visit https://github.com/sharadregoti/devops-cli"
}

# -----------------------------------------------------------------------------
# main
# -----------------------------------------------------------------------------
trap "fail_trap" EXIT

getSystemInfo
verifySupported
checkExistingDapr
checkHttpRequestCLI

getLatestRelease
downloadFile $ret_val
installFile
cleanup

installCompleted
`

func runUpgrade() error {
	// create a temporary file
	f, err := ioutil.TempFile("", "script")
	if err != nil {
		// handle error
		fmt.Println("failed to create temp file: Error ", err)
		return err
	}
	defer os.Remove(f.Name()) // delete the file when we're done

	// write the script to the file
	_, err = f.WriteString(installScript)
	if err != nil {
		// handle error
		f.Close()
		fmt.Println("failed to write data in temp file: Error ", err)
		return err
	}
	f.Close()

	// make the file executable
	err = os.Chmod(f.Name(), 0777)
	if err != nil {
		// handle error
		fmt.Println("failed to make the temp file executable: Error ", err)
		return err
	}

	// run the script
	cmd := exec.Command(f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		// handle error
		fmt.Println("failed to start cmd: Error ", err)
		return err
	}

	err = cmd.Wait()
	if err != nil {
		// handle error
		fmt.Println("failed to wait for cmd: Error ", err)
		return err
	}

	return nil
}
